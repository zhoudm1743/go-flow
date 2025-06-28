package resp

import "github.com/zhoudm1743/go-flow/pkg/types"

type SystemAdminResp struct {
	ID          uint         `json:"id" structs:"id"`
	Username    string       `json:"username" structs:"username"`
	Nickname    string       `json:"nickname" structs:"nickname"`
	Status      int8         `json:"status" structs:"status"`
	Role        []string     `json:"role" structs:"role"`
	Email       string       `json:"email" structs:"email"`
	Phone       string       `json:"phone" structs:"phone"`
	Avatar      string       `json:"avatar" structs:"avatar"`
	LastLoginAt types.TsTime `json:"last_login_at" structs:"last_login_at"`
	LastLoginIP string       `json:"last_login_ip" structs:"last_login_ip"`
	CreatedAt   types.TsTime `json:"created_at" structs:"created_at"`
	UpdatedAt   types.TsTime `json:"updated_at" structs:"updated_at"`
}
