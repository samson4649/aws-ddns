package aws_ddns

import (
  "net/http"
  "context"
)

type Client struct {
  ctx  *context.Context
  http *http.Client
  auth *Auth
}

func NewClient( ctx *context.Context, auth *Auth ) ( *Client, error ) {
  return &Client{
    ctx: ctx,
    http: &http.Client{},
    auth: auth,
  }, nil
}

