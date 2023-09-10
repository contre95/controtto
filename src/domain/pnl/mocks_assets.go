package pnl

// MockAssets is a mock implementation of the Assets interface for testing purposes.
type MockAssets struct {
	ListAssetsResponse func() ([]Asset, error)
	GetAssetResponse   func(symbol string) (*Asset, error)
}

// Asset represents individual assets like BTC, USD, EUR, etc. the Symbol property uniquely identifies an asset.
// AddAsset is a method of the Assets interface that adds an asset to the repository.
func (m MockAssets) AddAsset(a Asset) error {
	// Implement your desired behavior here.
	return nil
}

// ListAssets is a method of the Assets interface that lists all assets from the repository.
func (m MockAssets) ListAssets() ([]Asset, error) {
	if m.ListAssetsResponse != nil {
		return m.ListAssetsResponse()
	}
	panic("ListAssetsResponse not implemented in mock")
}

// GetAsset is a method of the Assets interface that retrieves an asset from the repository by its symbol.
func (m MockAssets) GetAsset(symbol string) (*Asset, error) {
	if m.GetAssetResponse != nil {
		return m.GetAssetResponse(symbol)
	}
	panic("GetAssetResponse not implemented in mock")
}
