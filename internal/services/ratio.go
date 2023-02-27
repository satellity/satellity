package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/session"
	"strings"
	"time"
)

type ratioImpl struct{}

func (impl *ratioImpl) Run(db *durable.Database) {
	log.Println("Ratio service started at:", time.Now())
	ctx := session.WithDatabase(context.Background(), db)
	for {
		err := impl.loopFetchRatios(ctx)
		if err != nil {
			log.Printf("services.loopFetchRatios error %s \n", err)
			time.Sleep(10 * time.Second)
			continue
		}
		time.Sleep(20 * time.Second)
	}
}

type Ratio struct {
	Symbol         string      `json:"symbol"`
	LongShortRatio string      `json:"longShortRatio"`
	LongAccount    string      `json:"longAccount"`
	ShortAccount   string      `json:"shortAccount"`
	Timestamp      json.Number `json:"timestamp"`
}

func (impl *ratioImpl) loopFetchRatios(ctx context.Context) error {
	indexes, err := fetchPremiumIndex(ctx)
	if err != nil {
		return err
	}
	assets, err := models.ReadAllAssets(ctx)
	if err != nil {
		return err
	}
	set := map[string]string{
		models.TopLongShortAccountRatio:    "https://fapi.binance.com/futures/data/topLongShortAccountRatio?symbol=%s&period=5m&limit=3",
		models.TopLongShortPositionRatio:   "https://fapi.binance.com/futures/data/topLongShortPositionRatio?symbol=%s&period=5m&limit=3",
		models.GlobalLongShortAccountRatio: "https://fapi.binance.com/futures/data/globalLongShortAccountRatio?symbol=%s&period=5m&limit=3",
	}
	for _, asset := range assets {
		symbol := strings.ToUpper(asset.Symbol)
		contract := indexes[fmt.Sprintf("%sUSDT", symbol)]
		if contract == "" {
			contract = indexes[fmt.Sprintf("%sBUSD", symbol)]
		}
		if contract == "" {
			contract = indexes[fmt.Sprintf("1000%sUSDT", symbol)]
		}
		if contract == "" {
			contract = indexes[fmt.Sprintf("1000%sBUSD", symbol)]
		}
		if contract == "" {
			contract = indexes[fmt.Sprintf("%s2USDT", symbol)]
		}
		if contract == "" {
			contract = indexes[fmt.Sprintf("%s2BUSD", symbol)]
		}
		if contract == "" {
			asset.Delete(ctx)
			continue
		}
		for k, v := range set {
			impl.fetchRatio(ctx, k, fmt.Sprintf(v, asset.Contract))
		}
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func (impl *ratioImpl) fetchRatio(ctx context.Context, category, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ratios []*Ratio
	err = json.NewDecoder(resp.Body).Decode(&ratios)
	if err != nil {
		return err
	}
	for _, r := range ratios {
		at, err := r.Timestamp.Int64()
		if err != nil {
			log.Printf("services.fetchRatio error %v \n", err)
			continue
		}
		_, err = models.UpsertRatio(ctx, category, r.Symbol, models.RatioPeriod5M, r.LongAccount, r.ShortAccount, r.LongShortRatio, at)
		if err != nil {
			log.Printf("services.UpsertRatio error %v \n", err)
			continue
		}
	}
	return nil
}
