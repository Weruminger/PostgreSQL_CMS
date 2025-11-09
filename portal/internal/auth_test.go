package internal

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func jwtWithPayload(payload string) string {
	h := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	p := base64.RawURLEncoding.EncodeToString([]byte(payload))
	return h + "." + p + "."
}

func TestParseJWTPayload_OK(t *testing.T) {
	j := jwtWithPayload(`{"realm_access":{"roles":["editor","other"]}}`)
	m, err := parseJWTPayload(j)
	require.NoError(t, err)
	ra := m["realm_access"].(map[string]any)
	roles := ra["roles"].([]any)
	require.Contains(t, roles, "editor")
}

func TestParseJWTPayload_Bad(t *testing.T) {
	_, err := parseJWTPayload("no-dot")
	require.Error(t, err)
	_, err = parseJWTPayload(strings.Repeat(".", 3))
	require.Error(t, err)
}
