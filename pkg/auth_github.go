package internetgolf

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// implements the interface `Permissions`
type GithubAuthChecker struct {
	oidcToken GitHubOIDCToken
	Db        Db
}

func (g *GithubAuthChecker) setReqData(_remoteAddr string, authHeader string) bool {
	headerComps := strings.Split(authHeader, " ")
	if len(headerComps) != 2 || headerComps[0] != "GithubOIDC" {
		return false
	}
	tokenData, tokenError := ParseGithubOidcToken(headerComps[1])
	if tokenError != nil {
		panic(tokenError)
	}
	g.oidcToken = tokenData
	return true
}

func (g *GithubAuthChecker) userHasFullPermissions() bool {
	externalUserId := g.oidcToken.ActorID
	externalUser, err := g.Db.GetExternalUser(externalUserId)
	if err != nil {
		return false
	}
	return externalUser.FullPermissions
}

func (g *GithubAuthChecker) canModifyDeployment(d *Deployment) bool {
	if g.userHasFullPermissions() {
		return true
	}
	if d.ExternalSourceType != Github {
		return false
	}
	repo := g.oidcToken.Repository
	branch := ""
	if g.oidcToken.RefType == "branch" {
		branchNameLocation := strings.LastIndex(g.oidcToken.Ref, "/")
		if branchNameLocation != -1 {
			branch = g.oidcToken.Ref[branchNameLocation+1:]
		}
	}

	return (d.ExternalSource == repo || d.ExternalSource == repo+"#"+branch)
}

func (g *GithubAuthChecker) canCreateDeployment() bool {
	return g.userHasFullPermissions()
}

func (g *GithubAuthChecker) canViewDeployment(d *Deployment) bool {
	if g.canModifyDeployment(d) {
		return true
	}
	return g.userHasFullPermissions()
}

func (g *GithubAuthChecker) canAddUser() bool {
	return g.userHasFullPermissions()
}

// example payload:
//
//	{
//	  "actor": "toBeOfUse",
//	  "actor_id": "49729978",
//	  "aud": "https://github.com/toBeOfUse",
//	  "base_ref": "",
//	  "event_name": "workflow_dispatch",
//	  "exp": 1756127261,
//	  "head_ref": "",
//	  "iat": 1756105661,
//	  "iss": "https://token.actions.githubusercontent.com",
//	  "job_workflow_ref": "toBeOfUse/internet-golf/.github/workflows/oidc-test.yml@refs/heads/main",
//	  "job_workflow_sha": "54fd8fbf5de6050880e24d97dd9870942c04f258",
//	  "jti": "975c1557-fd94-4a8b-ad76-89848e46cbfc",
//	  "nbf": 1756105361,
//	  "ref": "refs/heads/main",
//	  "ref_protected": "false",
//	  "ref_type": "branch",
//	  "repository": "toBeOfUse/internet-golf",
//	  "repository_id": "1034463833",
//	  "repository_owner": "toBeOfUse",
//	  "repository_owner_id": "49729978",
//	  "repository_visibility": "public",
//	  "run_attempt": "1",
//	  "run_id": "17201782829",
//	  "run_number": "2",
//	  "runner_environment": "github-hosted",
//	  "sha": "54fd8fbf5de6050880e24d97dd9870942c04f258",
//	  "sub": "repo:toBeOfUse/internet-golf:ref:refs/heads/main",
//	  "workflow": "Print and Post ID Token Variables",
//	  "workflow_ref": "toBeOfUse/internet-golf/.github/workflows/oidc-test.yml@refs/heads/main",
//	  "workflow_sha": "54fd8fbf5de6050880e24d97dd9870942c04f258"
//	}
type GitHubOIDCToken struct {
	Actor                string `json:"actor"`
	ActorID              string `json:"actor_id"`
	Aud                  string `json:"aud"`
	BaseRef              string `json:"base_ref"`
	EventName            string `json:"event_name"`
	Exp                  int64  `json:"exp"`
	HeadRef              string `json:"head_ref"`
	Iat                  int64  `json:"iat"`
	Iss                  string `json:"iss"`
	JobWorkflowRef       string `json:"job_workflow_ref"`
	JobWorkflowSha       string `json:"job_workflow_sha"`
	Jti                  string `json:"jti"`
	Nbf                  int64  `json:"nbf"`
	Ref                  string `json:"ref"`
	RefProtected         string `json:"ref_protected"`
	RefType              string `json:"ref_type"`
	Repository           string `json:"repository"`
	RepositoryID         string `json:"repository_id"`
	RepositoryOwner      string `json:"repository_owner"`
	RepositoryOwnerID    string `json:"repository_owner_id"`
	RepositoryVisibility string `json:"repository_visibility"`
	RunAttempt           string `json:"run_attempt"`
	RunID                string `json:"run_id"`
	RunNumber            string `json:"run_number"`
	RunnerEnvironment    string `json:"runner_environment"`
	Sha                  string `json:"sha"`
	Sub                  string `json:"sub"`
	Workflow             string `json:"workflow"`
	WorkflowRef          string `json:"workflow_ref"`
	WorkflowSha          string `json:"workflow_sha"`
}

// newJWKSet creates an auto-refreshing key set to validate JWT signatures.
// borrowed from example https://huma.rocks/how-to/oauth2-jwt/?h=ctx#huma-auth-middleware
func newJWKSet(jwkUrl string) (jwk.Set, error) {
	jwkCache := jwk.NewCache(context.Background())

	// register a minimum refresh interval for this URL.
	// when not specified, defaults to Cache-Control and similar resp headers
	err := jwkCache.Register(jwkUrl, jwk.WithMinRefreshInterval(10*time.Minute))
	if err != nil {
		return nil, errors.New("failed to register jwk location")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch once on application startup
	_, err = jwkCache.Refresh(ctx, jwkUrl)
	if err != nil {
		return nil, err
	}
	// create the cached key set
	return jwk.NewCachedSet(jwkCache, jwkUrl), nil
}

func ParseGithubOidcToken(token string) (GitHubOIDCToken, error) {
	keySet, keySetErr := newJWKSet("https://token.actions.githubusercontent.com/.well-known/jwks")
	if keySetErr != nil {
		return GitHubOIDCToken{}, keySetErr
	}

	_, err := jwt.ParseString(token,
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
		jwt.WithAudience("internet-golf"),
	)

	if err != nil {
		return GitHubOIDCToken{}, err
	}

	// after validating the token, accessing the raw payload data is somehow
	// the easiest way to get a struct out of it
	rawJson, rawJsonErr := base64.RawStdEncoding.DecodeString(strings.Split(token, ".")[1])
	if rawJsonErr != nil {
		return GitHubOIDCToken{}, rawJsonErr
	}
	var tokenData GitHubOIDCToken
	json.Unmarshal(rawJson, &tokenData)

	return tokenData, nil
}
