package middleware

import (
	"log"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var Valid = validator.New()

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("fiat", isFiat)
		v.RegisterValidation("crypto", isCrypto)
		v.RegisterValidation("currency", isCurrency)
		v.RegisterValidation("language", IsLanguage)
		v.RegisterValidation("country", isCountry)
		v.RegisterValidation("provider", isProvider)
	}
}

func isFiat(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		return IsFiatString(field.String())
	case reflect.Uint8:
		return IsFiatUint(uint8(field.Uint()))
	default:
		log.Printf("%v", field.Kind())
	}
	return false
}

func isCrypto(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		return IsCryptoString(field.String())
	case reflect.Uint8:
		return IsCryptoUint(uint8(field.Uint()))
	default:
		log.Printf("%v", field.Kind())
	}
	return false
}

func isCurrency(fl validator.FieldLevel) bool {
	return isCrypto(fl) || isFiat(fl)
}

func IsLanguage(fl validator.FieldLevel) bool {
	language := fl.Field().String()
	for _, lang := range Languages {
		if lang == language {
			return true
		}
	}
	return false
}

func isCountry(fl validator.FieldLevel) bool {
	c := fl.Field().String()
	for _, country := range Countries {
		if c == country {
			return true
		}
	}
	return false
}

func isProvider(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		key := field.String()
		for _, value := range Providers {
			if value == key {
				return true
			}
		}
		return false
	case reflect.Uint8:
		value := uint8(field.Uint())
		for key, _ := range Providers {
			if key == value {
				return true
			}
		}
		return false
	}
	return false
}
