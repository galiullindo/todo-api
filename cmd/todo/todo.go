package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	taskHandler "github.com/galiullindo/todo-api/internal/api"
	"github.com/galiullindo/todo-api/internal/config"
	taskRepository "github.com/galiullindo/todo-api/internal/repository"
	taskService "github.com/galiullindo/todo-api/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		root = filepath.Dir(os.Args[0])
		cfg  = config.New()
	)

	log.Println("open db")
	db, err := sql.Open("sqlite3", filepath.Join(root, cfg.DBPath))
	if err != nil {
		log.Fatalf("open db: %v\n", err)
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
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
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
