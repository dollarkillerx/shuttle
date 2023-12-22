package utils

import (
	"context"
	"net/http"

	"github.com/dollarkillerx/jwt"
	"github.com/pkg/errors"
	"google.dev/google/common/pkg/logger"
	"google.dev/google/shuttle/core/app/manager/pkg/enum"
)

func GetTokenFromHeader(header http.Header) (tokenString string, ok bool) {
	tk := header.Get("Authorization")
	if tk == "" {
		return "", false
	}

	return tk, true
}

func GetTokenByContext(ctx context.Context) (tokenString string, ok bool) {
	value := ctx.Value(enum.TokenCtxKey)

	if value == nil {
		return "", false
	}

	tokenString, ok = value.(string)

	return tokenString, true
}

func GetUserInformationFromContext(ctx context.Context) (*enum.AuthJWT, error) {
	tokenString, ok := GetTokenByContext(ctx)
	if !ok {
		return nil, errors.New("not auth token")
	}

	token, err := jwt.TokenFormatString(tokenString)
	if err != nil {
		logger.Info(err)
		return nil, err
	}

	err = JWT.VerificationSignature(token)
	if err != nil {
		logger.Info(err)
		return nil, err
	}

	var tk enum.AuthJWT
	err = token.Payload.Unmarshal(&tk)
	if err != nil {
		logger.Info(err)
		return nil, err
	}

	return &tk, nil
}

func GetRequestIdByContext(ctx context.Context) string {
	value := ctx.Value(enum.RequestId)
	if value == nil {
		return ""
	}

	requestId, ok := value.(string)
	if !ok {
		return ""
	}

	return requestId
}

func GetRequestIPByContext(ctx context.Context) string {
	value := ctx.Value(enum.ReqIP)
	if value == nil {
		return ""
	}

	requestId, ok := value.(string)
	if !ok {
		return ""
	}

	return requestId
}
