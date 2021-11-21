package config

import "time"

const (
	ContextBearerToken = "BearerToken"
)

const (
	DefaultClientTimeout = time.Second * 5
	DefaultCtxTimeout    = time.Second * 5
	DefaultIdleTimeout   = time.Second * 60
	DefaultReadTimeout   = time.Second * 15
	DefaultTokenLeeway   = time.Second * 30
	DefaultWriteTimeout  = time.Second * 15
)

const (
	MimeApplicationJson               = "application/json"
	MimeApplicationXWwwFormUrlencoded = "application/x-www-form-urlencoded"
	MimeTextHtml                      = "text/html"
)
