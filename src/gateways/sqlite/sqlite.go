package sqlite

import (
	"controtto/src/app/config"
	"controtto/src/domain/pnl"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

const tables string = `
        CREATE TABLE IF NOT EXISTS Pairs (
          ID TEXT PRIMARY KEY,
          BaseAsset TEXT,
          QuoteAsset TEXT
		);
		
        CREATE TABLE IF NOT EXISTS Trades (
          ID INTEGER PRIMARY KEY,
          Timestamp DATETIME,
          BaseAmount REAL,
          QuoteAmount REAL,
          TradeType TEXT,
          PairID TEXT,
          FeeInBase REAL,
          FeeInQuote REAL,
          FOREIGN KEY (PairID) REFERENCES Pairs (ID)
		);

        CREATE TABLE IF NOT EXISTS Asset (
          Symbol TEXT PRIMARY KEY,
          Name TEXT,
          Color TEXT,
          Type TEXT,
          CountryCode TEXT
        );

        -- Forex

        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('JPY', 'Japanese Yen', 'Forex', '#C91400', 'JP');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('GBP', 'British Pound', 'Forex', '#1C7CDD', 'GB');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('CAD', 'Canadian Dollar', 'Forex', '#CA9832', 'CA');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('AUD', 'Australian Dollar', 'Forex', '#029547', 'AU');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('CHF', 'Swiss Franc', 'Forex', '#FF8C00', 'CH');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('ARS', 'Argentine Peso', 'Forex', '#40AAF3', 'AR');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('EUR', 'Euro', 'Forex', '#004C8D', 'EU');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('USD', 'US Dollar', 'Forex', '#809862', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('NZD', 'New Zealand Dollar', 'Forex', '#009F3D', 'NZ');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('HKD', 'Hong Kong Dollar', 'Forex', '#6D4C42', 'HK');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('INR', 'Indian Rupee', 'Forex', '#FF9933', 'IN');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('CNY', 'Chinese Yuan', 'Forex', '#DE2910', 'CN');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('SGD', 'Singapore Dollar', 'Forex', '#FF6F61', 'SG');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('SEK', 'Swedish Krona', 'Forex', '#0076A8', 'SE');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('NOK', 'Norwegian Krone', 'Forex', '#C8102E', 'NO');


        -- Crypto

        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('BTC', 'Bitcoin', 'Crypto', '#F7931A', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('DOT', 'Polkadot', 'Crypto', '#DF0076', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('USDT', 'Tether', 'Crypto', '#009393', 'US');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('ETH', 'Ethereum', 'Crypto', '#c061cb', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('XRP', 'Ripple', 'Crypto', '#005FF9', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('LTC', 'Litecoin', 'Crypto', '#838383', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('ADA', 'Cardano', 'Crypto', '#1C366D', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('DOGE', 'Dogecoin', 'Crypto', '#C3A634', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Type, Color, CountryCode) VALUES ('XMR', 'Monero', 'Crypto', '#FF6600', '-');


        -- Stocks

        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('GOOGL', 'Alphabet', '#4285F4', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('MSFT', 'Microsoft', '#00A3E0', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('TSLA', 'Tesla', '#E82127', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('AAPL', 'Apple', '#5e5c64', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('AMZN', 'Amazon', '#ffa348', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('META', 'Meta Platforms', '#1877F2', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('NFLX', 'Netflix', '#E50914', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('GOOG', 'Alphabet Inc.', '#4285F4', 'Stock', '-');
        INSERT OR IGNORE INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES ('TSLA', 'Tesla', '#E82127', 'Stock', '-');
`

var demo string = `
			-- Trading Pairs
		INSERT OR IGNORE INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES ('BTCUSDT', 'BTC', 'USDT');
		INSERT OR IGNORE INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES ('ETHUSDT', 'ETH', 'USDT');
		INSERT OR IGNORE INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES ('EURUSD', 'EUR', 'USD');
		INSERT OR IGNORE INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES ('GBPJPY', 'GBP', 'JPY');
		INSERT OR IGNORE INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES ('AAPLUSD', 'AAPL', 'USD');
		INSERT OR IGNORE INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES ('TSLAUSD', 'TSLA', 'USD');

			INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:00:00', 1, 86500.00, 'Buy', 'BTCUSDT', 0.0005, 43.25);

		-- ETHUSDT
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:30:00', 1.5, 4800.00, 'Buy', 'ETHUSDT', 0.001, 2.40);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:35:00', 1.0, 4900.00, 'Sell', 'ETHUSDT', 0.0008, 2.45); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:40:00', 1.2, 4950.00, 'Buy', 'ETHUSDT', 0.0009, 2.47);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:45:00', 1.0, 5100.00, 'Sell', 'ETHUSDT', 0.0008, 2.55); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:50:00', 1.8, 5200.00, 'Buy', 'ETHUSDT', 0.0012, 2.60);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T10:55:00', 1.5, 5300.00, 'Sell', 'ETHUSDT', 0.001, 2.65); 

		-- EURUSD
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:00:00', 1000, 1080.00, 'Buy', 'EURUSD', 1.0, 1.08);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:05:00', 1000, 1100.00, 'Sell', 'EURUSD', 1.0, 1.10);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:10:00', 1200, 1120.00, 'Buy', 'EURUSD', 1.2, 1.12);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:15:00', 1100, 1150.00, 'Sell', 'EURUSD', 1.1, 1.15); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:20:00', 1500, 1180.00, 'Buy', 'EURUSD', 1.5, 1.18);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:25:00', 1400, 1200.00, 'Sell', 'EURUSD', 1.4, 1.20); 

		-- GBPJPY
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:30:00', 500, 87000.00, 'Buy', 'GBPJPY', 0.5, 43.5);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:35:00', 500, 88000.00, 'Sell', 'GBPJPY', 0.5, 44.0); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:40:00', 700, 89000.00, 'Buy', 'GBPJPY', 0.7, 44.5);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:45:00', 600, 90000.00, 'Sell', 'GBPJPY', 0.6, 45.0); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:50:00', 1000, 91500.00, 'Buy', 'GBPJPY', 1.0, 45.75);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T11:55:00', 1000, 92000.00, 'Sell', 'GBPJPY', 1.0, 46.0); 

		-- AAPLUSD
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T12:00:00', 20, 3600.00, 'Buy', 'AAPLUSD', 0.02, 3.60);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T12:05:00', 20, 3800.00, 'Sell', 'AAPLUSD', 0.02, 3.80); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T12:10:00', 30, 3850.00, 'Buy', 'AAPLUSD', 0.03, 3.85);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T12:15:00', 25, 3900.00, 'Sell', 'AAPLUSD', 0.025, 3.90); 
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T12:20:00', 50, 4000.00, 'Buy', 'AAPLUSD', 0.05, 4.00);
		INSERT OR IGNORE INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) 
		VALUES ('2025-04-13T12:25:00', 50, 4200.00, 'Sell', 'AAPLUSD', 0.05, 4.20); 
	`

// SQLiteStorage implements the Pairs interface using SQLiteStorage.
type SQLiteStorage struct {
	db *sql.DB
}

// NewStorage creates a new instance of SQLite.
func NewSQLite(cfg *config.Manager) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", cfg.Get().DBPath)
	if err != nil {
		return nil, err
	}
	// Create the Pair and Trade tables if they don't exist.
	qString := tables
	if cfg.Get().LoadSampleData {
		slog.Info("Loading sample data into the database.")
		qString += demo
	}
	_, err = db.Exec(qString)
	if err != nil {
		slog.Error("Error creating tables", "error", err)
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) AddAsset(a pnl.Asset) error {
	_, err := s.db.Exec("INSERT INTO Asset (Symbol, Name, Color, Type, CountryCode) VALUES (?, ?, ?, ?, ?)",
		a.Symbol, a.Name, a.Color, a.Type, a.CountryCode)
	if err != nil {
		slog.Error("Error adding asset.", "error", err)
		return err
	}
	return nil
}

// GetAsset retrieves an asset by symbol from the database.
func (s *SQLiteStorage) GetAsset(symbol string) (*pnl.Asset, error) {
	var asset pnl.Asset
	err := s.db.QueryRow("SELECT Symbol, Name, Color, Type, CountryCode FROM Asset WHERE Symbol = ?", symbol).
		Scan(&asset.Symbol, &asset.Name, &asset.Color, &asset.Type, &asset.CountryCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset with symbol %s not found", symbol)
		}
		return nil, err
	}
	return &asset, nil
}

// ListAssets retrieves a list of all assets from the database.
func (s *SQLiteStorage) ListAssets() ([]pnl.Asset, error) {
	rows, err := s.db.Query("SELECT Symbol, Name, Color, Type, CountryCode FROM Asset")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []pnl.Asset
	for rows.Next() {
		var asset pnl.Asset
		if err := rows.Scan(&asset.Symbol, &asset.Name, &asset.Color, &asset.Type, &asset.CountryCode); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

// Add adds a new trading pair to the database.
func (s *SQLiteStorage) AddPair(tp pnl.Pair) error {
	_, err := s.db.Exec("INSERT INTO Pairs (ID, BaseAsset, QuoteAsset) VALUES (?, ?, ?)", string(tp.ID), tp.BaseAsset.Symbol, tp.QuoteAsset.Symbol)
	if err != nil {
		slog.Error("Error adding trading pair", "error", err)
		return err
	}
	return nil
}

// List returns a list of all trading pairs.
func (s *SQLiteStorage) GetPair(tpid string) (*pnl.Pair, error) {
	var tp pnl.Pair
	var baseSymbol, quoteSymbol string
	err := s.db.QueryRow("SELECT ID, BaseAsset, QuoteAsset FROM Pairs WHERE ID = ?", tpid).Scan(&tp.ID, &baseSymbol, &quoteSymbol)
	if err != nil {
		slog.Error("Error retrieving Trading pair", "ID", tpid, "error", err)
		return nil, err
	}
	base, err := s.GetAsset(baseSymbol)
	if err != nil {
		return nil, err
	}
	quote, err := s.GetAsset(quoteSymbol)
	if err != nil {
		return nil, err
	}
	tp.BaseAsset = *base
	tp.QuoteAsset = *quote
	return &tp, nil
}

// List returns a list of all trading pairs.
func (s *SQLiteStorage) ListPairs() ([]pnl.Pair, error) {
	rows, err := s.db.Query("SELECT ID, BaseAsset, QuoteAsset FROM Pairs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tradingPairs []pnl.Pair
	for rows.Next() {
		var tp pnl.Pair
		var baseSymbol, quoteSymbol string
		if err := rows.Scan(&tp.ID, &baseSymbol, &quoteSymbol); err != nil {
			return nil, err
		}
		base, err := s.GetAsset(baseSymbol)
		if err != nil {
			return nil, err
		}
		quote, err := s.GetAsset(quoteSymbol)
		if err != nil {
			return nil, err
		}
		tp.BaseAsset = *base
		tp.QuoteAsset = *quote
		tradingPairs = append(tradingPairs, tp)
	}

	return tradingPairs, nil
}

// RecordTrade records a new trade for a trading pair.
func (s *SQLiteStorage) RecordTrade(t pnl.Trade, tpid pnl.PairID) error {
	_, err := s.db.Exec("INSERT INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, PairID, FeeInBase, FeeInQuote) VALUES (?, ?, ?, ?, ?, ?, ?)",
		t.Timestamp, t.BaseAmount, t.QuoteAmount, t.TradeType, string(tpid), t.FeeInBase, t.FeeInQuote)
	return err
}

// ListTrades returns a list of trades for a given trading pair ID.
func (s *SQLiteStorage) ListTrades(tpid string) ([]pnl.Trade, error) {
	rows, err := s.db.Query("SELECT ID, Timestamp, BaseAmount, QuoteAmount, TradeType, FeeInBase, FeeInQuote FROM Trades WHERE PairID = ?", tpid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []pnl.Trade
	for rows.Next() {
		var t pnl.Trade
		if err := rows.Scan(&t.ID, &t.Timestamp, &t.BaseAmount, &t.QuoteAmount, &t.TradeType, &t.FeeInBase, &t.FeeInQuote); err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}

	return trades, nil
}

// DeletePair deletes a trading pair by its ID.
func (s *SQLiteStorage) DeletePair(tpid string) error {
	_, err := s.db.Exec("DELETE FROM Pairs WHERE ID = ?", tpid)
	if err != nil {
		slog.Error("Error deleting trading pair", "error", err)
		return err
	}
	return nil
}

// DeleteTrade deletes a trade by its ID.
func (s *SQLiteStorage) DeleteTrade(tpid string) error {
	_, err := s.db.Exec("DELETE FROM Trades WHERE ID = ?", tpid)
	if err != nil {
		slog.Error("Error deleting trade", "error", err)
		return err
	}
	return nil
}
