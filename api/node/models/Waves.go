package models

type Waves struct {
	Account string `pg:",pk"`
	Address string `pg:",pk"`
}
