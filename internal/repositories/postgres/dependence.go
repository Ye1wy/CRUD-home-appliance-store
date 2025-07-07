package postgres

import "github.com/jackc/pgx/v5/pgconn"

type PgxCommanTag struct {
	tag pgconn.CommandTag
}

func (ct *PgxCommanTag) RowsAffected() int64 {
	return ct.RowsAffected()
}
