package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"go-rethinkdb/pkg/config"
	"go-rethinkdb/pkg/logger"
	pkg_rethinkdb "go-rethinkdb/pkg/rethinkdb"
	todo_http "go-rethinkdb/todo/delivery/http"
	todo_repository "go-rethinkdb/todo/repository"
	todod_service "go-rethinkdb/todo/service"
	"go-rethinkdb/utils"
	response "go-rethinkdb/utils/response"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger, // Log API request calls
		// middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	return router
}

// PrintAllRoutes - printing all routes
func PrintAllRoutes(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Error(err)
	}
}

func main() {
	utils.InitializeValidator()

	// Load environment variables
	err := config.LoadConfig()
	if err != nil {
		logger.Error(err)
	}

	// Init RethinkDB
	session, err := pkg_rethinkdb.InitRethinkDB()
	if err != nil {
		logger.Error(err)
	}

	router := Routes()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, response.H{
			"success": "true",
			"code":    200,
			"message": "Services run properly",
		})
	})

	// Repository
	todoRepo := todo_repository.NewRethinkDBTodoRepository(session)

	// Service
	todoService := todod_service.NewTodoService(todoRepo)

	// Handler
	todoHandler := todo_http.NewTodoHTTPHandler(router, todoService)
	todoHandler.RegisterRoutes()

	// Print
	PrintAllRoutes(router)

	addr := fmt.Sprintf("%s%s", ":", os.Getenv("PORT"))
	err = http.ListenAndServe(addr, router)
	if err != nil {
		logger.Error(err)
	}
}
