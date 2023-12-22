package middlewares

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"google.dev/google/shuttle/core/app/manager/pkg/enum"
	"google.dev/google/shuttle/core/app/manager/pkg/errs"
	"google.dev/google/shuttle/core/app/manager/utils"
)

func GetUserInformationFromCtx(ctx context.Context) (*enum.AuthJWT, error) {
	fromContext, err := utils.GetUserInformationFromContext(ctx)
	return fromContext, err
}

// HasLoginFunc 如果是请求需要登录才能访问的接口，则需要判断是否带有 token ，并检测 token 的合法性，如果失败拒绝请求
func HasLoginFunc(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	_, err := GetUserInformationFromCtx(ctx)
	if err != nil {
		return nil, errs.LoginFailed
	}

	return next(ctx)
}
