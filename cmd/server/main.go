package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/galiullindo/todo-api/internal/config"
	taskHandler "github.com/galiullindo/todo-api/internal/handler"
	taskRepository "github.com/galiullindo/todo-api/internal/repository"
	taskService "github.com/galiullindo/todo-api/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var cliArguments struct {
		configsDir string
		env        string
	}

	fs := flag.NewFlagSet("server", flag.ExitOnError)
	fs.StringVar(&cliArguments.configsDir, "configs-dir", "", "path to configs directory")
	fs.StringVar(&cliArguments.env, "env", "default", "enviroment (prod, test, docker)")
	fs.Parse(os.Args[1:])

	log.Println("parse cli argument 'env'")
	env, err := config.ParseEnv(cliArguments.env)
	if err != nil {
		log.Printf("parse cli argument 'env': %v\n", err)
		return
	}

	log.Println("load config")
	cfg, err := config.Load(cliArguments.configsDir, env)
	if err != nil {
		log.Printf("load config: %v\n", err)
		return
	}

	log.Println("open db")
	db, err := sql.Open("sqlite3", cfg.Database.URI)
	if err != nil {
		log.Printf("open db: %v\n", err)
		return
	}
	defer db.Close()

	log.Println("ping db")
	if err := db.Ping(); err != nil {
		log.Printf("ping db: %v\n", err)
		return
	}

	taskRepo := taskRepository.New(db)
	taskSvc := taskService.New(taskRepo)
	taskHndl := taskHandler.New(taskSvc)

	mux := http.NewServeMux()
	taskHndl.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	mainCtx, mainStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer mainStop()

	log.Printf("server listen and serve on http://%s\n", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Printf("server listen and serve: %v\n", err)
		}
	}()

	<-mainCtx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer shutdownCancel()

	log.Println("server shutdown")
	if err := server.Shutdown(shutdownCtx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		log.Printf("server shutdown: %v\n", err)
	}
}
