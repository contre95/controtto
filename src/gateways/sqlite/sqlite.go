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
		
        CREATE TABLE IF NOT EXISTS Transactions (
          ID INTEGER PRIMARY KEY,
          Timestamp DATETIME,
          BaseAmount REAL,
          QuoteAmount REAL,
          TransactionType TEXT,
          TradingPairID TEXT,
          FOREIGN KEY (TradingPairID) REFERENCES TradingPairs (ID)
		);

        CREATE TABLE IF NOT EXISTS Asset (
          Symbol TEXT PRIMARY KEY,
          Name TEXT,
          Color TEXT,
          CountryCode TEXT
        );

        INSERT OR IGNORE INTO Asset (Symbol, Name, CountryCode, Color) VALUES ('BTC', 'Bitcoin', '-', '#F7931A');
        INSERT OR IGNORE INTO Asset (Symbol, Name, CountryCode, Color) VALUES ('DOT', 'Polkadot', '-', '#DF0076');
        INSERT OR IGNORE INTO Asset (Symbol, Name, CountryCode, Color) VALUES ('ARS', 'Peso', 'AR', '#40AAF3');
        INSERT OR IGNORE INTO Asset (Symbol, Name, CountryCode, Color) VALUES ('EUR', 'Euro', 'EU', '#004C8D');
        INSERT OR IGNORE INTO Asset (Symbol, Name, CountryCode, Color) VALUES ('USDT', 'Tether', 'US', '#009393');
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
	// Create the TradingPair and Transaction tables if they don't exist.
	_, err = db.Exec(tables)
	if err != nil {
		slog.Error("Error creating tables", "error", err)
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) AddAsset(a pnl.Asset) error {
	_, err := s.db.Exec("INSERT INTO Asset (Symbol, Name, Color, CountryCode) VALUES (?, ?, ?, ?)",
		a.Symbol, a.Name, a.Color, a.CountryCode)
	if err != nil {
		slog.Error("Error adding asset.", "error", err)
		return err
	}
	return nil
}

// GetAsset retrieves an asset by symbol from the database.
func (s *SQLiteStorage) GetAsset(symbol string) (*pnl.Asset, error) {
	var asset pnl.Asset
	err := s.db.QueryRow("SELECT Symbol, Name, Color, CountryCode FROM Asset WHERE Symbol = ?", symbol).
		Scan(&asset.Symbol, &asset.Name, &asset.Color, &asset.CountryCode)
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
	rows, err := s.db.Query("SELECT Symbol, Name, Color, CountryCode FROM Asset")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var assets []pnl.Asset
	for rows.Next() {
		var asset pnl.Asset
		if err := rows.Scan(&asset.Symbol, &asset.Name, &asset.Color, &asset.CountryCode); err != nil {
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

// RecordTransaction records a new transaction for a trading pair.
func (s *SQLiteStorage) RecordTransaction(t pnl.Transaction, tpid pnl.TradingPairID) error {
	_, err := s.db.Exec("INSERT INTO Transactions (Timestamp, BaseAmount, QuoteAmount, TransactionType, TradingPairID) VALUES (?, ?, ?, ?, ?)",
		t.Timestamp, t.BaseAmount, t.QuoteAmount, t.TransactionType, string(tpid))
	return err
}

// ListTransactions returns a list of transactions for a given trading pair ID.
func (s *SQLiteStorage) ListTransactions(tpid string) ([]pnl.Transaction, error) {
	rows, err := s.db.Query("SELECT ID, Timestamp, BaseAmount, QuoteAmount, TransactionType FROM Transactions WHERE TradingPairID = ?", tpid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []pnl.Transaction
	for rows.Next() {
		var t pnl.Transaction
		if err := rows.Scan(&t.ID, &t.Timestamp, &t.BaseAmount, &t.QuoteAmount, &t.TransactionType); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
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

// DeleteTransaction deletes a transaction by its ID.
func (s *SQLiteStorage) DeleteTransaction(tpid string) error {
	_, err := s.db.Exec("DELETE FROM Transactions WHERE ID = ?", tpid)
	if err != nil {
		slog.Error("Error deleting transaction", "error", err)
		return err
	}
	return nil
}
