package consts

type ctxKey string

// ? does it better to be used as enum

const (
	SidCtxKey       ctxKey = ctxKey("sid")
	TokenCtxKey     ctxKey = ctxKey("token")
	LoginPathCtxKey ctxKey = ctxKey("login-path")
)
