package auth_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/go-modules/authentication/authenticationtest"
	"github.com/ankorstore/go-modules/httpserversecurity"
	"github.com/ankorstore/yokai-contrib/fxe2e/auth"
	"github.com/onsi/gomega"
)

// decodeToken pulls the bearer token off the request, base64-decodes it and
// returns the parsed JSON claims. The adapters render a JSON token and set it as
// "Authorization: Bearer <base64(token)>".
func decodeToken(g gomega.Gomega, req *http.Request) map[string]any {
	header := req.Header.Get(httpserversecurity.AuthorizationHTTPHeaderName)
	g.Expect(header).To(gomega.HavePrefix("Bearer "), "missing bearer Authorization header")

	raw, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(header, "Bearer "))
	g.Expect(err).NotTo(gomega.HaveOccurred(), "token is not valid base64")

	var claims map[string]any
	g.Expect(json.Unmarshal(raw, &claims)).To(gomega.Succeed(), "token is not valid JSON")

	return claims
}

func aks(g gomega.Gomega, claims map[string]any) map[string]any {
	a, ok := claims["aks"].(map[string]any)
	g.Expect(ok).To(gomega.BeTrue(), "token has no aks block")

	return a
}

func newRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/", nil)
}

// TestAnkorstore_Principals drives every principal type through the default
// dispatch (so it covers both routing and the adapter) and asserts the rendered
// token carries the expected idp, entity and principal-specific claims.
func TestAnkorstore_Principals(t *testing.T) {
	dispatch := auth.Ankorstore()

	cases := []struct {
		name   string
		raw    string
		idp    string
		entity string
		assert func(g gomega.Gomega, claims map[string]any)
	}{
		{
			name:   "brand",
			raw:    `{"type":"brand","accountUuid":"brand-1","accountEmail":"b@example.com"}`,
			idp:    "aks_user",
			entity: "user",
			assert: func(g gomega.Gomega, claims map[string]any) {
				acc := aks(g, claims)["account"].(map[string]any)
				g.Expect(acc["type"]).To(gomega.Equal("brand"))
				g.Expect(acc["id"]).To(gomega.Equal("brand-1"))
				g.Expect(acc["email"]).To(gomega.Equal("b@example.com"))
			},
		},
		{
			name:   "retailer",
			raw:    `{"type":"retailer","accountUuid":"retailer-1"}`,
			idp:    "aks_user",
			entity: "user",
			assert: func(g gomega.Gomega, claims map[string]any) {
				acc := aks(g, claims)["account"].(map[string]any)
				g.Expect(acc["type"]).To(gomega.Equal("retailer"))
				g.Expect(acc["id"]).To(gomega.Equal("retailer-1"))
			},
		},
		{
			name:   "admin",
			raw:    `{"type":"admin","roles":["customer-care"],"permissions":["leads.read"]}`,
			idp:    "aks_admin",
			entity: "admin",
			assert: func(g gomega.Gomega, claims map[string]any) {
				adminBlk := aks(g, claims)["admin"].(map[string]any)
				g.Expect(adminBlk["roles"]).To(gomega.ConsistOf("customer-care"))
				g.Expect(adminBlk["permissions"]).To(gomega.ConsistOf("leads.read"))
			},
		},
		{
			name:   "guest",
			raw:    `{"type":"guest","clientId":"client-x"}`,
			idp:    "aks_guest",
			entity: "guest",
			assert: func(g gomega.Gomega, claims map[string]any) {
				g.Expect(claims["cid"]).To(gomega.Equal("client-x"))
			},
		},
		{
			name:   "machine",
			raw:    `{"type":"machine","clientId":"client-y"}`,
			idp:    "aks_machine",
			entity: "machine",
			assert: func(g gomega.Gomega, claims map[string]any) {
				g.Expect(claims["cid"]).To(gomega.Equal("client-y"))
			},
		},
		{
			name:   "impersonation",
			raw:    `{"type":"impersonation","accountType":"brand","accountUuid":"acc-1","impersonatorUuid":"imp-1"}`,
			idp:    "aks_imp",
			entity: "user",
			assert: func(g gomega.Gomega, claims map[string]any) {
				acc := aks(g, claims)["account"].(map[string]any)
				g.Expect(acc["type"]).To(gomega.Equal("brand"))
				g.Expect(acc["id"]).To(gomega.Equal("acc-1"))

				imp := aks(g, claims)["imp"].(map[string]any)
				g.Expect(imp["sub"]).To(gomega.Equal("imp-1"))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := gomega.NewWithT(t)

			req := newRequest()
			dispatch(t, req, json.RawMessage(c.raw))

			claims := decodeToken(g, req)
			g.Expect(claims["idp"]).To(gomega.Equal(c.idp))
			g.Expect(aks(g, claims)["entity"]).To(gomega.Equal(c.entity))

			if c.assert != nil {
				c.assert(g, claims)
			}
		})
	}
}

// TestAnkorstore_DefaultsApplied checks that a bare principal fixture relies on
// the go-modules defaults rather than producing empty claims.
func TestAnkorstore_DefaultsApplied(t *testing.T) {
	g := gomega.NewWithT(t)

	req := newRequest()
	auth.Ankorstore()(t, req, json.RawMessage(`{"type":"brand"}`))

	acc := aks(g, decodeToken(g, req))["account"].(map[string]any)
	g.Expect(acc["id"]).To(gomega.Equal(authenticationtest.DefaultBrandAccountUUID))
	g.Expect(acc["email"]).To(gomega.Equal(authenticationtest.DefaultBrandAccountEmail))
}

// TestAnkorstore_Anonymous asserts no Authorization header is set for an absent
// or "none" principal.
func TestAnkorstore_Anonymous(t *testing.T) {
	dispatch := auth.Ankorstore()

	for _, raw := range []string{``, `{"type":"none"}`, `{}`} {
		req := newRequest()
		dispatch(t, req, json.RawMessage(raw))

		g := gomega.NewWithT(t)
		g.Expect(req.Header.Get(httpserversecurity.AuthorizationHTTPHeaderName)).To(gomega.BeEmpty(),
			"expected no Authorization header for raw %q", raw)
	}
}

// TestAnkorstore_Origin checks the orthogonal origin field sets the SPIFFE
// client-cert header, independently of (and composed with) the principal.
func TestAnkorstore_Origin(t *testing.T) {
	dispatch := auth.Ankorstore()

	t.Run("internal on anonymous", func(t *testing.T) {
		g := gomega.NewWithT(t)

		req := newRequest()
		dispatch(t, req, json.RawMessage(`{"type":"none","origin":"internal"}`))

		cert := req.Header.Get(httpserversecurity.ForwardedClientCertHTTPHeaderName)
		g.Expect(cert).NotTo(gomega.BeEmpty())
		g.Expect(cert).NotTo(gomega.ContainSubstring(httpserversecurity.ExternalSPIFFEURI))
	})

	t.Run("external composed with brand", func(t *testing.T) {
		g := gomega.NewWithT(t)

		req := newRequest()
		dispatch(t, req, json.RawMessage(`{"type":"brand","accountUuid":"brand-1","origin":"external"}`))

		g.Expect(req.Header.Get(httpserversecurity.ForwardedClientCertHTTPHeaderName)).
			To(gomega.ContainSubstring(httpserversecurity.ExternalSPIFFEURI))
		// The principal is still applied alongside the origin.
		g.Expect(req.Header.Get(httpserversecurity.AuthorizationHTTPHeaderName)).To(gomega.HavePrefix("Bearer "))
	})
}
