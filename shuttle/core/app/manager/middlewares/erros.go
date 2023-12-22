package middlewares

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.dev/google/common/pkg/logger"
	"google.dev/google/shuttle/core/app/manager/pkg/errs"
	"google.dev/google/shuttle/core/app/manager/utils"
)

// MiddleError 全局錯誤處理
func MiddleError(ctx context.Context, err error) *gqlerror.Error {
	logger.Errorf("[ReqID]:%s\n[Message]:%+v", utils.GetRequestIdByContext(ctx), err)

	var gqlErr *errs.Error
	if !errors.As(err, &gqlErr) {
		return gqlerror.WrapPath(graphql.GetPath(ctx), errs.SystemError(err))
	}

	return gqlerror.WrapPath(graphql.GetPath(ctx), err)
}
