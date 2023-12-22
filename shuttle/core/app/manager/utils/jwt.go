package utils

import (
	"context"

	"github.com/dollarkillerx/jwt"
	"github.com/pkg/errors"
	"google.dev/google/shuttle/core/app/manager/conf"
	"google.dev/google/shuttle/core/app/manager/pkg/enum"
)

var JWT *jwt.JWT

func InitJWT() {
	JWT = jwt.NewJwt(conf.CONFIG.JWTConfiguration.SecretKey)
}

// GetAuthModel GetAuthModel
func GetAuthModel(ctx context.Context) (enum.AuthJWT, error) {
	auth := ctx.Value(enum.TokenCtxKey.String())
	if auth == nil {
		return enum.AuthJWT{}, errors.New("what fuck JWTToken is not exists")
	}

	model, ok := auth.(enum.AuthJWT)
	if !ok {
		return enum.AuthJWT{}, errors.New("what fuck JWTToken is not exists 2")
	}

	return model, nil
}
