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
	Name        string // Optional
	CountryCode string // Optional
}

// Assets creator
type AssetsCreator struct {
	assets pnl.Assets
}

func NewAssetCreator(a pnl.Assets) *AssetsCreator {
	return &AssetsCreator{a}
}

func (ac *AssetsCreator) Create(req CreateAssetReq) (*CreateAssetResp, error) {
	var err error
	asset, err := pnl.NewAsset(req.Symbol, req.Color, req.Name, req.CountryCode)
	if err != nil {
		slog.Error("Creating new asset", "Name", req.Name, "Countr code", req.CountryCode, "error", err)
		return nil, err
	}
	err = ac.assets.AddAsset(*asset)
	// TODO: Validate already exists asset
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
