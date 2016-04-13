// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package fb

import (
	"fmt"
	"os"
	"strings"


	"jex/cn/longGo/fb/goIni"
"github.com/labstack/gommon/log"
	"jex/cn/longGo/fb/db"
)

type (
	Config struct {
	AppName       string // 应用名称
	Env           string   // 运行模式
	Version       string //版本
	LogLevel      uint8
	HttpAddr      string // 应用监听地址，默认为空，监听所有的网卡 IP
	HttpPort      int    // 应用监听端口，默认为 8080
	//TplSuffix     string // 模板后缀名
	//TplLeft       string // 模板左定界符
	//TplRight      string // 模板右定界符
	DefaultModule string // 默认模块的名称
	DBConfig db.Config //数据库配置
}

)

func getConfig() Config {
	cfg, err := goIni.Load("conf/app.cnf")
	if err != nil {
		fmt.Println("\n  请确保在项目目录下运行，且存在配置文件 conf/app.cnf")
		os.Exit(1)
	}
	cnf := new(Config)
	section := cfg.Section("")
	cnf.AppName=section.Key("appname").String()
        cnf.Env=section.Key("Env").String()
	cnf.LogLevel=strToLogLevel(section.Key("loglevel").String())
	cnf.HttpAddr=section.Key("httpaddr").MustString("0.0.0.0")
	cnf.HttpPort=section.Key("httpport").MustInt(8080)
	cnf.Version=section.Key("version").String()
	cnf.DefaultModule= SnakeString(strings.Trim(section.Key("default_module").MustString("home"), "/"))
       //读取数据库参数
	//err = cfg.Section(cnf.Env).MapTo(cnf.DBConfig)
	section = cfg.Section(cnf.Env)
	cnf.DBConfig=db.Config{}
	cnf.DBConfig.DriverName=section.Key("adapter").String()
	cnf.DBConfig.DataSourceName=section.Key("database").String()
	cnf.DBConfig.Host=section.Key("host").String()
	cnf.DBConfig.Encoding=section.Key("encoding").String()
	cnf.DBConfig.UserName=section.Key("username").String()
	cnf.DBConfig.Password=section.Key("password").String()
        return *cnf
}

func strToLogLevel(logLevelstr string)uint8{
	var logLevel uint8

	switch strings.ToUpper(logLevelstr) {
	//case "TRACE":
	//	logLevel = log.TRACE
	case "DEBUG":
		logLevel = log.DEBUG
	case "INFO":
		logLevel = log.INFO
	//case "NOTICE":
	//	logLevel = log.NOTICE
	case "WARN":
		logLevel = log.WARN
	case "ERROR":
		logLevel = log.ERROR
	case "FATAL":
		logLevel = log.FATAL
	case "OFF":
		logLevel = log.OFF
	default:
		logLevel = log.DEBUG
	}
	return logLevel
}
