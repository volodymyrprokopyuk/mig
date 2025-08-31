package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/volodymyrprokopyuk/mig"
	"github.com/volodymyrprokopyuk/mig/postgres"
)

func setupMigration() {
  type schemaData struct {
    Schema string
  }
  mig.SetFS(&postgres.FS)
  mig.SetURL(os.Getenv("POSTGRES_URL"))
  mig.SetSchema("util", "util", nil)
  mig.SetSchema("eu", "region", schemaData{Schema: "eu"})
  mig.SetSchema("us", "region", schemaData{Schema: "us"})
}

func migCmd() *cli.Command {
  setupMigration()
  cmd := &cli.Command{
    Name: "mig",
    Usage: "Apply and revert migrations to PostgreSQL",
    Version: os.Getenv("MIG_VERSION"),
    UseShortOptionHandling: true,
    Commands: []*cli.Command{mig.ApplyCmd(), mig.RevertCmd()},
  }
  return cmd
}

func main() {
  err := migCmd().Run(context.Background(), os.Args)
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err)
    os.Exit(1)
  }
}
