package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type AssetsQuerier struct {
	assets pnl.Assets
}

func NewAssetQuerier(a pnl.Assets) *AssetsQuerier {
	return &AssetsQuerier{a}
}

type ListAssetsReq struct{}

type ListAssetsResp struct {
	Assets []pnl.Asset
}

func (aq *AssetsQuerier) ListAssets(req ListAssetsReq) (*ListAssetsResp, error) {
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
	resp := ListAssetsResp{
		Assets: assets,
	}
	return &resp, nil
}

type QueryAssetReq struct {
	Symbol string
}

type QueryAssetResp struct {
	Asset pnl.Asset
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
