package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestLookupTenantByHost(t *testing.T) {
	q := mockQuerier{
		qr: func(ctx context.Context, sql string, args ...any) pgx.Row {
			host := args[0].(string)
			if host == "ok.local" {
				return mockRow{scan: func(dest ...any) error {
					*(dest[0].(*int64)) = 42
					return nil
				}}
			}
			return mockRow{scan: func(dest ...any) error { return errors.New("no rows") }}
		},
		ex: func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
			return pgconn.NewCommandTag("OK"), nil
		},
	}
	id, err := LookupTenantByHost(context.Background(), q, "ok.local")
	require.NoError(t, err)
	require.Equal(t, int64(42), id)

	_, err = LookupTenantByHost(context.Background(), q, "missing")
	require.Error(t, err)
}
