package sqlite

import (
	"controtto/src/domain/pnl"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

const tables string = `
        CREATE TABLE IF NOT EXISTS TradingPairs (
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
          TradingPairID TEXT,
          FOREIGN KEY (TradingPairID) REFERENCES TradingPairs (ID)
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

// SQLiteStorage implements the TradingPairs interface using SQLiteStorage.
type SQLiteStorage struct {
	db *sql.DB
}

// NewStorage creates a new instance of SQLite.
func NewSQLite(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Create the TradingPair and Trade tables if they don't exist.
	_, err = db.Exec(tables)
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
func (s *SQLiteStorage) AddTradingPair(tp pnl.TradingPair) error {
	_, err := s.db.Exec("INSERT INTO TradingPairs (ID, BaseAsset, QuoteAsset) VALUES (?, ?, ?)", string(tp.ID), tp.BaseAsset.Symbol, tp.QuoteAsset.Symbol)
	if err != nil {
		slog.Error("Error adding trading pair", "error", err)
		return err
	}
	return nil
}

// List returns a list of all trading pairs.
func (s *SQLiteStorage) GetTradingPair(tpid string) (*pnl.TradingPair, error) {
	var tp pnl.TradingPair
	var baseSymbol, quoteSymbol string
	err := s.db.QueryRow("SELECT ID, BaseAsset, QuoteAsset FROM TradingPairs WHERE ID = ?", tpid).Scan(&tp.ID, &baseSymbol, &quoteSymbol)
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
func (s *SQLiteStorage) ListTradingPairs() ([]pnl.TradingPair, error) {
	rows, err := s.db.Query("SELECT ID, BaseAsset, QuoteAsset FROM TradingPairs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tradingPairs []pnl.TradingPair
	for rows.Next() {
		var tp pnl.TradingPair
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
func (s *SQLiteStorage) RecordTrade(t pnl.Trade, tpid pnl.TradingPairID) error {
	_, err := s.db.Exec("INSERT INTO Trades (Timestamp, BaseAmount, QuoteAmount, TradeType, TradingPairID) VALUES (?, ?, ?, ?, ?)",
		t.Timestamp, t.BaseAmount, t.QuoteAmount, t.TradeType, string(tpid))
	return err
}

// ListTrades returns a list of trades for a given trading pair ID.
func (s *SQLiteStorage) ListTrades(tpid string) ([]pnl.Trade, error) {
	rows, err := s.db.Query("SELECT ID, Timestamp, BaseAmount, QuoteAmount, TradeType FROM Trades WHERE TradingPairID = ?", tpid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []pnl.Trade
	for rows.Next() {
		var t pnl.Trade
		if err := rows.Scan(&t.ID, &t.Timestamp, &t.BaseAmount, &t.QuoteAmount, &t.TradeType); err != nil {
			return nil, err
		}
		trades = append(trades, t)
	}

	return trades, nil
}

// DeleteTradingPair deletes a trading pair by its ID.
func (s *SQLiteStorage) DeleteTradingPair(tpid string) error {
	_, err := s.db.Exec("DELETE FROM TradingPairs WHERE ID = ?", tpid)
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
