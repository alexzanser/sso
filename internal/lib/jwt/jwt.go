package jwt

import (
	"time"

	"github.com/alexzanser/sso/internal/domain/models"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, ttl time.Duration) (string, error) {
	token := jwt5.New(jwt5.SigningMethodHS256)

	claims := token.Claims.(jwt5.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.ID

	signedToken, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
