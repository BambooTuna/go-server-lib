package mysql

import (
	"fmt"
	"github.com/BambooTuna/go-server-lib/config"
	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GormConnection() *gorm.DB {
	conf := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.GetEnvString("MYSQL_USER", "user"),
		config.GetEnvString("MYSQL_PASS", "pass"),
		config.GetEnvString("MYSQL_HOST", "127.0.0.1"),
		config.GetEnvString("MYSQL_PORT", "3306"),
		config.GetEnvString("MYSQL_DATABASE", "sample_database"),
	)
	connection, _ := gorm.Open("mysql", conf)
	return connection
}
