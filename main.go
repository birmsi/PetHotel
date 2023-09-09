package main

import (
	"PetHotel/handlers"
	"PetHotel/repositories"
	"PetHotel/services"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type DBInfo struct {
	address  string
	user     string
	password string
	name     string
}

type ApplicationHandlers struct {
	BoxHandlers handlers.BoxHandler
}

type ApplicationInfo struct {
	applicationPort string
	DBInfo          DBInfo
	handlers        ApplicationHandlers
}

func main() {
	slogHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	slogger := slog.New(slogHandler)

	applicationInfo, err := handleEnvVariables()
	if err != nil {
		slogger.Error(err.Error())
		os.Exit(1)
	}

	db, err := openDB(applicationInfo.DBInfo)
	if err != nil {
		slogger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	var version string
	if err = db.QueryRow("select version()").Scan(&version); err != nil {
		slogger.Error(err.Error())
		os.Exit(1)
	}

	slogger.Debug(fmt.Sprintf("version - %s", version))

	applicationInfo.handlers = initApplicationHandlers(db, slogger)

	r := chi.NewRouter()

	r.Get("/", home)
	r.Route("/box", func(r chi.Router) {
		r.Get("/", applicationInfo.handlers.BoxHandlers.GetBoxesView)

		r.Get("/create", applicationInfo.handlers.BoxHandlers.CreateBoxView)
		r.Post("/create", applicationInfo.handlers.BoxHandlers.CreateBoxPost)

		r.Get("/{boxID}/update", applicationInfo.handlers.BoxHandlers.GetBoxUpdateView)
		r.Post("/{boxID}/update", applicationInfo.handlers.BoxHandlers.GetBoxUpdatePut)

		r.Delete("/{boxID}", applicationInfo.handlers.BoxHandlers.BoxDelete)
	})

	server := http.Server{
		Addr:    ":4000",
		Handler: r,
	}

	slogger.Debug(fmt.Sprintf("Ready 2 go! Listening on :%s", applicationInfo.applicationPort))

	if err := server.ListenAndServe(); err != nil {
		slogger.Error(fmt.Sprintf("Failed to init server - %s", err.Error()))
		os.Exit(1)
	}
}

func openDB(dbInfo DBInfo) (*sql.DB, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", dbInfo.user, dbInfo.password, dbInfo.address, dbInfo.name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func handleEnvVariables() (ApplicationInfo, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbAddress := os.Getenv("DB_ADDRESS")
	if dbAddress == "" {
		log.Fatal("Missing DB_ADDRESS from .env")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("Missing DB_NAME from .env")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("Missing DB_USER from .env")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Missing DB_PASSWORD from .env")
	}

	applicationPort := os.Getenv("APPLICATION_PORT")
	if dbPassword == "" {
		log.Fatal("Missing APPLICATION_PORT from .env")
	}

	return ApplicationInfo{
		applicationPort: applicationPort,
		DBInfo: DBInfo{
			address:  dbAddress,
			user:     dbUser,
			password: dbPassword,
			name:     dbName,
		},
	}, err
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page :D"))
}

func initApplicationHandlers(db *sql.DB, slogger *slog.Logger) ApplicationHandlers {
	var applicationHandlers ApplicationHandlers

	boxRepository := repositories.NewBoxRepository(db, slogger)
	boxService := services.NewService(boxRepository, slogger)
	applicationHandlers.BoxHandlers = handlers.NewBoxHandler(boxService, slogger)

	return applicationHandlers
}
