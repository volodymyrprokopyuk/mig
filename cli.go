package mig

import (
	"context"
	"errors"
	"regexp"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/cli/v3"
)

var sqlFiles = []string{"*.apply.sql", "*.revert.sql"}

var reApplySchema = regexp.MustCompile(`^(?:\w{2,20}|all)$`)
var reApplyVersion = regexp.MustCompile(`^(?:\d{5}|latest)$`)

func applyAction(ctx context.Context, cmd *cli.Command) error {
  schema := cmd.String("schema")
  if !reApplySchema.MatchString(schema) {
    return errors.New("invalid schema format")
  }
  version := cmd.String("version")
  if !reApplyVersion.MatchString(version) {
    return errors.New("invalid version format")
  }
  dry := cmd.Bool("dry")
  pgw, err := pgxpool.New(ctx, urlPostgres)
  if err != nil {
    return err
  }
  defer pgw.Close()
  tpl, err := template.ParseFS(fs, sqlFiles...)
  if err != nil {
    return err
  }
  if schema == "all" {
    for _, schema = range schemas {
      err := Apply(ctx, tpl, pgw, schema, version, dry)
      if err != nil {
        return err
      }
    }
    return nil
  }
  return Apply(ctx, tpl, pgw, schema, version, dry)
}

func ApplyCmd() *cli.Command {
  cmd := &cli.Command{
    Name: "apply",
    Usage: "Apply not applied migrations to all or specific schema",
    Action: applyAction,
  }
  cmd.Flags = []cli.Flag{
    &cli.StringFlag{
      Name: "schema", Usage: "schema to migrate or all", Required: true,
    },
    &cli.StringFlag{
      Name: "version", Usage: "version to apply or latest", Required: true,
    },
    &cli.BoolFlag{
      Name: "dry", Usage: "show apply plan, but do not apply",
    },
  }
  return cmd
}

var reRevertVersion = regexp.MustCompile(`^\d{5}$`)

func revertAction(ctx context.Context, cmd *cli.Command) error {
  schema := cmd.String("schema")
  if !reApplySchema.MatchString(schema) {
    return errors.New("invalid schema format")
  }
  version := cmd.String("version")
  if !reRevertVersion.MatchString(version) {
    return errors.New("invalid version format")
  }
  dry := cmd.Bool("dry")
  pgw, err := pgxpool.New(ctx, urlPostgres)
  if err != nil {
    return err
  }
  defer pgw.Close()
  tpl, err := template.ParseFS(fs, sqlFiles...)
  if err != nil {
    return err
  }
  if schema == "all" {
    for i := len(schemas) - 1; i >= 0; i-- {
      schema = schemas[i]
      err := Revert(ctx, tpl, pgw, schema, version, dry)
      if err != nil {
        return err
      }
    }
    return nil
  }
  return Revert(ctx, tpl, pgw, schema, version, dry)
}

func RevertCmd() *cli.Command {
  cmd := &cli.Command{
    Name: "revert",
    Usage: "Revert already applied migrations from all or specific schema",
    Action: revertAction,
  }
  cmd.Flags = []cli.Flag{
    &cli.StringFlag{
      Name: "schema", Usage: "schema to revert or all", Required: true,
    },
    &cli.StringFlag{
      Name: "version", Usage: "version to revert", Required: true,
    },
    &cli.BoolFlag{
      Name: "dry", Usage: "show revert plan, but do not revert",
    },
  }
  return cmd
}
