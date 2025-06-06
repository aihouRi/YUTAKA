package domain

import (
	"database/sql"
	"time"
)

// User 结构体对应数据库中的 users 表
type User struct {
	ID        int64 // 用 int64 对应数据库的 INT 或 BIGINT 主键
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time // DATETIME 也能被 parseTime=True 解析为 time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}