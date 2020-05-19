package postgres

import (
	"fmt"

	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) GetCountries(language string) (countries []*models.Country, err error) {
	_, err = p.Query(&countries, fmt.Sprintf(`SELECT "%s"."id", "%s"."value" FROM %s`, language, language, language))
	return
}
