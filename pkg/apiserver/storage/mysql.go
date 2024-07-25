package storage

import (
	"database/sql"
	"fmt"
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/options"
	"github.com/HugoWw/x_apiserver/pkg/clog"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func InitMysqlDBObj(option *options.MysqlClientOptions) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		option.UserName,
		option.Password,
		option.Host,
		option.Port,
		option.Database)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		clog.Logger.Errorf("create mysql connect err: %v", err)
		return nil, err
	}

	db.SetMaxIdleConns(option.MaxIdleConns)
	db.SetMaxOpenConns(option.MaxOpenConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(option.ConnMaxLifetime))
	err = db.Ping()
	if err != nil {
		clog.Logger.Errorf("ping mysql connection failed: %v", err)
		return nil, err
	}

	return db, nil
}
