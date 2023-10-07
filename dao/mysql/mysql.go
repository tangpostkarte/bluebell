package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("Connect mysql failed", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.max_connection"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_connection"))

	return
}

func Close() {
	_ = db.Close()
}
