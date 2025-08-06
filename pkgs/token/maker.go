package token

type Maker interface {
	CreateToken(payload *Payload) (string, error)
	VerifyUserToken(tokenString string) (*Payload, error)
}
