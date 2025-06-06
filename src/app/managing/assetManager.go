package managing

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type CreateAssetResp struct {
	Symbol string
	Msg    string
}

type CreateAssetReq struct {
	Symbol      string
	Color       string
	Type        string
	Name        string // Optional
	CountryCode string // Optional
}

// Assets creator
type AssetsManager struct {
	assets pnl.Assets
}

func NewAssetCreator(a pnl.Assets) *AssetsManager {
	return &AssetsManager{a}
}

func (ac *AssetsManager) Create(req CreateAssetReq) (*CreateAssetResp, error) {
	var err error
	asset, err := pnl.NewAsset(req.Symbol, req.Color, req.Name, req.CountryCode, req.Type)
	if err != nil {
		slog.Error("Creating new asset", "Name", req.Name, "CountryCode", req.CountryCode, "error", err)
		return nil, err
	}
	err = ac.assets.AddAsset(*asset)
	// TODO: Validate already exists asset ?
	if err != nil {
		slog.Error("Creating saving asset", "asset", asset, "error", err)
		return nil, err
	}
	slog.Info("Asset created", "Asset", asset)
	return &CreateAssetResp{
		Symbol: asset.Symbol,
		Msg:    "New asset created",
	}, nil
}
