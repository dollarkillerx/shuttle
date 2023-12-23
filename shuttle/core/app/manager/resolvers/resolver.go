package resolvers

import (
	"context"
	"github.com/patrickmn/go-cache"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/afocus/captcha"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"google.dev/google/shuttle/core/app/manager/generated"
	"google.dev/google/shuttle/core/app/manager/storage"
	"google.dev/google/socks5_discovery/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Resolver ...
type Resolver struct {
	Storage storage.Interface
	cache   *cache.Cache
	captcha *captcha.Captcha

	socks5DiscoveryClient proto.Socks5DiscoveryClient
}

func NewResolver(db storage.Interface, socks5DiscoveryClient proto.Socks5DiscoveryClient) *Resolver {
	return &Resolver{
		Storage:               db,
		cache:                 cache.New(15*time.Minute, 30*time.Minute),
		captcha:               captchaInit(),
		socks5DiscoveryClient: socks5DiscoveryClient,
	}
}

func captchaInit() (cca *captcha.Captcha) {
	cca = captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	cca.SetFont("./static/comic.ttf")
	// 设置验证码大小
	cca.SetSize(150, 64)
	// 设置干扰强度
	cca.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cca.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cca.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	return
}

// Mutation is the root mutation resolver
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

// Query is the root query resolver
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Healthcheck(ctx context.Context) (string, error) {
	return "ack", nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Healthcheck(ctx context.Context) (string, error) {
	return "ack", nil
}

func (r *queryResolver) Now(ctx context.Context) (*timestamppb.Timestamp, error) {
	return &timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
	}, nil
}

// field resolver
type fileInfoResolver struct{ *Resolver }

// UploadFile ...
func (r *mutationResolver) UploadFile(ctx context.Context, file graphql.Upload) (string, error) {
	// 检查上传的文件是否符合要求
	maxSize := int64(8 * 1024 * 1024) // 8MB
	if file.Size > maxSize {
		return "", errors.New("file size too large")
	}

	// 生成新的文件名
	newFileName := xid.New().String() + filepath.Ext(file.Filename)

	// 打开文件，准备写入
	targetFile, err := os.OpenFile("./static/"+newFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	// 从上传的文件中读取数据，并写入到打开的文件中
	fileData, err := ioutil.ReadAll(file.File)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(targetFile, strings.NewReader(string(fileData)))
	if err != nil {
		return "", err
	}

	// 返回新的文件名
	return newFileName, nil
}
