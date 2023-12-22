package resolvers

import (
	"context"

	"google.dev/google/shuttle/core/app/manager/generated"
	"google.dev/google/shuttle/core/app/manager/pkg/models"
)

// App ...
func (r *queryResolver) App(ctx context.Context, appID string) (*generated.AppInfo, error) {
	var app models.App
	err := r.Storage.DB().Model(&models.App{}).Where("app_id = ?", appID).First(&app).Error
	if err != nil {
		return nil, err
	}

	return &generated.AppInfo{
		AppID:              app.AppID,
		AppVersion:         app.AppVersion,
		MinimumVersion:     app.MinimumVersion,
		State:              app.State,
		ErrorNotification:  app.ErrorNotification,
		NormalNotification: app.NormalNotification,
	}, nil
}
