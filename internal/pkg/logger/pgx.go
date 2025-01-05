package logger

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (l Logger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l.Debug("pgx query", "sql", data.SQL, "args", data.Args)
	return ctx
}

func (l Logger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		l.Error("pgx query failed", "err", data.Err)
	} else {
		l.Debug("pgx query success", "command_tag", data.CommandTag)
	}
}
