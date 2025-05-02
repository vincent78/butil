package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	Issuer = "utilInfo"
	Secret = "utilInfo"
)

const (
	TokenTypeBearer = "bearer"
)

const (
	ClaimsTypeRefresh = "RefreshToken"
	ClaimsTypeAccess  = "AccessToken"
)

type Payload struct {
	UserID        string
	Username      string
	NickName      string
	UserStatus    int32
	UserType      int32
	AccountType   int32
	NamespaceIds  []string
	DownloadCount int
	Admin         string
}
type CustomClaims struct {
	Type string
	*Payload
	jwt.RegisteredClaims
}

type CustomToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresTime  time.Time `json:"expires_time"`
	RefreshToken string    `json:"refresh_token"`
	Payload      *Payload  `json:"payload"`
}

func generate(claimsType string, payload *Payload, expiresTime time.Duration) (string, time.Time, error) {
	nowTime := time.Now()
	expiresT := nowTime.Add(expiresTime)
	claims := CustomClaims{
		Type:    claimsType,
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			ExpiresAt: jwt.NewNumericDate(expiresT),
			NotBefore: &jwt.NumericDate{},
			IssuedAt:  jwt.NewNumericDate(nowTime),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	newToken, err := t.SignedString([]byte(Secret))
	return newToken, expiresT, err
}

func VerifyToken(token string, claimsType string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	}, jwt.WithJSONNumber())
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Type != claimsType {
		return nil, errors.New("invalid " + claimsType)
	}

	return claims, nil
}

func GenerateToken(payload *Payload, accessTokenExpiresTime, refreshTokenExpiresTime time.Duration) (*CustomToken, error) {

	// 生成accessToken
	// 存入payload
	accessToken, expiresTime, err := generate(ClaimsTypeAccess, payload, accessTokenExpiresTime)
	if err != nil {
		return nil, err
	}

	// 生成refreshToken
	// 存入payload
	refreshToken, _, err := generate(ClaimsTypeRefresh, payload, refreshTokenExpiresTime)
	if err != nil {
		return nil, err
	}

	return &CustomToken{TokenType: TokenTypeBearer, ExpiresTime: expiresTime, AccessToken: accessToken, RefreshToken: refreshToken, Payload: payload}, nil
}
