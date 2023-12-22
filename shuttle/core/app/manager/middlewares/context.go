package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rs/xid"
	"google.dev/google/shuttle/core/app/manager/pkg/enum"
	"google.dev/google/shuttle/core/app/manager/utils"
)

// Context  get user from jwt and put user into ctx
func Context() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := xid.New().String()
			ctx := context.WithValue(r.Context(), enum.RequestId, reqID)

			ctx = context.WithValue(ctx, enum.RequestReceivedAtCtxKey, time.Now())

			// user agent
			userAgent := r.Header.Get("User-Agent")
			ctx = context.WithValue(ctx, enum.UserAgentCtxKey, userAgent)

			// ip
			reqIP := r.Header.Get("X-Forwarded-For")
			if reqIP == "" {
				reqIP = r.Header.Get("X-Real-Ip")
			}
			if reqIP == "" {
				index := strings.Index(r.RemoteAddr, ":")
				if index != -1 {
					reqIP = r.RemoteAddr[:index]
				} else {
					reqIP = r.RemoteAddr
				}
			}
			ctx = context.WithValue(ctx, enum.ReqIP, reqIP)

			// token
			tokenString, _ := utils.GetTokenFromHeader(r.Header)
			if tokenString != "" {
				ctx = context.WithValue(ctx, enum.TokenCtxKey, tokenString)
			}

			r = r.WithContext(ctx)

			log.Println(r.URL.String())
			next.ServeHTTP(w, r)
		})
	}
}

func Safety() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			salt := utils.Md5Encode(fmt.Sprintf("%d-9776e538-59ba-473f-8ccf-1d72031e360f", time.Now().UnixMilli()/10000))
			reqSalt := r.Header.Get("salt")
			if salt == reqSalt {
				next.ServeHTTP(w, r)
				return
			}

			respData := map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"message": "500 server exception",
						"path":    []string{"main"},
					},
				},
			}

			marshal, err := json.Marshal(respData)
			if err == nil {
				w.Write(marshal)
			}
			//w.WriteHeader(500)
		})
	}
}
