package db

import (
	"fmt"
	"github.com/jinzhu/gorm"

	//_ "github.com/lib/pq"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/mattn/go-sqlite3"
	//_ "github.com/nakagami/firebirdsql"
)
var (
     db_Tables  []interface{}
)

type LongDB struct {
	*gorm.DB
	Cnf  Config
}




type Config struct {
	//"sqlite3"
	DriverName string
	//"/tmp/post_db.bin"
	DataSourceName string
	UserName string
	Password string
	Host string
	Encoding string
}

//初始化数据接口
type IDB interface {
	InitData()
}


func (this *LongDB)Init(cnf Config)*LongDB{
	var err error

	switch cnf.DriverName {
	case "sqlite3":
		this.DB, err = gorm.Open("sqlite3", cnf.DataSourceName)
		// logger.CheckFatal(err, "Got error when connect database")
	if err!=nil{
		fmt.Println(err)
	}
	//this.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	case "mysql":
		this.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			cnf.UserName,
			cnf.Password,
			cnf.Host,
			"mysql",
			cnf.Encoding))
		//创建数据库
		this.DB.Debug().Exec("CREATE DATABASE IF NOT EXISTS `"+cnf.DataSourceName+"` DEFAULT CHARSET utf8 COLLATE utf8_general_ci;")

		//db, err := sql.Open("mysql", "user:password@tcp(localhost:5555)/dbname?charset=utf8&parseTime=True&loc=Local")
		this.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			cnf.UserName,
			cnf.Password,
			cnf.Host,
			cnf.DataSourceName,
			cnf.Encoding))
		//logger.CheckFatal(err,"Got error when connect database")
		if err!=nil{
			fmt.Println("db open faild:",err)

		}

	//this.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	default:
		fmt.Println("The types of database does not support:"+cnf.DriverName , err)
	}

	this.DB.DB()
	this.DB.DB().Ping()
	this.DB.DB().SetMaxIdleConns(10)
	this.DB.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	this.SingularTable(true)

	// construct a gorp DbMap

	//
	this.Cnf=cnf
	return  this
}



func (this *LongDB)AutoMigrate(tables ...interface{}) *LongDB{
	if len(tables)==0{
		tables=db_Tables
	}
	if this.Cnf.DriverName=="mysql" {
		this.DB.Set("gorm:table_options", fmt.Sprintf("ENGINE=InnoDB CHARSET=%s",this.Cnf.Encoding)).AutoMigrate(tables...)

	}else{
		this.DB.AutoMigrate(tables...)
	}

	//初始化数据表数据
	for _,value:=range  tables {
		if idb,ok:=value.(IDB);ok{
			idb.InitData()
		}
	}
	return this
}



func AddTable(table interface{}) {
	db_Tables = append(db_Tables, table)
}