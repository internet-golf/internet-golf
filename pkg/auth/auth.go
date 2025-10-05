package auth

import (
	"fmt"
	"net"
	"strings"

	"github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager struct {
	Db db.Db
}

func (a *AuthManager) GetPermissionsForRequest(remoteAddr string, authHeader string) (Permissions, error) {
	if l := (LocalReqAuthChecker{}); l.setReqData(remoteAddr, authHeader) {
		fmt.Println("automatically trusting request from " + remoteAddr)
		return &l, nil
	} else if g := (GithubAuthChecker{Db: a.Db}); g.setReqData(remoteAddr, authHeader) {
		return &g, nil
	} else if b := (BearerTokenAuthChecker{Db: a.Db}); b.setReqData(remoteAddr, authHeader) {
		return &b, nil
	}
	return nil, fmt.Errorf("could not check auth for header value \"%s\"", authHeader)
}

func (a *AuthManager) RegisterExternalUser(e db.ExternalUser) {
	a.Db.SaveExternalUser(e)
}

func (a *AuthManager) CreateBearerToken(fullPermissions bool) (string, error) {
	return (&BearerTokenAuthChecker{Db: a.Db}).CreateBearerToken(fullPermissions)
}

type Permissions interface {
	// returns false if the given concrete implementation of Permissions is not
	// suitable for the given request data
	setReqData(remoteAddr string, authHeader string) bool
	CanCreateDeployment() bool
	CanModifyDeployment(d *db.Deployment) bool
	CanViewDeployment(d *db.Deployment) bool
	// can add external users and bearer tokens
	CanCreateCredentials() bool
}

// if a request comes from the same machine as the server (i.e. comes from
// 127.0.0.1), this lets it do whatever it wants.
//
// this is similar to how you can access caddy's admin api from the same machine
// of it and just do whatever.
//
// implements the interface `Permissions`.
type LocalReqAuthChecker struct{}

func (l *LocalReqAuthChecker) setReqData(remoteAddr string, authHeader string) bool {
	// this allows requests from localhost or from other entities that have the
	// ability to give themselves the hostname "golf-client" on the local
	// network (which pretty much just means docker containers in the same
	// virtual network)
	serverAddr, _ := net.LookupIP("golf-client")
	return (remoteAddr == "127.0.0.1" || strings.HasPrefix(remoteAddr, "127.0.0.1:") ||
		(len(serverAddr) > 0 && strings.HasPrefix(remoteAddr, serverAddr[0].String()+":")))
}
func (l *LocalReqAuthChecker) CanModifyDeployment(_ *db.Deployment) bool {
	return true
}
func (l *LocalReqAuthChecker) CanCreateDeployment() bool {
	return true
}
func (l *LocalReqAuthChecker) CanViewDeployment(_ *db.Deployment) bool {
	return true
}
func (l *LocalReqAuthChecker) CanCreateCredentials() bool {
	return true
}

// provides bearer token-based authorization. implements the Permissions interface.
type BearerTokenAuthChecker struct {
	Db    db.Db
	token db.BearerToken
}

func (b *BearerTokenAuthChecker) CreateBearerToken(fullPermissions bool) (string, error) {
	var token, id string
	for {
		token, id = utils.GetRandomToken()
		existing, err := b.Db.GetBearerToken(id)
		if err != nil && len(existing.Id) == 0 {
			break
		}
	}
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	if err != nil {
		return "", err
	}
	b.Db.SaveBearerToken(db.BearerToken{Id: id, TokenHash: tokenHash, FullPermissions: fullPermissions})
	return id + "." + token, nil
}

func (b *BearerTokenAuthChecker) setReqData(remoteAddr string, authHeader string) bool {
	comps := strings.Split(authHeader, " ")
	if len(comps) != 2 || comps[0] != "Bearer" {
		return false
	}
	tokenWithId := comps[1]
	tokenComps := strings.Split(tokenWithId, ".")
	if len(tokenComps) != 2 {
		// should this be a panic? what would the implications be?
		panic("invalid token value received; it should have the format [id].[content]")
	}
	id := tokenComps[0]
	token := tokenComps[1]
	tokenStruct, tokenErr := b.Db.GetBearerToken(id)
	if tokenErr != nil {
		panic("could not find bearer token matching input")
	}
	compareErr := bcrypt.CompareHashAndPassword(tokenStruct.TokenHash, []byte(token))
	if compareErr != nil {
		panic("token content is not a match")
	}
	b.token = tokenStruct
	return true
}
func (b *BearerTokenAuthChecker) CanModifyDeployment(_ *db.Deployment) bool {
	return b.token.FullPermissions
}
func (b *BearerTokenAuthChecker) CanCreateDeployment() bool {
	return b.token.FullPermissions
}
func (b *BearerTokenAuthChecker) CanViewDeployment(_ *db.Deployment) bool {
	return b.token.FullPermissions
}
func (b *BearerTokenAuthChecker) CanCreateCredentials() bool {
	return b.token.FullPermissions
}
