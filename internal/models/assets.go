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
	UpdatedAt             time.Time

	Ratios []*Ratio
}

var assetColumns = []string{"asset_id", "api_id", "symbol", "name", "image", "current_price", "high_24h", "low_24h", "market_cap", "market_cap_rank", "fully_diluted_valuation", "total_volume", "circulating_supply", "total_supply", "max_supply", "ath", "atl", "contract", "updated_at"}

func assetFromRows(row durable.Row) (*Asset, error) {
	var a Asset
	err := row.Scan(&a.AssetID, &a.AppID, &a.Symbol, &a.Name, &a.Image, &a.CurrentPrice, &a.High24h, &a.Low24h, &a.MarketCap, &a.MarketCapRank, &a.FullyDilutedValuation, &a.TotalVolume, &a.CirculatingSupply, &a.TotalSupply, &a.MaxSupply, &a.ATH, &a.ATL, &a.Contract, &a.UpdatedAt)
	return &a, err
}

func (a *Asset) values() []any {
	return []any{a.AssetID, a.AppID, a.Symbol, a.Name, a.Image, a.CurrentPrice, a.High24h, a.Low24h, a.MarketCap, a.MarketCapRank, a.FullyDilutedValuation, a.TotalVolume, a.CirculatingSupply, a.TotalSupply, a.MaxSupply, a.ATH, a.ATL, a.Contract, a.UpdatedAt}
}

func UpsertAsset(ctx context.Context, appID, symbol, name, image, price, high, low, marketCap string, marketCapRank int64, valuation, volumn, supply, totalSupply, maxSupply, ath, atl string) (*Asset, error) {
	set := fetchContractMap()
	symbol = strings.ToUpper(symbol)
	contract := set[fmt.Sprintf("%sUSDT", symbol)]
	if contract == "" {
		contract = set[fmt.Sprintf("%sBUSD", symbol)]
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
		UpdatedAt:             time.Now(),
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findAsset(ctx, tx, asset.AssetID)
		if err != nil {
			return err
		}
		if old != nil {
			cols, posits := durable.PrepareColumnsAndExpressions([]string{"symbol", "name", "image", "current_price", "high_24h", "low_24h", "market_cap", "market_cap_rank", "fully_diluted_valuation", "total_volume", "circulating_supply", "total_supply", "max_supply", "ath", "atl", "contract", "updated_at"}, 1)
			values := []any{asset.AssetID, asset.Symbol, asset.Name, asset.Image, asset.CurrentPrice, asset.High24h, asset.Low24h, asset.MarketCap, asset.MarketCapRank, asset.FullyDilutedValuation, asset.TotalVolume, asset.CirculatingSupply, asset.TotalSupply, asset.MaxSupply, asset.ATH, asset.ATL, asset.Contract, asset.UpdatedAt}
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

var contracts = []string{"BTCUSDT", "ETHUSDT", "BCHUSDT", "XRPUSDT", "EOSUSDT", "LTCUSDT", "TRXUSDT", "ETCUSDT", "LINKUSDT", "XLMUSDT", "ADAUSDT", "XMRUSDT", "DASHUSDT", "ZECUSDT", "XTZUSDT", "BNBUSDT", "ATOMUSDT", "ONTUSDT", "IOTAUSDT", "BATUSDT", "VETUSDT", "NEOUSDT", "QTUMUSDT", "IOSTUSDT", "THETAUSDT", "ALGOUSDT", "ZILUSDT", "KNCUSDT", "ZRXUSDT", "COMPUSDT", "OMGUSDT", "DOGEUSDT", "SXPUSDT", "KAVAUSDT", "BANDUSDT", "RLCUSDT", "WAVESUSDT", "MKRUSDT", "SNXUSDT", "DOTUSDT", "DEFIUSDT", "YFIUSDT", "BALUSDT", "CRVUSDT", "TRBUSDT", "RUNEUSDT", "SUSHIUSDT", "EGLDUSDT", "SOLUSDT", "ICXUSDT", "STORJUSDT", "BLZUSDT", "UNIUSDT", "AVAXUSDT", "FTMUSDT", "HNTUSDT", "ENJUSDT", "FLMUSDT", "TOMOUSDT", "RENUSDT", "KSMUSDT", "NEARUSDT", "AAVEUSDT", "FILUSDT", "RSRUSDT", "LRCUSDT", "MATICUSDT", "OCEANUSDT", "BELUSDT", "CTKUSDT", "AXSUSDT", "ALPHAUSDT", "ZENUSDT", "SKLUSDT", "GRTUSDT", "INCHUSDT", "BTCBUSD", "CHZUSDT", "SANDUSDT", "ANKRUSDT", "LITUSDT", "UNFIUSDT", "REEFUSDT", "RVNUSDT", "SFPUSDT", "XEMUSDT", "COTIUSDT", "CHRUSDT", "MANAUSDT", "ALICEUSDT", "HBARUSDT", "ONEUSDT", "LINAUSDT", "STMXUSDT", "DENTUSDT", "CELRUSDT", "HOTUSDT", "MTLUSDT", "OGNUSDT", "NKNUSDT", "DGBUSDT", "SHIBUSDT", "BAKEUSDT", "GTCUSDT", "ETHBUSD", "BNBBUSD", "ADABUSD", "XRPBUSD", "IOTXUSDT", "DOGEBUSD", "AUDIOUSDT", "CUSDT", "MASKUSDT", "ATAUSDT", "SOLBUSD", "DYDXUSDT", "XECUSDT", "GALAUSDT", "CELOUSDT", "ARUSDT", "KLAYUSDT", "ARPAUSDT", "CTSIUSDT", "LPTUSDT", "ENSUSDT", "PEOPLEUSDT", "ANTUSDT", "ROSEUSDT", "DUSKUSDT", "FLOWUSDT", "IMXUSDT", "APIUSDT", "GMTUSDT", "APEUSDT", "WOOUSDT", "JASMYUSDT", "DARUSDT", "GALUSDT", "AVAXBUSD", "NEARBUSD", "GMTBUSD", "APEBUSD", "GALBUSD", "FTMBUSD", "DODOBUSD", "GALABUSD", "TRXBUSD", "LUNCBUSD", "LUNABUSD", "OPUSDT", "DOTBUSD", "TLMBUSD", "ICPBUSD", "WAVESBUSD", "LINKBUSD", "SANDBUSD", "LTCBUSD", "MATICBUSD", "CVXBUSD", "FILBUSD", "SHIBBUSD", "LEVERBUSD", "ETCBUSD", "LDOBUSD", "UNIBUSD", "INJUSDT", "STGUSDT", "FOOTBALLUSDT", "SPELLUSDT", "LUNCUSDT", "LUNAUSDT", "AMBBUSD", "PHBBUSD", "LDOUSDT", "CVXUSDT", "ICPUSDT", "APTUSDT", "QNTUSDT", "APTBUSD", "FETUSDT", "AGIXBUSD", "FXSUSDT", "HOOKUSDT", "MAGICUSDT", "TUSDT", "RNDRUSDT", "HIGHUSDT", "MINAUSDT", "ASTRUSDT", "AGIXUSDT", "PHBUSDT", "GMXUSDT"}

func fetchContractMap() map[string]string {
	m := make(map[string]string, 0)
	for _, c := range contracts {
		m[c] = c
	}
	return m
}
