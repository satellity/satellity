package views

import "satellity/internal/models"

type RatioView struct {
	Category       string `json:"category"`
	Symbol         string `json:"symbol"`
	LongAccount    string `json:"long_account"`
	ShortAccount   string `json:"short_account"`
	LongShortRatio string `json:"long_short_ratio"`
	TimestampAt    int64  `json:"timestamp"`
}

func buildRadio(r *models.Ratio) RatioView {
	return RatioView{
		Category:       r.Category,
		Symbol:         r.Symbol,
		LongAccount:    r.LongAccount,
		ShortAccount:   r.ShortAccount,
		LongShortRatio: r.LongShortRatio,
		TimestampAt:    r.TimestampAt,
	}
}
