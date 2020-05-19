package models

type BitcoinCash struct {
	Account string `pg:",unique:acc"`
	Address string `pg:",unique:acc"`
}
