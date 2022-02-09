package zcode

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = ""

// Open 连接数据库
func OpenDB() (db *gorm.DB, err error) {

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		AllowGlobalUpdate:      false, // 不允许全局update
		PrepareStmt:            true,  // 缓存预编译语句，执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		SkipDefaultTransaction: true,  // 禁用默认事务(使用事务能确保数据完整，但会降低性能)，这将获得大约 30%+ 性能提升
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(8)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(128)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return
}
