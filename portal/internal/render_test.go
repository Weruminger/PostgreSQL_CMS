package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestGetTheme_Cache(t *testing.T) {
	calls := 0
	q := mockQuerier{
		qr: func(ctx context.Context, sql string, args ...any) pgx.Row {
			calls++
			if calls == 1 {
				return mockRow{scan: func(dest ...any) error {
					*(dest[0].(*string)) = `{{ define "layout" }}ok{{ end }}`
					*(dest[1].(*string)) = `{"header":"<h1>H</h1>"}`
					return nil
				}}
			}
			return mockRow{scan: func(dest ...any) error { return errors.New("db down") }}
		},
		ex: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag("OK"), nil
		},
	}
	themeCache = map[int64]cacheEntry{} // reset
	th, err := getTheme(context.Background(), q, 1)
	require.NoError(t, err)
	require.Contains(t, th.Layout, "layout")
	require.Equal(t, "<h1>H</h1>", th.Partials["header"])

	th2, err := getTheme(context.Background(), q, 1)
	require.NoError(t, err)
	require.Equal(t, th.Layout, th2.Layout)
}
