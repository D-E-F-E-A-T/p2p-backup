package models

type Ethereum struct {
	Account string `pg:",unique"`
	Address string `pg:",unique"`
}
