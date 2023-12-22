package generated

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Resolver struct{}

// // foo
func (r *mutationResolver) Healthcheck(ctx context.Context) (string, error) {
	panic("not implemented")
}

// // foo
func (r *mutationResolver) UploadFile(ctx context.Context, file graphql.Upload) (string, error) {
	panic("not implemented")
}

// // foo
func (r *mutationResolver) SendEmail(ctx context.Context, input *SendEmailRequest) (*SendEmailResponse, error) {
	panic("not implemented")
}

// // foo
func (r *mutationResolver) UserRegistration(ctx context.Context, input *UserRegistration) (*AuthPayload, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) Healthcheck(ctx context.Context) (string, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) Now(ctx context.Context) (*timestamppb.Timestamp, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) App(ctx context.Context, appID string) (*AppInfo, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) Combos(ctx context.Context) (*Combos, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) Nodes(ctx context.Context) (*Nodes, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) NodeToken(ctx context.Context, nodeID string, mountToken *string) (*NodeToken, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) User(ctx context.Context) (*UserInformation, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) Captcha(ctx context.Context) (*Captcha, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
