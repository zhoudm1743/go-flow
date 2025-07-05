package resp

import (
	"time"
)

// AccountResp 响应结构体
type AccountResp struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	LastLogin string    `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
