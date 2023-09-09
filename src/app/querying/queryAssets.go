package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type QueryAssetReq struct {
	Symbol string
}
type QueryAssetResp struct {
	Asset pnl.Asset
}

type QueryAssetsReq struct{}
type QueryAssetsResp struct {
	Assets []pnl.Asset
}

type AssetsQuerier struct {
	assets pnl.Assets
}

func NewAssetQuerier(a pnl.Assets) *AssetsQuerier {
	return &AssetsQuerier{a}
}

func (aq *AssetsQuerier) ListAssets(req QueryAssetsReq) (*QueryAssetsResp, error) {
	var err error
	assets, err := aq.assets.ListAssets()
	if err != nil {
		slog.Error("Error listing assets from DB", "error", err)
		return nil, err
	}
	for _, a := range assets {
		if _, err := a.Validate(); err != nil {
			slog.Error("Invalid asset "+a.Symbol, "error", err)
		}
	}
	resp := QueryAssetsResp{
		Assets: assets,
	}
	return &resp, nil
}

func (aq *AssetsQuerier) GetAsset(req QueryAssetReq) (*QueryAssetResp, error) {
	var err error
	asset, err := aq.assets.GetAsset(req.Symbol)
	if err != nil {
		slog.Error("Error retrieving Asset from DB", "error", err)
		return nil, err
	}
	if _, err := asset.Validate(); err != nil {
		slog.Error("Invalid asset retrieved from DB", "error", err)
		return nil, err
	}
	resp := QueryAssetResp{
		Asset: *asset,
	}
	return &resp, nil
}
