package aws_ddns

import (

)

type Auth struct {
  AccessKey string
  SecretKey string
}

func NewAuth( access, secret string) (*Auth) {
  // add test to validate auth here
  return &Auth{
    access,
    secret,
  }
}
