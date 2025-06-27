package utils

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	constants "github.com/nelsonin-research-org/clenz-auth/const"
	"github.com/nelsonin-research-org/clenz-auth/globals"
	"github.com/nelsonin-research-org/clenz-auth/models/appschema"
)

func GenerateJWTClaims(data *appschema.JwtData, ttl time.Duration) jwt.MapClaims {
	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["email"] = data.Email
	claims["first_name"] = data.FirstName
	claims["last_name"] = data.LastName
	claims["id"] = data.ID
	claims["type"] = data.Type
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["token_type"] = constants.LONG_LIVE_TOKEN
	claims["jwt_created"] = now.Unix()
	return claims
}

func GenerateToken(ttl time.Duration, data *appschema.JwtData) (map[string]string, error) {
	accessClaims := GenerateJWTClaims(data, ttl)
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims).SignedString(globals.AppKeys.PrivateKey)
	if err != nil {
		fmt.Println("Error generating access token: ", err)
		return nil, err
	}

	refreshClaims := GenerateRefreshTokenClaims(ttl+7*24*time.Hour, data.ID, data.Email)
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(globals.AppKeys.PrivateKey)
	if err != nil {
		fmt.Println("Error generating refresh token: ", err)
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return tokens, nil
}

func GenerateTempToken(ttl time.Duration, data *appschema.TempJwtData) (map[string]string, error) {
	quickClaims := GenerateTempTokenClaims(data, ttl)
	quickToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, quickClaims).SignedString(globals.AppKeys.PrivateKey)
	if err != nil {
		fmt.Println("Error generating access token: ", err)
		return nil, err
	}

	token := map[string]string{
		"token": quickToken,
	}

	return token, nil
}

func GenerateTempTokenClaims(data *appschema.TempJwtData, ttl time.Duration) jwt.MapClaims {
	now := time.Now().UTC()
	claims := make(jwt.MapClaims)

	claims["email"] = data.Email
	claims["token_type"] = constants.TEMP_TOKEN
	claims["exp"] = now.Add(ttl).Unix()
	claims["jwt_created"] = now.Unix()

	return claims
}

func ParseToken(token string) (jwt.MapClaims, error) {

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return globals.AppKeys.PublicKeyPem, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, err
	}

	return claims, nil
}

func GenerateRefreshTokenClaims(ttl time.Duration, userId, email string) jwt.MapClaims {
	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["id"] = userId
	claims["email"] = FormatStringToLowerCase(email)
	claims["exp"] = now.Add(ttl).Unix()
	claims["jwt_created"] = now.Unix()
	claims["token_type"] = constants.REFRESH_TOKEN
	return claims
}
