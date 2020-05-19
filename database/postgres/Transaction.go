package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

type SearchUserTransactionUser struct {
	tableName struct{} `pg:"users"`
	ID        uint64   `json:"id" pg:",pk"`
	Login     string   `json:"login" pg:",unique,notnull"`
}

type SearchUserTransactionModel struct {
	*models.Transaction `pg:",inherit"`
	User                *SearchUserTransactionUser `json:"user"`
}

func (p postgres) SearchUserTransactions(id uint64, limit, offset int, sortMethod string) (count int, transactions []*SearchUserTransactionModel, err error) {
	query := p.Model(&transactions).
		Relation("User").
		WhereOr("transaction.id = ?", id).
		WhereOr("user_id = ?", id).
		Limit(limit).
		Offset(offset)
	switch sortMethod {
	case "New":
		query = query.Order("timestamp DESC")
	case "Old":
		query = query.Order("timestamp ASC")
	}
	count, err = query.SelectAndCount()
	return
}

func (p postgres) AddUserTransaction(transaction *models.Transaction) error {
	return p.Insert(transaction)
}
