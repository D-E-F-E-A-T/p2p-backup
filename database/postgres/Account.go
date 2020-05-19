package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) RegisterNewUser(user *models.User) error {
	return p.Insert(user)
}

func (p postgres) SearchUser(credentials string) (user *models.User, err error) {
	user = new(models.User)
	err = p.Model(user).
		WhereOr("login = ?", credentials).
		WhereOr("email = ?", credentials).
		Select()
	return
}

func (p postgres) SearchUserByActivation(activation string) (user *models.User, err error) {
	user = new(models.User)
	err = p.Model(user).
		Where("activation = ?", activation).
		Select()
	return
}

func (p postgres) DeleteUserActivation(id uint64) error {
	_, err := p.Model((*models.User)(nil)).
		Where("id = ?", id).
		Set("activation = ?", nil).
		Update()
	return err
}

func (p postgres) RestoreUserPassword(id uint64, newPassword string) error {
	_, err := p.Model((*models.User)(nil)).
		Where("id = ?", id).
		Set("password = ?", newPassword).
		Update()
	return err
}

func (p postgres) SearchUserAndReturnEmail(id uint64) (email string, err error) {
	err = p.Model((*models.User)(nil)).
		Column("email").
		Where("id = ?", id).
		Select(&email)
	return
}

type User struct {
	*models.User `pg:",inherit"`
	Statistic    *models.Statistic `json:"statistic"`
}

func (p postgres) GetUserInfo(id uint64) (user *User, err error) {
	user = new(User)
	return user, p.Model(user).
		Relation("Statistic").
		Where("user_id = ?", id).
		Column("user.description", "user.login", "user.groups", "user.id").
		Select()
}

func (p postgres) GetUserLogin(id uint64) (login string, err error) {
	return login, p.Model((*models.User)(nil)).
		Where("id = ?", id).
		Column("login").
		Select(&login)
}

func (p postgres) SetUserPassword(id uint64, encodedPassword string) error {
	_, err := p.Model((*models.User)(nil)).
		Where("id = ?", id).
		Set("password = ?", encodedPassword).
		Update()
	return err
}

type Users struct {
	User
	Settings *models.Settings `json:"settings"`
}

func (p postgres) GetUsers(search string, limit, offset int) (count int, users []*Users, err error) {
	count, err = p.Model(&users).
		Relation("Statistic").
		Relation("Settings").
		Column("user.description",
			"user.login",
			"user.id",
			"user.email",
			"user.groups").
		WhereOr("login ~ ?", search).
		WhereOr("email ~ ?", search).
		Limit(limit).
		Offset(offset).
		SelectAndCount()
	return
}

func (p postgres) AddUserBan(ban *models.Ban) error {
	return p.Insert(ban)
}

func (p postgres) DeleteSessions(id uint64) error {
	_, err := p.Model((*models.Session)(nil)).
		Where("id = ?", id).
		Delete()
	return err
}

func (p postgres) SearchUserBan(id uint64) (ban *models.Ban, err error) {
	ban = new(models.Ban)
	return ban, p.Model(ban).
		Where("user_id = ?", id).
		Select()
}

func (p postgres) UpdateUserGroups(id uint64, groups []int) error {
	_, err := p.Model((*models.User)(nil)).
		Where("id = ?", id).
		Set("groups = ?", groups).
		Update()
	return err
}
