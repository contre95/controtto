package presenters

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"
	"encoding/csv"
	"io"
	"strconv"

	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
)

func dashboardSection(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Dashboard requested")
		req := querying.ListTradingPairsReq{}
		resp, err := aq.ListTradingPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("dashboardSection", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}

}

func pairsSection(c *fiber.Ctx) error {
	slog.Info("Pairs Section")
	return c.Render("pairsSection", fiber.Map{})
}

func pairSection(c *fiber.Ctx) error {
	slog.Info("Pair Section")
	return c.Render("pairSection", fiber.Map{
		"PairID": c.Params("id"),
	})
}

// Home hanlder reders the homescreen
func Home(c *fiber.Ctx) error {
	slog.Info("HOME")
	// render index template
	return c.Render("main", fiber.Map{})
}

func checkPrice(mq querying.MarketsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		base := c.Query("base")
		quote := c.Query("quote")
		req := querying.QueryMarketReq{
			AssetSymbolA: base,
			AssetSymbolB: quote,
		}
		resp, err := mq.GetMarketPrice(req)
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("%f", resp.Price))
	}
}

func newPairForm(aq querying.AssetsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := querying.QueryAssetsReq{}
		resp, err := aq.ListAssets(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("newPairForm", fiber.Map{
			"Title":  "New Pair",
			"Assets": resp.Assets,
		})
	}
}

func newAssetForm(c *fiber.Ctx) error {
	slog.Info("Create Asset UI requested")
	return c.Render("newAssetForm", fiber.Map{
		"Title": "New Asset",
	})
}

// Tables handler that renderizer the tables view and returns it to the client
func pairsTable(aq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Pairs table requested")
		req := querying.ListTradingPairsReq{}
		resp, err := aq.ListTradingPairs(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("pairsTable", fiber.Map{
			"Title": "Trading Pairs",
			"Pairs": resp.Pairs,
		})
	}
}

func avgBuyPrice(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

// marketPrice handler that reutnr a string witht the total amount of money.
func marketPrice(aq querying.MarketsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		id := c.Params("id")
		slog.Info("Requesting market price", "market", id)
		req := querying.QueryMarketReq{
			AssetSymbolA: id,
		}
		resp, err := aq.GetMarketPrice(req)
		if err != nil {
			return c.SendString("--")
		}
		// render index template
		return c.SendString(fmt.Sprintf("%.2f", resp.Price))
	}
}

// TotalUsers handler that reutnr a string witht the total amount of users.
func TotalUsers(c *fiber.Ctx) error {
	slog.Info("USERS")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := rand.Intn(101)
	return c.SendString(fmt.Sprintf("%d", randomNumber))
}

// Tables handler that renderizer the tables view and returns it to the client
func Tables(c *fiber.Ctx) error {
	slog.Info("TABLES")
	return c.Render("tables", fiber.Map{
		"Title":    "Hello, World!",
		"Selected": "tables",
	})
}

func NewPair(c *fiber.Ctx) error {
	slog.Info("Creating new asset")
	slog.Info(string(c.Body()))
	return c.Send(c.Body())
}

func Assets(c *fiber.Ctx) error {
	slog.Info("Assets requested")
	return nil
}

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
		return c.Render("toastErr", fiber.Map{
			"Title": "Error",
			"Msg":   "Deleted",
		})
	}
}

func deleteTradingPair(tpm managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := managing.DeleteTradingPairReq{
			ID: c.Params("id"),
		}
		slog.Info("Delete", "id", req.ID)
		resp, err := tpm.DeleteTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slog.Info("Trading pair deleted", "trading-pair", resp.ID)
		c.Append("HX-Trigger", "newPair")
		return c.Render("toastErr", fiber.Map{
			"Title": "Error",
			"Msg":   "Deleted",
		})
	}
}

func newTradingPair(tpc managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Creating new pair")
		payload := struct {
			Quote string `form:"quote"`
			Base  string `form:"base"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			fmt.Println(err)
			return err
		}
		req := managing.CreateTradingPairReq{
			BaseAsset:  payload.Base,
			QuoteAsset: payload.Quote,
		}
		slog.Info("Creating", "req", req)
		resp, err := tpc.Create(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		c.Append("HX-Trigger", "newPair")
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   resp.Msg,
		})
	}
}

func newAsset(ac managing.AssetsCreator) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Creating new asset")
		payload := struct {
			Symbol string `form:"symbol"`
			Name   string `form:"name"`
			Color  string `form:"color"`
		}{}
		if err := c.BodyParser(&payload); err != nil {
			fmt.Println(err)
			return err
		}
		req := managing.CreateAssetReq{
			Symbol:      payload.Symbol,
			Color:       payload.Color,
			Name:        payload.Name,
			CountryCode: "-",
		}
		resp, err := ac.Create(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		return c.Render("toastOk", fiber.Map{
			"Title": "Created",
			"Msg":   resp.Msg,
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
		return c.Render("transactionsResponse", fiber.Map{
			"Title":   "Info",
			"TErr":    len(failedTransactions),
			"TOk":     tok,
			"TCount":  tCount,
			"TFailed": failedTransactions,
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

func pairCards(tpq querying.TradingPairsQuerier) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := querying.GetTradingPairReq{
			TPID:             id,
			WithCalculations: true,
		}
		resp, err := tpq.GetTradingPair(req)
		if err != nil {
			return c.Render("toastErr", fiber.Map{
				"Title": "Error",
				"Msg":   err,
			})
		}
		slices.Reverse(resp.Pair.Transactions)
		return c.Render("pairCards", fiber.Map{
			"Today": time.Now().Format("Mon Jan 02 15:04 2006"),
			"Pair":  resp.Pair,
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
