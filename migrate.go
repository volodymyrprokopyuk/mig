package mig

import (
	"context"
	"embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var fs *embed.FS

func SetFS(efs *embed.FS) {
  fs = efs
}

var urlPostgres string

func SetURL(url string) {
  urlPostgres = url
}

var schemas = make([]string, 0, 10)
var schemaPrefix = make(map[string]string, 10)
var schemaData = make(map[string]any, 10)

func SetSchema(schema, prefix string, data any) {
  schemas = append(schemas, schema)
  schemaPrefix[schema] = prefix
  schemaData[schema] = data
}

type migration struct {
  Schema string `db:"schema"`
  Version string `db:"version"`
  Description string `db:"description"`
  Applied time.Time `db:"applied"`
}

func printMigrations(migrations []migration) {
  for i := len(migrations) - 1; i >= 0; i-- {
    mig := migrations[i]
    fmt.Printf(
      "%10s %s %-50s %s\n",
      mig.Schema, mig.Version, mig.Description,
      mig.Applied.UTC().Format("2006-01-02 15:04:05"),
    )
  }
}

var reFileVersion = regexp.MustCompile(`_(\d{5})_`)

func fsReadMigrations(
  schema, suffix string,
) (map[string]string, string, error) {
  prefix := schemaPrefix[schema]
  files, err := fs.ReadDir(".")
  if err != nil {
    return nil, "", err
  }
  migrations := make(map[string]string, len(files) / 2)
  var latest string
  for _, file := range files {
    name := file.Name()
    if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
      m := reFileVersion.FindStringSubmatch(name)
      if len(m) < 2 {
        continue
      }
      version := m[1]
      migrations[version] = name
      if version > latest {
        latest = version
      }
    }
  }
  return migrations, latest, err
}

func Apply(
  ctx context.Context, tpl *template.Template, pgp *pgxpool.Pool,
  schema, version string, dry bool,
) error {
  lastApplied, err := qryLastApplied(ctx, pgp, schema, 5)
  if err != nil {
    return err
  }
  if dry {
    printMigrations(lastApplied)
  }
  migs, latest, err := fsReadMigrations(schema, ".apply.sql")
  if err != nil {
    return err
  }
  start, end := -1, 0;
  if len(lastApplied) > 0 {
    start, err = strconv.Atoi(lastApplied[0].Version)
    if err != nil {
      return err
    }
  }
  if version == "latest" {
    if len(latest) > 0 {
      end, err = strconv.Atoi(latest)
      if err != nil {
        return err
      }
    }
  } else {
    end, err = strconv.Atoi(version)
    if err != nil {
      return err
    }
  }
  if start >= end {
    fmt.Printf(
      "version %s for schema %s is already applied\n", version, schema,
    )
    return nil
  }
  format := "=> %s applying %05d\n"
  if dry {
    format = "=> %s will apply %05d\n"
  }
  for i := start + 1; i <= end; i++ {
    fmt.Printf(format, schema, i)
    version := fmt.Sprintf("%05d", i)
    name, exist := migs[version]
    if !exist {
      return fmt.Errorf(
        "version %s for schema %s does not exist", version, schema,
      )
    }
    if dry {
      continue
    }
    err := qryApplyMigration(ctx, tpl, pgp, schema, name)
    if err != nil {
      return err
    }
  }
  return nil
}

func Revert(
  ctx context.Context, tpl *template.Template, pgp *pgxpool.Pool,
  schema, version string, dry bool,
) error {
  lastApplied, err := qryLastApplied(ctx, pgp, schema, 5)
  if err != nil {
    return err
  }
  if dry {
    printMigrations(lastApplied)
  }
  migs, _, err := fsReadMigrations(schema, ".revert.sql")
  if err != nil {
    return err
  }
  start, end := 0, -1;
  if len(lastApplied) > 0 {
    end, err = strconv.Atoi(lastApplied[0].Version)
    if err != nil {
      return err
    }
  }
  start, err = strconv.Atoi(version)
  if err != nil {
    return err
  }
  if start > end {
    return fmt.Errorf(
      "version %s for schema %s is not applied", version, schema,
    )
  }
  format := "=> %s reverting %05d\n"
  if dry {
    format = "=> %s will revert %05d\n"
  }
  for i := end; i >= start; i-- {
    fmt.Printf(format, schema, i)
    version := fmt.Sprintf("%05d", i)
    name, exist := migs[version]
    if !exist {
      return fmt.Errorf(
        "version %s for schema %s does not exist", version, schema,
      )
    }
    if dry {
      continue
    }
    err := qryApplyMigration(ctx, tpl, pgp, schema, name)
    if err != nil {
      return err
    }
  }
  return nil
}
