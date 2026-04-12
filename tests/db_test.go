package tests

import (
	"log"
	"testing"

	"github.com/golgys0621/lgystool/db"
)

func TestStart(t *testing.T) {
	// 测试配置
	dbConfigs := map[string]map[string]string{
		"MySQL": {
			"DBType":       "MySQL",
			"RunMode":      "dev",
			"HostDev":      "localhost",
			"Port":         "3306",
			"UsernameDev":  "root",
			"PasswordDev":  "root",
			"DatabaseName": "test",
			"Charset":      "utf8mb4",
			"TablePrefix":  "",
			"MaxOpenConns": "10",
			"MaxIdleConns": "5",
			"MaxLifetime":  "3600",
		},
		"MSSQL": {
			"DBType":       "MSSQL",
			"RunMode":      "dev",
			"HostDev":      "localhost",
			"Port":         "1433",
			"UsernameDev":  "sa",
			"PasswordDev":  "password",
			"DatabaseName": "test",
			"TablePrefix":  "",
			"MaxOpenConns": "10",
			"MaxIdleConns": "5",
			"MaxLifetime":  "3600",
		},
		"PostgreSQL": {
			"DBType":       "PostgreSQL",
			"RunMode":      "dev",
			"HostDev":      "localhost",
			"Port":         "5432",
			"UsernameDev":  "postgres",
			"PasswordDev":  "postgres",
			"DatabaseName": "test",
			"TablePrefix":  "",
			"MaxOpenConns": "10",
			"MaxIdleConns": "5",
			"MaxLifetime":  "3600",
		},
		"SQLite": {
			"DBType":       "SQLite",
			"RunMode":      "dev",
			"DatabaseName": ":memory:",
			"TablePrefix":  "",
			"MaxOpenConns": "10",
			"MaxIdleConns": "5",
			"MaxLifetime":  "3600",
		},
	}

	// 初始化连接池
	db.Start(dbConfigs)

	// 测试获取数据库连接
	testDB := db.Init("SQLite")
	if testDB == nil {
		t.Error("Failed to initialize SQLite database")
	}

	// 测试连接池状态
	status := db.GetPoolStatus("SQLite")
	if status == nil {
		t.Error("Failed to get pool status")
	} else {
		log.Printf("SQLite pool status: %v", status)
	}

	// 测试关闭指定连接池
	db.Close("SQLite")

	// 测试关闭所有连接池
	db.Close()

	log.Println("Test completed successfully")
}
