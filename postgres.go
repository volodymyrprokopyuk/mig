package mig

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func qryLastApplied(
  ctx context.Context, pgp *pgxpool.Pool, schema string, count int,
) ([]migration, error) {
  sql := `select mig.* from util.migration_last_applied($1, $2) mig;`
  rows, _ := pgp.Query(ctx, sql, schema, count)
  migrations, err := pgx.CollectRows(rows, pgx.RowToStructByName[migration])
  if err != nil {
    if !strings.Contains(err.Error(), `schema "util" does not exist`) {
      return nil, err
    }
  }
  return migrations, nil
}

func qryApplyMigration(
  ctx context.Context, tpl *template.Template, pgp *pgxpool.Pool,
  schema, name string,
) error {
  var sql bytes.Buffer
  data := schemaData[schema]
  err := tpl.ExecuteTemplate(&sql, name, data)
  if err != nil {
    return err
  }
  opt := pgx.TxOptions{
    IsoLevel: pgx.Serializable,
    AccessMode: pgx.ReadWrite,
    DeferrableMode: pgx.NotDeferrable,
  }
  return pgx.BeginTxFunc(ctx, pgp, opt, func(tx pgx.Tx) error {
    _, err := tx.Exec(ctx, sql.String())
    return err
  })
}
