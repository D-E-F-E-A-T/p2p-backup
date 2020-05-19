package postgres

import (
	"time"

	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) SearchUserSettings(id uint64) (settings *models.Settings, err error) {
	settings = new(models.Settings)
	err = p.Model(settings).
		Where("id = ?", id).
		Select()
	if err == pg.ErrNoRows {
		settings.ID = id
		settings.UserID = id
		settings.Language = LanguageEnglish
		settings.Timezone, _ = time.Now().Zone()
		err = p.Insert(settings)
	}
	return
}

func (p postgres) GetUserAvailableAddresses(id uint64) (addresses []string, err error) {
	err = p.Model((*models.Settings)(nil)).
		Column("available_addresses").
		Where("id = ?", id).
		Select(&addresses)
	return
}

func (p postgres) SetUserAvailableAddresses(id uint64, addresses []string) error {
	_, err := p.Model((*models.Settings)(nil)).
		Where("id = ?", id).
		Set("available_addresses = ?", addresses).
		Update()
	return err
}

type UserSettngs struct {
	*models.User `pg:",inherit"`
	Settings     *models.Settings
}

func (p postgres) GetUserSettings(id uint64) (userSettings *UserSettngs, err error) {
	userSettings = new(UserSettngs)
	return userSettings, p.Model(userSettings).
		Relation("Settings").
		Where("user_id = ?", id).
		Select()
}
