package resolvers

import (
	"context"
	"log"

	"google.dev/google/shuttle/core/app/manager/generated"
	"google.dev/google/shuttle/core/app/manager/pkg/enum"
	"google.dev/google/shuttle/core/app/manager/pkg/errs"
	"google.dev/google/shuttle/core/app/manager/pkg/models"
	"google.dev/google/shuttle/core/app/manager/utils"
)

func (r *mutationResolver) SendEmail(ctx context.Context, input *generated.SendEmailRequest) (*generated.SendEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UserRegistration(ctx context.Context, input *generated.UserRegistration) (*generated.AuthPayload, error) {
	// search app info
	var appInfo models.App
	err := r.Storage.DB().Model(&models.App{}).Where("app_id = ?", input.AppID).First(&appInfo).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !appInfo.NoAuthenticationRequired && input.CaptchaID == "" {
		return nil, errs.CaptchaCode2
	} else {
		// TODO: 验证captcha
	}

	// 注册用户
	var user = models.User{
		OS:         input.Os.String(),
		Token:      input.Token,
		DeviceID:   input.DeviceID,
		DeviceName: input.DeviceName,
		AppID:      input.AppID,
		RegIP:      utils.GetRequestIPByContext(ctx),
	}

	err = r.Storage.DB().Model(&models.User{}).
		Where("app_id = ?", input.AppID).
		Where("token = ?", input.Token).
		FirstOrCreate(&user).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token, err := utils.JWT.CreateToken(&enum.AuthJWT{
		generated.UserInformation{
			Os:         input.Os,
			Token:      input.Token,
			DeviceName: input.DeviceName,
			DeviceID:   input.DeviceID,
			AppID:      input.AppID,
		},
	}, 0)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &generated.AuthPayload{
		AccessTokenString: token,
	}, err
}

func (r *queryResolver) User(ctx context.Context) (*generated.UserInformation, error) {
	authJWT, err := utils.GetUserInformationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return &generated.UserInformation{
		Os:            authJWT.Os,
		Token:         authJWT.Token,
		DeviceName:    authJWT.DeviceName,
		DeviceID:      authJWT.DeviceID,
		AppID:         authJWT.AppID,
		Vip:           authJWT.Vip,
		ComboID:       authJWT.ComboID,
		DaysLeft:      authJWT.DaysLeft,
		RemainingFlow: authJWT.RemainingFlow,
	}, nil
}
