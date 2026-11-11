package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/galiullindo/todolist-rest-api/internal/config"
	taskHandler "github.com/galiullindo/todolist-rest-api/internal/handler"
	"github.com/galiullindo/todolist-rest-api/internal/model"
	taskRepository "github.com/galiullindo/todolist-rest-api/internal/repository"
	taskService "github.com/galiullindo/todolist-rest-api/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func testServer(db *sql.DB) *httptest.Server {
	taskRepo := taskRepository.New(db)
	taskSvc := taskService.New(taskRepo)
	taskHndl := taskHandler.New(taskSvc)

	mux := http.NewServeMux()
	taskHndl.RegisterRoutes(mux)

	return httptest.NewServer(mux)
}

func testClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func trimTestFlags(args []string) []string {
	for i, a := range args {
		if a == "--" {
			return args[i+1:]
		}
	}
	return []string{}
}

func createTask(t *testing.T, url string, task model.Task) {
	t.Helper()
	client := testClient()

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&task); err != nil {
		t.Errorf("encode response body: %v\n", err)
	}

	resp, err := client.Post(url, "application/json", &body)
	if err != nil {
		t.Errorf("client post: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusCreated)
	}
}

func TestAPI(t *testing.T) {
	args := trimTestFlags(os.Args[1:])

	var configsDir string
	fs := flag.NewFlagSet("e2e_test", flag.ExitOnError)
	fs.StringVar(&configsDir, "configs-dir", ".", "path to configs directory")
	fs.Parse(args)

	cfg, err := config.Load(configsDir, config.EnvTest)
	if err != nil {
		t.Fatalf("load config: %v\n", err)
	}

	db, err := sql.Open("sqlite3", cfg.Database.URI)
	if err != nil {
		t.Fatalf("open db: %v\n", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v\n", err)
	}

	// SAVEPOINT для отката базы к исходному состоянию после тестов
	if _, err := db.Exec("SAVEPOINT e2e_test"); err != nil {
		t.Fatalf("set savepoint: %v\n", err)
	}
	rollback := func() {
		if _, err := db.Exec("ROLLBACK TO e2e_test"); err != nil {
			t.Errorf("rollback to savepoint: %v\n", err)
		}
	}
	release := func() {
		if _, err := db.Exec("RELEASE SAVEPOINT e2e_test"); err != nil {
			t.Errorf("release savepoint: %v\n", err)
		}
	}
	defer release()

	server := testServer(db)
	defer server.Close()

	client := testClient()

	t.Run("GET /tasks возвращает пустой список", func(t *testing.T) {
		want := []model.Task{}

		resp, err := client.Get(server.URL + "/tasks")
		if err != nil {
			t.Errorf("client get: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusOK)
		}

		var tasks []model.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Errorf("decode response body: %v\n", err)
		}

		if !slices.Equal(tasks, want) {
			t.Errorf("got %v, want %v\n", tasks, want)
		}
	})

	t.Run("POST /tasks создает задачу", func(t *testing.T) {
		defer rollback()
		want := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}

		var body bytes.Buffer
		if err := json.NewEncoder(&body).Encode(&want); err != nil {
			t.Errorf("encode response body: %v\n", err)
		}

		resp, err := client.Post(server.URL+"/tasks", "application/json", &body)
		if err != nil {
			t.Errorf("client post: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusCreated)
		}

		var task model.Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			t.Errorf("decode response body: %v\n", err)
		}

		if task != want {
			t.Errorf("got %v, want %v\n", task, want)
		}
	})

	t.Run("POST /tasks возвращает 400 для невалидного JSON", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{invalid json}`))

		resp, err := client.Post(server.URL+"/tasks", "application/json", body)
		if err != nil {
			t.Errorf("client post: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("GET /tasks возвращает список задач", func(t *testing.T) {
		defer rollback()
		want := []model.Task{
			{ID: 1, Title: "Новая задача 1", Status: "Статус задачи 1"},
			{ID: 2, Title: "Новая задача 2", Status: "Статус задачи 2"},
			{ID: 3, Title: "Новая задача 3", Status: "Статус задачи 3"},
		}

		for _, task := range want {
			createTask(t, server.URL+"/tasks", task)
		}

		resp, err := client.Get(server.URL + "/tasks")
		if err != nil {
			t.Errorf("client get: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusOK)
		}

		var tasks []model.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Errorf("decode response body: %v\n", err)
		}

		if !slices.Equal(tasks, want) {
			t.Errorf("got %v, want %v\n", tasks, want)
		}
	})

	t.Run("GET /tasks/{id} возвращает задачу по id", func(t *testing.T) {
		defer rollback()
		want := model.Task{ID: 1, Title: "Новая задача", Status: "Статус задачи"}

		createTask(t, server.URL+"/tasks", want)

		resp, err := client.Get(server.URL + "/tasks/1")
		if err != nil {
			t.Errorf("client get: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusOK)
		}

		var task model.Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			t.Errorf("decode response body: %v\n", err)
		}

		if task != want {
			t.Errorf("got %v, want %v\n", task, want)
		}
	})

	t.Run("GET /tasks/{id} возвращает 400 для невалидного id", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/tasks/abc")
		if err != nil {
			t.Errorf("client get: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusBadRequest)
		}
	})

	t.Run("GET /tasks/{id} возвращает 404 для несуществующего id", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/tasks/1")
		if err != nil {
			t.Errorf("client get: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("response status code got %d, want %d\n", resp.StatusCode, http.StatusNotFound)
		}
	})
}
