package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhil-datla/Presence/internal/config"
	"github.com/akhil-datla/Presence/internal/server"
	"github.com/akhil-datla/Presence/internal/store"
	"github.com/pterm/pterm"
)

//go:embed web/static
var staticFiles embed.FS

var version = "dev"

func main() {
	// CLI flags (env vars take precedence via config.Load)
	port := flag.Int("port", 0, "server port (default: 8080, or PORT env)")
	dbPath := flag.String("db", "", "database path (default: presence.db, or DATABASE_PATH env)")
	flag.Parse()

	cfg := config.Load()

	// CLI flags override defaults only if explicitly set
	if *port != 0 {
		cfg.Port = *port
	}
	if *dbPath != "" {
		cfg.DatabasePath = *dbPath
	}

	banner()

	// Database
	db, err := store.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal("failed to open database: ", err)
	}
	defer func() { _ = db.Close() }()

	// Server
	srv := server.New(cfg, db, staticFiles)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		fmt.Println("\nShutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal("shutdown error: ", err)
		}
	}()

	pterm.Info.Printfln("Server running on http://localhost:%d", cfg.Port)
	pterm.Info.Printfln("API base: http://localhost:%d/api/v1", cfg.Port)

	if err := srv.Start(); err != nil {
		// ErrServerClosed is expected after graceful shutdown
		if err.Error() != "http: Server closed" {
			log.Fatal(err)
		}
	}
}

func banner() {
	pterm.DefaultCenter.Print(
		pterm.DefaultHeader.
			WithFullWidth().
			WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).
			WithMargin(10).
			Sprint("Presence"),
	)
	pterm.Info.Printfln("Attendance tracking API  v%s", version)
	fmt.Println()
}
