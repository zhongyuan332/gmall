package model

import (
	"database/sql"
	"time"
)

// AdminUser 用户模型
type AdminUser struct {
	ID            int64          `json:"id"`
	Username      string         `json:"username"`
	Password      string         `json:"-"` // 密码不返回给前端
	Email         string         `json:"email"`
	RealName      string         `json:"real_name,omitempty"`
	Mobile        string         `json:"mobile,omitempty"`
	Avatar        sql.NullString `json:"avatar,omitempty"`
	Role          sql.NullString `json:"role"`
	Status        bool           `json:"status"`
	LastLoginTime *time.Time     `json:"last_login_time,omitempty"`
	LastLoginIP   sql.NullString `json:"last_login_ip,omitempty"`
	LoginCount    int            `json:"login_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     `json:"-"`
}
