package querying

import (
	"controtto/src/domain/pnl"
	"log/slog"
)

type QueryAssetReq struct{}
type QueryAssetResp struct {
	Assets []pnl.Asset
}

type AssetsQuerier struct {
	assets pnl.Assets
}

func NewAssetQuerier(a pnl.Assets) *AssetsQuerier {
	return &AssetsQuerier{a}
}

func (aq *AssetsQuerier) ListAssets(req QueryAssetReq) (*QueryAssetResp, error) {
	var err error
	assets, err := aq.assets.ListAssets()
	if err != nil {
		slog.Error("Error listing assets from DB", "error", err)
		return nil, err
	}
	resp := QueryAssetResp{
		Assets: assets,
	}
	return &resp, nil
}
