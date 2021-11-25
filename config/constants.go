package config

import "time"

const (
	ContextBearerToken = "BearerToken"
)

const (
	DefaultClientTimeout = time.Second * 10
	DefaultCtxTimeout    = time.Second * 10
	DefaultIdleTimeout   = time.Second * 60
	DefaultReadTimeout   = time.Second * 10
	DefaultTokenLeeway   = time.Second * 30
	DefaultWriteTimeout  = time.Second * 10
)

const (
	MimeApplicationJson               = "application/json"
	MimeApplicationXWwwFormUrlencoded = "application/x-www-form-urlencoded"
	MimeTextHtml                      = "text/html"
)
