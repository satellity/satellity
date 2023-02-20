package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"satellity/internal/configs"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/session"
	"time"
)

type AssetImpl struct{}

func (impl *AssetImpl) Run(db *durable.Database) {
	log.Println("Asset service started at:", time.Now())
	ctx := session.WithDatabase(context.Background(), db)
	for {
		err := impl.loopFetchAssets(ctx)
		if err != nil {
			log.Printf("services.loopFetchAssets error %s \n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		time.Sleep(10 * time.Second)
	}
}

type Asset struct {
	AppID                 string      `json:"id"`
	Symbol                string      `json:"symbol"`
	Name                  string      `json:"name"`
	Image                 string      `json:"image"`
	CurrentPrice          json.Number `json:"current_price"`
	High24h               json.Number `json:"high_24h"`
	Low24h                json.Number `json:"low_24h"`
	MarketCap             json.Number `json:"market_cap"`
	MarketCapRank         int64       `json:"market_cap_rank"`
	FullyDilutedValuation json.Number `json:"fully_diluted_valuation"`
	TotalVolume           json.Number `json:"total_volume"`
	CirculatingSupply     json.Number `json:"circulating_supply"`
	TotalSupply           json.Number `json:"total_supply"`
	MaxSupply             json.Number `json:"max_supply"`
	ATH                   json.Number `json:"ath"`
	ATL                   json.Number `json:"atl"`
}

func (impl *AssetImpl) loopFetchAssets(ctx context.Context) error {
	resp, err := http.Get(fmt.Sprintf("https://pro-api.coingecko.com/api/v3/coins/markets?vs_currency=usd&per_page=250&x_cg_pro_api_key=%s", configs.AppConfig.Cgc))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var assets []*Asset
	err = json.NewDecoder(resp.Body).Decode(&assets)
	if err != nil {
		return err
	}

	for _, a := range assets {
		_, err = models.UpsertAsset(ctx, a.AppID, a.Symbol, a.Name, a.Image, a.CurrentPrice.String(), a.High24h.String(), a.Low24h.String(), a.MarketCap.String(), a.MarketCapRank, a.FullyDilutedValuation.String(), a.TotalVolume.String(), a.CirculatingSupply.String(), a.TotalSupply.String(), a.MaxSupply.String(), a.ATH.String(), a.ATL.String())
		if err != nil {
			log.Printf("services.loopFetchAssets error %s \n", err)
			time.Sleep(time.Second)
			continue
		}
	}
	return nil
}
