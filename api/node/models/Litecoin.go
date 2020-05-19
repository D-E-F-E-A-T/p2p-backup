package models

type Litecoin struct {
	Account string `pg:",unique:acc"`
	Address string `pg:",unique:acc"`
}
