package internetgolf

import (
	"fmt"
	"strings"
)

func getPermissionsForRequest(remoteAddr string, authHeader string) (Permissions, error) {
	if l := (LocalReqAuthChecker{}); l.setReqData(remoteAddr, authHeader) {
		return &l, nil
	} else if g := (GithubAuthChecker{}); g.setReqData(remoteAddr, authHeader) {
		return &g, nil
	}
	return nil, fmt.Errorf("could not check auth for header value \"%s\"", authHeader)
}

type Permissions interface {
	// returns false if the given concrete implementation of Permissions is not
	// suitable for the given request data
	setReqData(remoteAddr string, authHeader string) bool
	canCreateDeployment() bool
	canModifyDeployment(d *Deployment) bool
	canViewDeployment(d *Deployment) bool
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
	return remoteAddr == "127.0.0.1" || strings.HasPrefix(remoteAddr, "127.0.0.1:")
}
func (l *LocalReqAuthChecker) canModifyDeployment(_ *Deployment) bool {
	return true
}
func (l *LocalReqAuthChecker) canCreateDeployment() bool {
	return true
}
func (l *LocalReqAuthChecker) canViewDeployment(_ *Deployment) bool {
	return true
}
