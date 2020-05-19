package l18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var savedTranslations map[string]map[string]string

func init() {
	savedTranslations = make(map[string]map[string]string)
	if err := filepath.Walk(viper.GetString("api.l18n"),
		func(path string, info os.FileInfo, err error) (err2 error) {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return
			}
			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			var result map[string]string
			if err := json.Unmarshal(raw, &result); err != nil {
				return err
			}
			index := strings.Index(info.Name(), "-")
			savedTranslations[info.Name()[:index]] = result
			return
		}); err != nil {
		log.Panicf("Error with reading l18n files: %v", err)
	}
}

func T(language string, key string, args ...interface{}) string {
	translations, ok := savedTranslations[language]
	if !ok {
		return ""
	}
	result, ok := translations[key]
	if !ok {
		return ""
	}
	if len(args) == 0 {
		return result
	}
	return fmt.Sprintf(result, args)
}
