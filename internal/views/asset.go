package views

import (
	"net/http"
	"satellity/internal/models"
)

type AssetView struct {
	Symbol                string      `json:"symbol"`
	Name                  string      `json:"name"`
	Image                 string      `json:"image"`
	CurrentPrice          string      `json:"current_price"`
	High24h               string      `json:"high_24h"`
	Low24h                string      `json:"low_24h"`
	MarketCap             string      `json:"market_cap"`
	MarketCapRank         int64       `json:"market_cap_rank"`
	FullyDilutedValuation string      `json:"fully_diluted_valuation"`
	TotalVolume           string      `json:"total_volume"`
	CirculatingSupply     string      `json:"circulating_supply"`
	TotalSupply           string      `json:"total_supply"`
	MaxSupply             string      `json:"max_supply"`
	ATH                   string      `json:"ath"`
	ATL                   string      `json:"atl"`
	Contract              string      `json:"contract"`
	GlobalRatio           string      `json:"global_ratio"`
	FundingRate           string      `json:"funding_rate"`
	Ratios                []RatioView `json:"ratios"`
}

func buildAsset(a *models.Asset) AssetView {
	view := AssetView{
		Symbol:                a.Symbol,
		Name:                  a.Name,
		Image:                 a.Image,
		CurrentPrice:          a.CurrentPrice,
		High24h:               a.High24h,
		Low24h:                a.Low24h,
		MarketCap:             a.MarketCap,
		MarketCapRank:         a.MarketCapRank,
		FullyDilutedValuation: a.FullyDilutedValuation,
		TotalVolume:           a.TotalVolume,
		CirculatingSupply:     a.CirculatingSupply,
		TotalSupply:           a.TotalSupply,
		MaxSupply:             a.MaxSupply,
		ATH:                   a.ATH,
		ATL:                   a.ATL,
		Contract:              a.Contract,
		FundingRate:           a.FundingRate,
		GlobalRatio:           "0",
	}
	view.Ratios = make([]RatioView, len(a.Ratios))
	for i, r := range a.Ratios {
		if r.Category == models.GlobalLongShortAccountRatio {
			view.GlobalRatio = r.LongShortRatio
		}
		view.Ratios[i] = buildRadio(r)
	}
	return view
}

func RenderAssets(w http.ResponseWriter, r *http.Request, assets []*models.Asset) {
	views := make([]AssetView, len(assets))
	for i, a := range assets {
		views[i] = buildAsset(a)
	}
	RenderResponse(w, r, views)
}
