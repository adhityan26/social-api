package helper

import (
	"github.com/jinzhu/gorm"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"regexp"
)

func SetupDB() *gorm.DB {
	conn := Getenv("USERNAME_DB", "root") + ":" +
		Getenv("PASSWORD_DB", "") + "@tcp(" +
		Getenv("DATABASE_HOST", "localhost") + ":" +
		Getenv("DATABASE_PORT", "3306") + ")/" +
		Getenv("DATABASE_NAME", "social-api") + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(Getenv("DATABASE_TYPE", "mysql"), conn)

	if err != nil {
		panic(err)
	}

	return db
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ValidateEMail(email string) bool {
	emailChecker := regexp.MustCompile(`([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})`)
	emailCheckerList := emailChecker.FindAllString(email, -1)

	return len(emailCheckerList) == 0
}