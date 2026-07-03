package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	OpenID string `json:"open_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

func NewJWT(secret string, expiration time.Duration, issuer string) *JWT {
	return &JWT{
		secret:     []byte(secret),
		expiration: expiration,
		issuer:     issuer,
	}
}

func (j *JWT) GenerateToken(userID int64, openID, phone string) (string, error) {
	claims := Claims{
		UserID: userID,
		OpenID: openID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
