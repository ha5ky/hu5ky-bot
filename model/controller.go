/**
 * @Author Nil
 * @Description model/controller.go
 * @Date 2023/3/28 17:06
 **/

package model

import (
	"errors"
	"fmt"
	"github.com/ha5ky/hu5ky-bot/model/base"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mysqllogger "gorm.io/gorm/logger"
	"sync"
)

var (
	once sync.Once
	db   *gorm.DB

	MySQLLogLevel = map[string]mysqllogger.LogLevel{
		"silent": mysqllogger.Silent,
		"error":  mysqllogger.Error,
		"warn":   mysqllogger.Warn,
		"info":   mysqllogger.Info,
	}
)

type Controller struct {
	// normal db and transaction
	controller *gorm.DB
	// isTx
	isTx bool
}

func NewController() *Controller {
	return &Controller{
		controller: GetDB(),
	}
}

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
			config.SysCache.DB.Mysql.User,
			config.SysCache.DB.Mysql.Pwd,
			config.SysCache.DB.Mysql.Host,
			config.SysCache.DB.Mysql.Port,
			config.SysCache.DB.Mysql.DBName,
			config.SysCache.DB.Mysql.Charset,
			config.SysCache.DB.Mysql.ParseTime,
			config.SysCache.DB.Mysql.Loc,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: mysqllogger.Default.LogMode(MySQLLogLevel[config.SysCache.DB.Mysql.LogLevel]),
		})
		if err != nil {
			panic(err)
		}
	})
	return db
}

func (c *Controller) CreateTables(tables []base.AutoRegister) (err error) {
	//err := c.Begin()
	for i := range tables {
		if !c.controller.Migrator().HasTable(tables[i]) {
			err = c.controller.Migrator().CreateTable(tables[i])
			if err != nil {
				logger.Errorf("can not create table: %s: %s", tables[i].TableName(), err.Error())
				_ = c.Rollback()
				panic(err)
			}
		}
	}
	//c.Commit()
	return
}

func (c *Controller) Begin() error {
	c.controller = c.controller.Begin()
	if c.controller.Error != nil {
		return c.controller.Error
	}
	c.isTx = true
	return nil
}

func (c *Controller) Commit() error {
	if !c.isTx {
		return errors.New("it is not a transaction")
	}
	c.controller = c.controller.Commit()
	if c.controller.Error != nil {
		return c.controller.Error
	}
	return nil
}

func (c *Controller) Rollback() error {
	if !c.isTx {
		return errors.New("it is not a transaction")
	}
	c.controller = c.controller.Rollback()
	if c.controller.Error != nil {
		return c.controller.Error
	}
	return nil
}
