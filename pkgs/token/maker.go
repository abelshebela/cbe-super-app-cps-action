package token

type Maker interface {
	VerifyUserToken(tokenString string) (*Payload, error)
}
