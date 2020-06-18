package aws_ddns

type Auth struct {
  AccessKey string
  SecretKey string
}

func NewAuth( access, secret string) (*Auth) {
  return &Auth{
    access,
    secret,
  }
}
