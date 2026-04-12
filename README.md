# LgyTool

> LgyTool 是一套 Golang 工具包，提供了丰富的开发过程中所需的工具包，助力 Golang 高效开发 ~

## 项目简介

LgyTool 是一个功能丰富的 Go 语言工具库，包含多个实用模块，为 Go 开发者提供便捷的工具函数和组件。

## 支持的数据库

- MySQL
- MSSQL (SQL Server 2008 R2+)
- PostgreSQL
- SQLite

## 主要功能模块

### 1. 数据库操作 (db)
- 数据库连接池管理
- 支持多种数据库类型
- 连接池状态监控
- 统一的数据库操作接口

### 2. 云服务 (cloud)
- 阿里云 OSS
- 腾讯云 COS
- 阿里云短信服务

### 3. 日期时间处理 (datetime)
- 日期时间格式化
- 时间戳转换
- 时间计算

### 4. 图像处理 (gimage)
- 图片处理
- 验证码生成

### 5. Web 工具 (gintool)
- Gin 框架工具
- 参数处理
- 文件上传

### 6. 网络工具
- TCP 消息处理 (tcpMessage)
- HTTP 请求 (request)
- IP 处理 (ip)

### 7. 其他工具
- 压缩功能 (gZip)
- 文件系统操作 (gfs)
- 字符串处理 (gstring)
- 加密工具 (gmd5)
- 缓存工具 (mapCache, mapDetailCache)
- 随机数生成 (random)
- 切片操作 (slice)
- 第三方登录 (thirdPartyLogin)

## 安装

```bash
go get github.com/golgys0621/lgystool
```

## 快速开始

### 数据库连接池示例

```go
import (
    "github.com/golgys0621/lgystool/db"
)

func main() {
    // 初始化数据库连接池
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
    }
    db.Start(dbConfigs)

    // 获取数据库连接
    mysqlDB := db.Init("MySQL")

    // 获取连接池状态
    status := db.GetPoolStatus("MySQL")
    fmt.Printf("连接池状态: %v\n", status)

    // 关闭连接池
    db.Close()
}
```

### TCP 消息处理示例

```go
import (
    "github.com/golgys0621/lgystool"
    "net"
)

func main() {
    // 服务器端
    listener, _ := net.Listen("tcp", ":8080")
    for {
        conn, _ := listener.Accept()
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    // 读取消息
    msg, _ := lgystool.ReadTCPResponse(conn)
    fmt.Printf("收到消息: %s\n", msg)
    
    // 发送响应
    response := []byte("Hello, Client!")
    lgystool.WriteTCPResponse(conn, response)
}
```

## 版本要求

- Go 1.26.2+

## 依赖管理

项目使用 Go Modules 进行依赖管理，主要依赖包括：

- gorm.io/gorm
- gorm.io/driver/mysql
- gorm.io/driver/sqlserver
- gorm.io/driver/postgres
- gorm.io/driver/sqlite
- github.com/gin-gonic/gin
- github.com/redis/go-redis/v9

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 许可证

MIT License
