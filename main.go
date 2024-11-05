package caddy_req_id

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/google/uuid"
)

func init() {
	caddy.RegisterModule(ReqID{})
	httpcaddyfile.RegisterHandlerDirective("req_id", parseCaddyfile)
}

type ReqID struct {
	Enabled bool `json:"enabled,omitempty"`
}

func (ReqID) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.req_id",
		New: func() caddy.Module { return new(ReqID) },
	}
}

func enhancedUUID() string {
	uuid := uuid.New().String()
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(append([]byte(uuid), randomBytes...))
	return hex.EncodeToString(hash[:])
}

func (m ReqID) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	if m.Enabled {
		reqID := enhancedUUID()
		r.Header.Set("Req-ID", reqID)
		w.Header().Set("Req-ID", reqID)
	}
	return next.ServeHTTP(w, r)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u ReqID
	if !h.Next() {
		return nil, h.ArgErr()
	}
	remainingArgs := h.RemainingArgs()
	if len(remainingArgs) > 0 {
		u.Enabled = (remainingArgs[0] == "true")
	}
	return u, nil
}

var (
	_ caddyhttp.MiddlewareHandler = (*ReqID)(nil)
)

