package common

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var print = fmt.Print

type MysqlConfig struct {
	Host      string `json:"host"`
	Port      int32  `json:"port"`
	DbName    string `json:"db_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ShowSql   bool   `json:"show_sql"`
	ShowDebug bool   `json:"show_debug"`
	ShowErr   bool   `json:"show_err"`
	ShowWarn  bool   `json:"show_warn"`
}

func NewMysqlConf(m json.RawMessage) (conf *MysqlConfig, err error) {
	conf = &MysqlConfig{}
	err = json.Unmarshal(m, conf)
	return
}

func EngineForCfg(cfg *MysqlConfig) (engine *xorm.Engine, err error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	engine, err1 := xorm.NewEngine("mysql", dataSourceName)
	if err1 != nil {
		err = err1
		return
	}

	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "t_")
	engine.SetTableMapper(tbMapper)
	engine.SetColumnMapper(core.SameMapper{})
	engine.ShowSQL = cfg.ShowSql
	engine.ShowDebug = cfg.ShowDebug
	engine.ShowErr = cfg.ShowErr
	engine.ShowWarn = cfg.ShowWarn
	err = engine.Ping()
	return
}

func NewDBService(m json.RawMessage) (db *DBService, err error) {
	conf, err1 := NewMysqlConf(m)
	if err1 != nil {
		err = err1
		return
	}
	db = &DBService{DBCfg: conf}
	return
}

type DBService struct {
	Engine *xorm.Engine
	DBCfg  *MysqlConfig
}

func (db *DBService) Init() {
	engine, err := EngineForCfg(db.DBCfg)

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db.DBCfg.Username, db.DBCfg.Password, db.DBCfg.Host, db.DBCfg.Port, db.DBCfg.DbName)

	if err != nil {
		panic(fmt.Sprintf("connect %s failed,reason:%s", dataSourceName, err.Error()))
	}

	logger.Printf("connect %s success\n", dataSourceName)
	db.Engine = engine
}

func (d *DBService) AfterInit() {

}

func (d *DBService) BeforeDestroy() {

}

func (d *DBService) Destroy() {

}
