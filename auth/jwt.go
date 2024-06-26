package auth

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go_todo_app/clock"
	"go_todo_app/entity"
	"net/http"
	"time"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

type userIDKey struct{}
type roleKey struct{}

func NewJWTer(store Store, clocker clock.Clocker) (*JWTer, error) {
	jwt := &JWTer{Store: store}
	privateKey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}
	publicKey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}

	jwt.PrivateKey = privateKey
	jwt.PublicKey = publicKey
	jwt.Clocker = clocker

	return jwt, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (j *JWTer) GenerateToken(ctx context.Context, user entity.User) ([]byte, error) {
	token, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/kyeongminlee/go_todo_app`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, user.Role).
		Claim(UserNameKey, user.Name).
		Build()

	if err != nil {
		return nil, fmt.Errorf("GenerateToken: failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx, token.JwtID(), user.ID); err != nil {
		return nil, err
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, request *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		request,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}

	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired: %w", token.JwtID(), err)
	}
	return token, nil
}

func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserID(ctx context.Context) (entity.UserID, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return id, ok
}

func SetRole(ctx context.Context, token jwt.Token) context.Context {
	get, ok := token.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}
	return context.WithValue(ctx, roleKey{}, get)
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

func (j *JWTer) FillContext(request *http.Request) (*http.Request, error) {
	token, err := j.GetToken(request.Context(), request)
	if err != nil {
		return nil, err
	}
	uid, err := j.Store.Load(request.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}

	ctx := SetUserID(request.Context(), uid)
	ctx = SetRole(ctx, token)
	clone := request.Clone(ctx)
	return clone, nil
}

func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	return role == "admin"
}
