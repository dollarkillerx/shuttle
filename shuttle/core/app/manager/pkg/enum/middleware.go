package enum

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

// contextKey ...
const (
	TokenCtxKey             ContextKey = "Token"
	AuthorizationCtxKey     ContextKey = "Authorization"
	UserAgentCtxKey         ContextKey = "UserAgent"
	RequestId               ContextKey = "requestID"
	RequestReceivedAtCtxKey ContextKey = "ReqReceivedAt"
	ReqIP                   ContextKey = "ReqIP"
)
