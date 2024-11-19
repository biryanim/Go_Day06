package auth

import (
	"context"
	"errors"
	"go_day06/pkg/entities/admin"
	"go_day06/pkg/entities/token"
)

type Auth struct {
	admin *admin.Admin
}

func New(admin *admin.Admin) *Auth {
	return &Auth{admin: admin}
}

func (a *Auth) Login(c context.Context, username, password string) (string, error) {
	if username == a.admin.Username && password == a.admin.Password {
		return token.BuildJwtString()
	}
	return "", errors.New("no such user or password")
}

func (a *Auth) VerifyToken(ctx context.Context, t string) error {
	return token.Verify(t)
}
