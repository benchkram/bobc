package authenticator

import (
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/benchkram/errz"
	"github.com/labstack/echo/v4"
)

type Authenticator struct {
	apiKey []byte
}

func New(apiKey []byte) *Authenticator {
	return &Authenticator{
		apiKey: apiKey,
	}
}

func (a *Authenticator) Authenticate(ctx echo.Context) (err error) {
	defer errz.Recover(&err)

	token, err := a.extractTokenFromRequest(ctx)
	errz.Fatal(err)

	return a.extractUserFromToken(token)
}

func (a *Authenticator) extractUserFromToken(token string) error {
	if subtle.ConstantTimeCompare([]byte(token), a.apiKey) != 1 {
		return errors.New("invalid token")
	}

	return nil
}

func (a *Authenticator) extractTokenFromRequest(ctx echo.Context) (token string, err error) {
	defer errz.Recover(&err)

	r := ctx.Request()

	auth := r.Header.Get("Authorization")
	if auth == "" {
		errz.Fatal(errors.New("invalid auth token"))
	}

	parts := strings.Split(auth, "Bearer")
	if len(parts) != 2 {
		errz.Fatal(errors.New("invalid auth token"))
	}

	return strings.TrimSpace(parts[1]), nil
}
