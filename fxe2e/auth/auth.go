// Package auth provides request-authentication for the e2e engine. It knows the
// Ankorstore principal types out of the box: a case declares its principal in
// 2_request.json and the framework applies the matching token — no wiring in the
// consumer's Boot or test entrypoint.
//
//	"auth": {"type": "brand", "accountUuid": "..."}
//
// e2e.Runner defaults its Auth to Ankorstore(), so cases authenticate purely by
// declaring a "type". The supported types and their (all optional — go-modules
// fills sensible defaults) fields are:
//
//	guest          clientId
//	machine        clientId
//	admin          clientId, entityUuid, roles, permissions
//	brand          clientId, entityUuid, accountUuid, accountEmail
//	retailer       clientId, entityUuid, accountUuid, accountEmail
//	impersonation  accountType (brand|retailer), accountUuid, accountEmail,
//	               entityUuid, clientId, impersonatorUuid, impersonatorRoles,
//	               impersonatorPermissions
//
// An absent or "none" type sends the request anonymously. An optional "origin"
// field ("internal" or "external") is applied on top of any type, including
// anonymous, to set the request's SPIFFE origin.
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ankorstore/go-modules/authentication"
	"github.com/ankorstore/go-modules/authentication/authenticationtest"
	"github.com/ankorstore/go-modules/httpserversecurity/httpserversecuritytest"
)

// Func applies the authentication described by a case's raw auth fixture to req.
// It matches the e2e.Runner.Auth field shape.
type Func func(tb testing.TB, req *http.Request, raw json.RawMessage)

// Ankorstore returns the default auth Func: it routes on the "type" field of a
// case's auth fixture to the matching Ankorstore principal adapter. e2e.Runner
// uses it automatically when its Auth field is nil.
func Ankorstore() Func {
	return Dispatch(map[string]Func{
		"guest":         Guest,
		"machine":       Machine,
		"admin":         Admin,
		"brand":         Brand,
		"retailer":      Retailer,
		"impersonation": Impersonation,
	})
}

// Dispatch returns an auth function that routes on the "type" field of the raw
// auth fixture to one of the provided adapters. An absent or "none" type is a
// no-op (anonymous request); an unknown type fails the test.
//
// An optional "origin" field ("internal" or "external") is applied
// orthogonally to the principal — including for anonymous requests — so a case
// can exercise routes that gate on the request's SPIFFE origin:
//
//	"auth": {"type": "brand", "accountUuid": "...", "origin": "internal"}
func Dispatch(adapters map[string]Func) Func {
	return func(tb testing.TB, req *http.Request, raw json.RawMessage) {
		tb.Helper()

		if len(raw) == 0 {
			return
		}

		var head struct {
			Type   string `json:"type"`
			Origin string `json:"origin"`
		}
		if err := json.Unmarshal(raw, &head); err != nil {
			tb.Fatalf("auth: decode type from %s: %v", raw, err)
		}

		applyOrigin(tb, req, head.Origin)

		switch head.Type {
		case "", "none":
			return
		default:
			adapter, ok := adapters[head.Type]
			if !ok {
				tb.Fatalf("auth: unsupported type %q", head.Type)
			}
			adapter(tb, req, raw)
		}
	}
}

// applyOrigin marks the request's SPIFFE origin when the fixture sets one.
func applyOrigin(tb testing.TB, req *http.Request, origin string) {
	tb.Helper()

	switch origin {
	case "":
		return
	case "internal":
		prepare(tb, req, httpserversecuritytest.WithInternalRequestOrigin())
	case "external":
		prepare(tb, req, httpserversecuritytest.WithExternalRequestOrigin())
	default:
		tb.Fatalf("auth: unknown origin %q (want %q or %q)", origin, "internal", "external")
	}
}

// Guest sends the request as a guest entity.
//
//	"auth": {"type": "guest", "clientId": "..."}
func Guest(tb testing.TB, req *http.Request, raw json.RawMessage) {
	tb.Helper()

	var d struct {
		ClientID string `json:"clientId"`
	}
	decode(tb, raw, &d)

	prepare(tb, req, httpserversecuritytest.WithGuestRequestEntity(authenticationtest.GuestTokenData{
		ClientID: d.ClientID,
	}))
}

// Machine sends the request as a machine entity.
//
//	"auth": {"type": "machine", "clientId": "..."}
func Machine(tb testing.TB, req *http.Request, raw json.RawMessage) {
	tb.Helper()

	var d struct {
		ClientID string `json:"clientId"`
	}
	decode(tb, raw, &d)

	prepare(tb, req, httpserversecuritytest.WithMachineRequestEntity(authenticationtest.MachineTokenData{
		ClientID: d.ClientID,
	}))
}

// Admin sends the request as an admin entity with the given roles/permissions.
//
//	"auth": {"type": "admin", "roles": ["..."], "permissions": ["..."]}
func Admin(tb testing.TB, req *http.Request, raw json.RawMessage) {
	tb.Helper()

	var d struct {
		ClientID    string   `json:"clientId"`
		EntityUUID  string   `json:"entityUuid"`
		Roles       []string `json:"roles"`
		Permissions []string `json:"permissions"`
	}
	decode(tb, raw, &d)

	prepare(tb, req, httpserversecuritytest.WithAdminRequestEntity(authenticationtest.AdminTokenData{
		ClientID:         d.ClientID,
		EntityUUID:       d.EntityUUID,
		AdminRoles:       d.Roles,
		AdminPermissions: d.Permissions,
	}))
}

// Brand sends the request as a brand account.
//
//	"auth": {"type": "brand", "accountUuid": "..."}
func Brand(tb testing.TB, req *http.Request, raw json.RawMessage) {
	tb.Helper()

	d := decodeAccount(tb, raw)

	prepare(tb, req, httpserversecuritytest.WithBrandRequestAccount(authenticationtest.BrandTokenData{
		ClientID:     d.ClientID,
		EntityUUID:   d.EntityUUID,
		AccountUUID:  d.AccountUUID,
		AccountEmail: d.AccountEmail,
	}))
}

// Retailer sends the request as a retailer account.
//
//	"auth": {"type": "retailer", "accountUuid": "..."}
func Retailer(tb testing.TB, req *http.Request, raw json.RawMessage) {
	tb.Helper()

	d := decodeAccount(tb, raw)

	prepare(tb, req, httpserversecuritytest.WithRetailerRequestAccount(authenticationtest.RetailerTokenData{
		ClientID:     d.ClientID,
		EntityUUID:   d.EntityUUID,
		AccountUUID:  d.AccountUUID,
		AccountEmail: d.AccountEmail,
	}))
}

// Impersonation sends the request as an admin impersonating a brand or retailer
// account. accountType selects the impersonated account kind.
//
//	"auth": {"type": "impersonation", "accountType": "brand", "accountUuid": "..."}
func Impersonation(tb testing.TB, req *http.Request, raw json.RawMessage) {
	tb.Helper()

	var d struct {
		ClientID                string   `json:"clientId"`
		EntityUUID              string   `json:"entityUuid"`
		AccountType             string   `json:"accountType"`
		AccountUUID             string   `json:"accountUuid"`
		AccountEmail            string   `json:"accountEmail"`
		ImpersonatorUUID        string   `json:"impersonatorUuid"`
		ImpersonatorRoles       []string `json:"impersonatorRoles"`
		ImpersonatorPermissions []string `json:"impersonatorPermissions"`
	}
	decode(tb, raw, &d)

	accountType, err := parseAccountType(d.AccountType)
	if err != nil {
		tb.Fatalf("auth: impersonation fixture %s: %v", raw, err)
	}

	prepare(tb, req, httpserversecuritytest.WithImpersonationRequest(authenticationtest.ImpersonationTokenData{
		ClientID:                     d.ClientID,
		EntityUUID:                   d.EntityUUID,
		AccountType:                  accountType,
		AccountUUID:                  d.AccountUUID,
		AccountEmail:                 d.AccountEmail,
		ImpersonatorUUID:             d.ImpersonatorUUID,
		ImpersonatorAdminRoles:       d.ImpersonatorRoles,
		ImpersonatorAdminPermissions: d.ImpersonatorPermissions,
	}))
}

// accountFixture is the shared field set of the brand and retailer adapters.
type accountFixture struct {
	ClientID     string `json:"clientId"`
	EntityUUID   string `json:"entityUuid"`
	AccountUUID  string `json:"accountUuid"`
	AccountEmail string `json:"accountEmail"`
}

func decodeAccount(tb testing.TB, raw json.RawMessage) accountFixture {
	tb.Helper()

	var d accountFixture
	decode(tb, raw, &d)

	return d
}

// parseAccountType maps the JSON accountType to the go-modules enum; only brand
// and retailer accounts can be impersonated.
func parseAccountType(s string) (authentication.AccountType, error) {
	switch s {
	case "brand":
		return authentication.BrandAccount, nil
	case "retailer":
		return authentication.RetailerAccount, nil
	default:
		var zero authentication.AccountType

		return zero, fmt.Errorf("accountType must be %q or %q, got %q", "brand", "retailer", s)
	}
}

// decode unmarshals the raw auth fixture into dst, failing the test on error.
func decode(tb testing.TB, raw json.RawMessage, dst any) {
	tb.Helper()

	if err := json.Unmarshal(raw, dst); err != nil {
		tb.Fatalf("auth: decode fixture from %s: %v", raw, err)
	}
}

// prepare applies a single request-security option to req, failing on error.
func prepare(tb testing.TB, req *http.Request, fn httpserversecuritytest.RequestSecurityFn) {
	tb.Helper()

	if err := httpserversecuritytest.PrepareRequestSecurity(req, fn); err != nil {
		tb.Fatalf("auth: prepare request security: %v", err)
	}
}
