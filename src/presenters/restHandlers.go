package presenters

import (
	"controtto/src/app/managing"
	"controtto/src/app/querying"

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
			return c.Render("toastCustom", fiber.Map{
				"ToastColor": "#6F6F6F",
				"Msg":        "Price not found.",
			})
		}
		return c.Render("toastCustom", fiber.Map{
			"ToastColor": resp.ProviderColor,
			"Msg":        resp.ProviderName,
		})
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

func newTransaction(tpm managing.TradingPairsManager) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		slog.Info("Recording new transaction")
		payload := struct {
			Base   float64 `form:"base"`
			Quote  float64 `form:"quote"`
			WFee   float64 `form:"wfee"`
			TFee   float64 `form:"tfee"`
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
			TradingFee:    payload.WFee,
			WithdrawalFee: payload.WFee,
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
			"Time":  resp.RecordTime.Format("15h 04m 05s"),
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
		// TODO: Should this happen here? or in the app layer ?
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
