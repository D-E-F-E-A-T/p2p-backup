package postgres

import (
	"fmt"

	"github.com/dmitriy-vas/p2p/models"
)

type Comment struct {
	*models.Comment `pg:",inherit"`
	FromUser        string `json:"from_user" pg:"from_user"`
	ToUser          string `json:"to_user" pg:"to_user"`
}

//func (p postgres) SearchUserComments(id uint64, limit, offset int, sortMethod string) (count int, comments []*Comment, err error) {
//	query := p.Model(&comments).
//		Relation("User.login").
//		Where("comment.id = ?", id).
//		Limit(limit).
//		Offset(offset)
//	switch sortMethod {
//	case "New":
//		query = query.Order("timestamp DESC")
//	case "Old":
//		query = query.Order("timestamp ASC")
//	}
//	count, err = query.SelectAndCount()
//	return
//}

// TODO add default relation
func (p postgres) SearchUserComments(id uint64, limit, offset int, sortMethod string) (count int, comments []*Comment, err error) {
	queryTemplate1 := `SELECT "comments".id,
       user_id,
       deal,
       message,
       rating,
       timestamp,
      "from_user"."login" from_user,
       "to_user"."login"   to_user`
	queryTemplate2 := ` FROM comments
                LEFT JOIN users AS from_user ON from_user.id = comments.user_id
                LEFT JOIN users AS to_user ON to_user.id = comments.id
			WHERE comments.id = ?0`
	switch sortMethod {
	case "New":
		sortMethod = "DESC"
	case "Old":
		sortMethod = "ASC"
	}
	if _, err = p.Query(&comments, queryTemplate1+queryTemplate2+fmt.Sprintf(` ORDER BY timestamp %s `, sortMethod)+`LIMIT ?1 OFFSET ?2;`, id, limit, offset); err != nil {
		return count, comments, err
	}
	_, err = p.Query(&count, `SELECT COUNT(*)`+queryTemplate2, id)
	return
}

func (p postgres) AddNewComment(comment *models.Comment) error {
	return p.Insert(comment)
}
