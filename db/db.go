package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var GoToolDBMap = make(map[string]*gorm.DB)

// 初始化数据库连接池
func Start(dbConfigs map[string]map[string]string) {
	// 遍历数据库配置初始化连接池
	for k, conf := range dbConfigs {
		// 日志级别
		var loggerType logger.LogLevel
		if conf["RunMode"] == "dev" {
			loggerType = logger.Info
		} else {
			loggerType = logger.Silent
		}
		// 配置文件
		options := &gorm.Config{
			// 开启事务保证数据一致性
			SkipDefaultTransaction: false,
			// 日志
			Logger: logger.Default.LogMode(loggerType),
			// 命名策略
			NamingStrategy: schema.NamingStrategy{
				// 表前缀
				TablePrefix: conf["TablePrefix"],
				// 单数表名称
				SingularTable: true,
			},
			// 建表时候是否忽略外键
			DisableForeignKeyConstraintWhenMigrating: true,
		}
		// 创建连接池
		var DSN string = ""
		var err error
		switch conf["DBType"] {
		// MySQL
		case "MySQL":
			if conf["RunMode"] == "dev" {
				DSN = conf["UsernameDev"] + ":" +
					conf["PasswordDev"] + "@" +
					"tcp(" + conf["HostDev"] + ":" + conf["Port"] + ")/" +
					conf["DatabaseName"] + "?charset=" +
					conf["Charset"] + "&parseTime=True&loc=Local"
			} else {
				DSN = conf["Username"] + ":" +
					conf["Password"] + "@" +
					"tcp(" + conf["Host"] + ":" + conf["Port"] + ")/" +
					conf["DatabaseName"] + "?charset=" +
					conf["Charset"] + "&parseTime=True&loc=Local"
			}
			GoToolDBMap[k], err = gorm.Open(mysql.Open(DSN), options)
			if err != nil {
				log.Println("✘ 连接库连接池 : " + k + " 初始化失败 ( " + err.Error() + " )")
			} else {
				log.Println("✔ 连接库连接池 : " + k + " 初始化成功 ")
			}
		// MSSQL
		case "MSSQL":
			if conf["RunMode"] == "dev" {
				DSN = "sqlserver://" + conf["UsernameDev"] + ":" +
					conf["PasswordDev"] + "@" +
					conf["HostDev"] + ":" + conf["Port"] + "?database=" +
					conf["DatabaseName"] + "&encrypt=disable&trustservercertificate=true"
			} else {
				DSN = "sqlserver://" + conf["Username"] + ":" +
					conf["Password"] + "@" +
					conf["Host"] + ":" + conf["Port"] + "?database=" +
					conf["DatabaseName"] + "&encrypt=disable&trustservercertificate=true"
			}
			GoToolDBMap[k], err = gorm.Open(sqlserver.Open(DSN), options)
			if err != nil {
				log.Println("✘ 连接库连接池 : " + k + " 初始化失败 ( " + err.Error() + " )")
			} else {
				log.Println("✔ 连接库连接池 : " + k + " 初始化成功 ")
			}
		// PostgreSQL
		case "PostgreSQL":
			if conf["RunMode"] == "dev" {
				DSN = "host=" + conf["HostDev"] + " user=" + conf["UsernameDev"] + " password=" + conf["PasswordDev"] + " dbname=" + conf["DatabaseName"] + " port=" + conf["Port"] + " sslmode=disable"
			} else {
				DSN = "host=" + conf["Host"] + " user=" + conf["Username"] + " password=" + conf["Password"] + " dbname=" + conf["DatabaseName"] + " port=" + conf["Port"] + " sslmode=disable"
			}
			GoToolDBMap[k], err = gorm.Open(postgres.Open(DSN), options)
			if err != nil {
				log.Println("✘ 连接库连接池 : " + k + " 初始化失败 ( " + err.Error() + " )")
			} else {
				log.Println("✔ 连接库连接池 : " + k + " 初始化成功 ")
			}
		// SQLite
		case "SQLite":
			DSN = conf["DatabaseName"]
			GoToolDBMap[k], err = gorm.Open(sqlite.Open(DSN), options)
			if err != nil {
				log.Println("✘ 连接库连接池 : " + k + " 初始化失败 ( " + err.Error() + " )")
			} else {
				log.Println("✔ 连接库连接池 : " + k + " 初始化成功 ")
			}
		}
		// 获取基础数据库操作接口
		sqlDB, _ := GoToolDBMap[k].DB()
		//设置数据库连接池参数
		// 最大连接数
		MaxOpenConns, err := strconv.Atoi(conf["MaxOpenConns"])
		if err != nil || MaxOpenConns <= 0 {
			MaxOpenConns = 10 // 默认值
		}
		sqlDB.SetMaxOpenConns(MaxOpenConns)
		// 最大空闲连接数
		MaxIdleConns, err := strconv.Atoi(conf["MaxIdleConns"])
		if err != nil || MaxIdleConns <= 0 {
			MaxIdleConns = 5 // 默认值
		}
		sqlDB.SetMaxIdleConns(MaxIdleConns)
		// 最大连接时间
		MaxLifetime, err := strconv.Atoi(conf["MaxLifetime"])
		if err != nil || MaxLifetime <= 0 {
			MaxLifetime = 3600 // 默认值（秒）
		}
		sqlDB.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
		// 连接最大空闲时间
		MaxIdleTime, err := strconv.Atoi(conf["MaxIdleTime"])
		if err != nil || MaxIdleTime <= 0 {
			MaxIdleTime = 600 // 默认值（秒）
		}
		sqlDB.SetConnMaxIdleTime(time.Duration(MaxIdleTime) * time.Second)
	}
}

// 获取数据库操作对象
func Init(configName ...string) *gorm.DB {
	if len(configName) < 1 {
		configName = append(configName, "DB")
	}
	gormDB, ok := GoToolDBMap[configName[0]]
	if ok {
		return gormDB
	}
	panic("✘ 数据库连接池 [ " + configName[0] + " ] 初始化失败")
}

// 将 map 对象转换为 sql 条件
func MapToWhere(mapData map[string][]any) (string, []any) {
	var whereSql = make([]string, 0)
	var whereVal = make([]any, 0)
	for k, item := range mapData {
		whereSql = append(whereSql, fmt.Sprintf("%v %v %v ?", item[0], k, item[1]))
		whereVal = append(whereVal, item[2])
	}
	return strings.Join(whereSql, " "), whereVal
}

// 获取连接池状态
func GetPoolStatus(configName ...string) map[string]int {
	if len(configName) < 1 {
		configName = append(configName, "DB")
	}
	gormDB, ok := GoToolDBMap[configName[0]]
	if !ok {
		return nil
	}
	sqlDB, _ := gormDB.DB()
	stats := sqlDB.Stats()
	return map[string]int{
		"OpenConnections":   stats.OpenConnections,
		"InUse":             stats.InUse,
		"Idle":              stats.Idle,
		"WaitCount":         int(stats.WaitCount),
		"WaitDuration":      int(stats.WaitDuration.Seconds()),
		"MaxIdleClosed":     int(stats.MaxIdleClosed),
		"MaxLifetimeClosed": int(stats.MaxLifetimeClosed),
	}
}

// 关闭指定数据库连接池
func Close(configName ...string) {
	if len(configName) < 1 {
		// 关闭所有连接池
		for k, db := range GoToolDBMap {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
				log.Println("✔ 连接库连接池 : " + k + " 已关闭")
			}
		}
		GoToolDBMap = make(map[string]*gorm.DB)
	} else {
		// 关闭指定连接池
		if db, ok := GoToolDBMap[configName[0]]; ok {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
				log.Println("✔ 连接库连接池 : " + configName[0] + " 已关闭")
			}
			delete(GoToolDBMap, configName[0])
		}
	}
}
