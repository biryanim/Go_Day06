package admin

import (
	"go_day06/pkg/config"
)

type Admin struct {
	Username string
	Password string
}

func NewAdmin(cfg *config.AppConfig) *Admin {
	return &Admin{Username: cfg.AdminUser, Password: cfg.AdminPassword}
}
