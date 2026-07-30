package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/dakaii/vibegopher/db"
	"github.com/dakaii/vibegopher/internal/envvar"
)

func main() {
	if len(os.Args) < 2 {
		fatalUsage()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	dsn := databaseURL()
	var err error
	switch strings.ToLower(os.Args[1]) {
	case "up":
		err = db.Up(ctx, dsn)
	case "down":
		err = db.Down(ctx, dsn)
	case "status":
		err = db.Status(ctx, dsn)
	default:
		fatalUsage()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
}

func databaseURL() string {
	return envvar.PostgresDSN()
}

func fatalUsage() {
	fmt.Fprintln(os.Stderr, "usage: migrate <up|down|status>")
	fmt.Fprintln(os.Stderr, "  DATABASE_URL or POSTGRES_* must be set")
	os.Exit(2)
}
