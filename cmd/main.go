package main

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-task/internal/config"
	"test-task/internal/repository"
	"test-task/internal/transport/grpc_contract"
	"test-task/internal/transport/grpc_server"
	"test-task/internal/usecase/garantex"
	"test-task/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

const databaseURLFlagName string = "database_url"

func main() {
	app := &cli.App{
		Name:   "exchange",
		Usage:  "grpc server for getting rates from garantex",
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if c.IsSet(databaseURLFlagName) {
		cfg.DatabaseURL = c.String(databaseURLFlagName)
	}

	repo, err := initPostgreSQL(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	if err := initMigrations(cfg.DatabaseURL); err != nil {
		return err
	}

	garantexUC := garantex.New(repo)

	server := grpc_server.New(garantexUC)

	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpc_contract.RegisterAPIServer(grpcServer, server)

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`OK`))
	})

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	httpServer.Shutdown(ctx)
	grpcServer.GracefulStop()

	return nil
}

func initPostgreSQL(ctx context.Context, dbURL string) (*repository.Repository, error) {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	dbConfig.MaxConns = 5

	conn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return repository.New(conn), nil
}

func initMigrations(dbURL string) error {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}
