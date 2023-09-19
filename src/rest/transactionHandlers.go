package rest

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"encoding/csv"
	"io"
	"strconv"

	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

func deleteTransaction(tpm managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := managing.DeleteTransactionReq{
			ID: c.Params("id"),
		}
		slog.Info("Transaction delete requested.", "id", req.ID)
		resp, err := tpm.DeleteTransaction(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slog.Info("Transaction deleted", "trasnaction", resp.ID)
		c.Append("HX-Trigger", "newTransaction")
		return c.Render("toastOk", fiber.Map{
			"Title": "Deleted",
			"Msg":   "Transaction deleted",
		})
	}
}

func newTransactionImport(tpm managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Importing transactions")
		file, err := c.FormFile("trancsv")
		if err != nil {
			slog.Error("error", err)
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintln("Error reading CSV req:", err),
			})

		}
		uploadedFile, err := file.Open()
		if err != nil {
			slog.Error("error", err)
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   fmt.Sprintln("Error reading CSV req:", err),
			})

		}
		defer uploadedFile.Close()

		csvReader := csv.NewReader(uploadedFile)

		// Iterate over the CSV records
		csvReader.Read() // Skip column
		tCount := 0
		reqs := []managing.RecordTransactionReq{}
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return c.Render("toastErr", fiber.Map{
					"Title": "Error",
					"Msg":   fmt.Sprintln("Error reading CSV req:", err),
				})
			}
			req := managing.RecordTransactionReq{}
			req.TradingPairID = c.Params("id")
			req.Type = line[0]
			req.Timestamp, err = time.Parse("2006-01-02 15:04", line[1])
			if err != nil {
				slog.Error("Error parsing imported transaction.", "error", err)
				return c.Render("toastErr", fiber.Map{"Title": "Created", "Msg": fmt.Sprintf("Error on row %d col %d", tCount, 0)})
			}
			req.BaseAmount, err = strconv.ParseFloat(line[2], 64)
			if err != nil {
				slog.Error("Error parsing imported transaction.", "error", err)
				return c.Render("toastErr", fiber.Map{"Title": "Created", "Msg": fmt.Sprintf("Error on row %d col %d", tCount, 1)})
			}
			req.QuoteAmount, err = strconv.ParseFloat(line[3], 64)
			if err != nil {
				slog.Error("Error parsing imported transaction.", "error", err)
				return c.Render("toastErr", fiber.Map{"Title": "Created", "Msg": fmt.Sprintf("Error on row %d col %d", tCount, 2)})
			}
			req.FeeInBase, err = strconv.ParseFloat(line[4], 64)
			if err != nil {
				slog.Error("Error parsing imported transaction.", "error", err)
				return c.Render("toastErr", fiber.Map{"Title": "Created", "Msg": fmt.Sprintf("Error on row %d col %d", tCount, 3)})
			}
			req.FeeInQuote, err = strconv.ParseFloat(line[5], 64)
			if err != nil {
				slog.Error("Error parsing imported transaction.", "error", err)
				return c.Render("toastErr", fiber.Map{"Title": "Created", "Msg": fmt.Sprintf("Error on row %d col %d", tCount, 4)})
			}
			reqs = append(reqs, req)
			tCount++
		}
		failedTransactions := []string{}
		for i, r := range reqs {
			_, err := tpm.RecordTransaction(r)
			if err != nil {
				slog.Error("Attempt to import transaction failed.", "error", err, "row", i)
				failedTransactions = append(failedTransactions, fmt.Sprint(i))
			}
		}
		tok := int(tCount - len(failedTransactions))
		c.Append("HX-Trigger", "newTransaction")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Extra": "",
			"Msg":   fmt.Sprintf("%d/%d imported.\n%d/%d failed.", tok, tCount, len(failedTransactions), tCount),
		})
	}
}

func newTransaction(tpm managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Recording new transaction")
		payload := struct {
			Base   float64 `form:"base"`
			Quote  float64 `form:"quote"`
			BFee   float64 `form:"bfee"`
			QFee   float64 `form:"qfee"`
			TType  string  `form:"ttype"`
			TTDate string  `form:"tdate"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			fmt.Println(err)
			return err
		}
		tdate, err := time.Parse("2006-01-02", payload.TTDate)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		req := managing.RecordTransactionReq{
			TradingPairID: c.Params("id"),
			Timestamp:     tdate,
			BaseAmount:    payload.Base,
			QuoteAmount:   payload.Quote,
			FeeInBase:     payload.BFee,
			FeeInQuote:    payload.QFee,
			Type:          payload.TType,
		}
		resp, err := tpm.RecordTransaction(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		c.Append("HX-Trigger", "newTransaction")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   resp.Msg,
			"Extra": resp.RecordTime.Format("15h 04m 05s"),
		})

	}
}

func newTransactionForm(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Create Transaction UI requested")
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: true,
			WithTransactions:     false,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("transactionForm", fiber.Map{
			"Pair":  resp.Pair,
			"Today": time.Now().Format("2006-01-02"),
		})
	}
}

func transactionTable(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: false,
			WithTransactions:     true,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		// TODO: Should this happen here? or in the app layer ?
		slices.Reverse(resp.Pair.Transactions)
		slog.Info("Pair Section requested", "Pair", resp.Pair.ID)
		c.Append("HX-Trigger", "refreshTransaction")
		return c.Render("transactionTable", fiber.Map{
			"Today":      time.Now().Format(time.UnixDate),
			"TodayShort": time.Now().Format("02/01/2006"),
			"Pair":       resp.Pair,
		})
	}
}

func transactionExport(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:                 id,
			WithCurrentBasePrice: false,
			WithTransactions:     true,
			WithCalculations:     false,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.SendString(fmt.Sprintf("Error exporting transactions. %s", err))
		}
		file := "TransactionType,Timestamp,BaseAmount,QuoteAmount,FeeInBase,FeeInQuote\n"
		for _, t := range resp.Pair.Transactions {
			file += fmt.Sprintf("%s,%s,%f,%f,%f,%f\n", t.TransactionType, t.Timestamp.Format("2006-01-02 15:04"), t.BaseAmount, t.QuoteAmount, t.FeeInBase, t.FeeInQuote)
		}
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=export_%s_%s.csv", resp.Pair.BaseAsset.Symbol, resp.Pair.QuoteAsset.Symbol))
		c.Set("Content-Type", "application/octet-stream")
		// Send the string as a response
		return c.SendString(file)
	}
}
