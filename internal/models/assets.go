package models

import (
	"context"
	"fmt"
	"log"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Asset struct {
	AssetID               string
	AppID                 string
	Symbol                string
	Name                  string
	Image                 string
	CurrentPrice          string
	High24h               string
	Low24h                string
	MarketCap             string
	MarketCapRank         int64
	FullyDilutedValuation string
	TotalVolume           string
	CirculatingSupply     string
	TotalSupply           string
	MaxSupply             string
	ATH                   string
	ATL                   string
	Contract              string
	FundingRate           string
	UpdatedAt             time.Time

	Ratios []*Ratio
}

var assetColumns = []string{"asset_id", "api_id", "symbol", "name", "image", "current_price", "high_24h", "low_24h", "market_cap", "market_cap_rank", "fully_diluted_valuation", "total_volume", "circulating_supply", "total_supply", "max_supply", "ath", "atl", "contract", "funding_rate", "updated_at"}

func assetFromRows(row durable.Row) (*Asset, error) {
	var a Asset
	err := row.Scan(&a.AssetID, &a.AppID, &a.Symbol, &a.Name, &a.Image, &a.CurrentPrice, &a.High24h, &a.Low24h, &a.MarketCap, &a.MarketCapRank, &a.FullyDilutedValuation, &a.TotalVolume, &a.CirculatingSupply, &a.TotalSupply, &a.MaxSupply, &a.ATH, &a.ATL, &a.Contract, &a.FundingRate, &a.UpdatedAt)
	return &a, err
}

func (a *Asset) values() []any {
	return []any{a.AssetID, a.AppID, a.Symbol, a.Name, a.Image, a.CurrentPrice, a.High24h, a.Low24h, a.MarketCap, a.MarketCapRank, a.FullyDilutedValuation, a.TotalVolume, a.CirculatingSupply, a.TotalSupply, a.MaxSupply, a.ATH, a.ATL, a.Contract, a.FundingRate, a.UpdatedAt}
}

func UpsertAsset(ctx context.Context, appID, symbol, name, image, price, high, low, marketCap string, marketCapRank int64, valuation, volumn, supply, totalSupply, maxSupply, ath, atl string, indexes map[string]string) (*Asset, error) {
	symbol = strings.ToUpper(symbol)
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
		return nil, nil
	}
	asset := Asset{
		AssetID:               generateUniqueID("ASSETS:", appID),
		AppID:                 appID,
		Symbol:                symbol,
		Name:                  name,
		Image:                 image,
		CurrentPrice:          price,
		High24h:               high,
		Low24h:                low,
		MarketCap:             marketCap,
		MarketCapRank:         marketCapRank,
		FullyDilutedValuation: valuation,
		TotalVolume:           volumn,
		CirculatingSupply:     supply,
		TotalSupply:           totalSupply,
		MaxSupply:             maxSupply,
		ATH:                   ath,
		ATL:                   atl,
		Contract:              contract,
		FundingRate:           indexes[contract],
		UpdatedAt:             time.Now(),
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findAsset(ctx, tx, asset.AssetID)
		if err != nil {
			return err
		}
		if old != nil {
			cols, posits := durable.PrepareColumnsAndExpressions([]string{"symbol", "name", "image", "current_price", "high_24h", "low_24h", "market_cap", "market_cap_rank", "fully_diluted_valuation", "total_volume", "circulating_supply", "total_supply", "max_supply", "ath", "atl", "contract", "funding_rate", "updated_at"}, 1)
			values := []any{asset.AssetID, asset.Symbol, asset.Name, asset.Image, asset.CurrentPrice, asset.High24h, asset.Low24h, asset.MarketCap, asset.MarketCapRank, asset.FullyDilutedValuation, asset.TotalVolume, asset.CirculatingSupply, asset.TotalSupply, asset.MaxSupply, asset.ATH, asset.ATL, asset.Contract, asset.FundingRate, asset.UpdatedAt}
			_, err := tx.Exec(ctx, fmt.Sprintf("UPDATE assets SET (%s)=(%s) WHERE asset_id=$1", cols, posits), values...)
			return err
		}
		rows := [][]interface{}{
			asset.values(),
		}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"assets"}, assetColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		log.Println("session.TransactionError:", err)
		return nil, session.TransactionError(ctx, err)
	}
	return &asset, nil
}

func findAsset(ctx context.Context, tx pgx.Tx, id string) (*Asset, error) {
	if uuid.FromStringOrNil(id).String() != id {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM assets WHERE asset_id=$1", strings.Join(assetColumns, ",")), id)
	a, err := assetFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func ReadAllAssets(ctx context.Context) ([]*Asset, error) {
	var assets []*Asset
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		assets, err = readAssets(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return assets, nil
}

func ReadAssetsWithRatios(ctx context.Context) ([]*Asset, error) {
	assets, err := ReadAllAssets(ctx)
	if err != nil {
		return nil, err
	}
	return FillAssetWithRatios(ctx, assets)
}

func readAssets(ctx context.Context, tx pgx.Tx) ([]*Asset, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM assets LIMIT 1000", strings.Join(assetColumns, ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*Asset
	for rows.Next() {
		asset, err := assetFromRows(rows)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, rows.Err()
}

func (asset *Asset) Delete(ctx context.Context) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "DELETE FROM assets WHERE asset_id=$1", asset.AssetID)
		return err
	})
	return err
}
