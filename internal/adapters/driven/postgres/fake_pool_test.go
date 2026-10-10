package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// fakeRows is a minimal implementation of pgx.Rows for testing.
type fakeRows struct {
	pgx.Rows
	next      bool
	err       error
	scanErr   error
	scanVals  []any
	scanCalls int
}

func (f *fakeRows) Next() bool {
	n := f.next
	f.next = false
	return n
}

func (f *fakeRows) Close() {}

func (f *fakeRows) Err() error { return f.err }

func (f *fakeRows) Scan(dest ...any) error {
	f.scanCalls++
	if f.scanErr != nil {
		return f.scanErr
	}
	for i, val := range f.scanVals {
		if i < len(dest) {
			switch d := dest[i].(type) {
			case *string:
				if s, ok := val.(string); ok {
					*d = s
				}
			case *pgtype.Text:
				if s, ok := val.(string); ok {
					d.String = s
					d.Valid = true
				}
			case **any:
				**d = val
			}
		}
	}
	return nil
}

// fakePool is a minimal implementation of poolIface for testing.
type fakePool struct {
	queryErr   error
	execErr    error
	execNoRows bool
	rows       pgx.Rows
	row        pgx.Row
	queryCalls int
	execCalls  int
	pingErr    error
}

func (f *fakePool) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	f.queryCalls++
	if f.rows != nil {
		return f.rows, f.queryErr
	}
	if f.queryErr != nil {
		var nilRows pgx.Rows
		return nilRows, f.queryErr
	}
	return &fakeRows{}, nil
}

func (f *fakePool) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return f.row
}

func (f *fakePool) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	f.execCalls++
	if f.execErr != nil {
		return pgconn.NewCommandTag(""), f.execErr
	}
	if f.execNoRows {
		return pgconn.NewCommandTag(""), nil
	}
	return pgconn.NewCommandTag("INSERT 0 1"), nil
}

func (f *fakePool) Ping(_ context.Context) error {
	return f.pingErr
}

func (f *fakePool) Close() {}

// fakeRow is a minimal implementation of pgx.Row for testing.
type fakeRow struct {
	err error
	val any
}

func (f *fakeRow) Scan(dest ...any) error {
	if f.err != nil {
		return f.err
	}
	if len(dest) > 0 && f.val != nil {
		switch d := dest[0].(type) {
		case *string:
			if s, ok := f.val.(string); ok {
				*d = s
			}
		case *int64:
			if i, ok := f.val.(int64); ok {
				*d = i
			}
		case **any:
			**d = f.val
		}
	}
	return nil
}

// errDummy is a sentinel error for fake pool tests.
var errDummy = errors.New("dummy error")
