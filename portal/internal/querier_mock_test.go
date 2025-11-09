package internal

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type rowFunc func(dest ...any) error
type mockRow struct{ scan rowFunc }

func (m mockRow) Scan(dest ...any) error { return m.scan(dest...) }

type mockRows struct {
	i    int
	data [][]any
	err  error
}

func (m *mockRows) Next() bool { m.i++; return m.err == nil && m.i <= len(m.data) }
func (m *mockRows) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	for i := range dest {
		switch d := dest[i].(type) {
		case *int64:
			*d = m.data[m.i-1][i].(int64)
		case *string:
			*d = m.data[m.i-1][i].(string)
		default:
			return errors.New("unsupported scan type")
		}
	}
	return nil
}
func (m *mockRows) Err() error { return m.err }
func (m *mockRows) Close()     {}

type mockQuerier struct {
	qr func(ctx context.Context, sql string, args ...any) pgx.Row
	q  func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	ex func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func (m mockQuerier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.qr(ctx, sql, args...)
}
func (m mockQuerier) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.q == nil {
		return nil, errors.New("not implemented")
	}
	return m.q(ctx, sql, args...)
}
func (m mockQuerier) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if m.ex == nil {
		return pgconn.CommandTag(""), errors.New("not implemented")
	}
	return m.ex(ctx, sql, args...)
}
