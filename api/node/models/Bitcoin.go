package models

type Bitcoin struct {
	Account string `pg:",unique:acc"`
	Address string `pg:",unique:acc"`
}
