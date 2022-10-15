package models

type Faucet struct {
	FaucetID string
	Name     string
	Chain    string
	Currency string
	Faucets  []string
}
