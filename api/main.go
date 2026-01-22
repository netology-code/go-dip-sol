package main

import (
	"advanced-blog-management-system/internal/handler"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/repository"
	"advanced-blog-management-system/internal/service"
	"advanced-blog-management-system/pkg/auth"
	"advanced-blog-management-system/pkg/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	projectRoot, err := findProjectRoot()
	if err == nil && projectRoot != "" {
		log.Printf("Project root detected: %s", projectRoot)
		if err := os.Chdir(projectRoot); err != nil {
			log.Printf("Warning: Failed to change working directory: %v", err)
		}
	} else {
		log.Printf("Using current working directory (Docker environment)")
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found in project root")
	}

	cfg := loadConfig()

	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	log.Println("Running database migrations...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	userRepo := repository.NewUserRepo(db)
	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)

	userService := service.NewUserService(userRepo, jwtManager)
	postService := service.NewPostService(postRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	eventLogger := logger.NewEventLogger("logs.txt")
	eventLogger.Start()

	authHandler := handler.NewAuthHandler(userService)
	postHandler := handler.NewPostHandler(postService, eventLogger)
	commentHandler := handler.NewCommentHandler(commentService, eventLogger)

	loggingMiddleware := middleware.NewLoggingMiddleware(log.New(os.Stdout, "", log.LstdFlags))
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router := chi.NewRouter()

	router.Use(loggingMiddleware.Recovery)
	router.Use(loggingMiddleware.Logger)
	router.Use(loggingMiddleware.CORS)

	router.Post("/api/register", authHandler.Register)
	router.Post("/api/login", authHandler.Login)
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	apiRouter := chi.NewRouter()

	apiRouter.Group(func(r chi.Router) {
		r.Get("/posts", postHandler.GetAll)
		r.Get("/posts/{id}", postHandler.GetByID)
		r.Get("/posts/{postId}/comments", commentHandler.GetByPost)
	})

	apiRouter.Group(func(r chi.Router) {
		r.Use(middleware.ToMiddleware(authMiddleware.RequireAuth))
		r.Post("/posts", postHandler.Create)
		r.Post("/posts/{postId}/comments", commentHandler.Create)
	})

	router.Mount("/api", apiRouter)

	server := &http.Server{
		Addr:    cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	eventLogger.Stop()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(currentDir, "migrations")); err == nil {
		return currentDir, nil
	}

	oneLevelUp := filepath.Join(currentDir, "..")
	if _, err := os.Stat(filepath.Join(oneLevelUp, "migrations")); err == nil {
		return filepath.Clean(oneLevelUp), nil
	}

	twoLevelsUp := filepath.Join(currentDir, "..", "..")
	if _, err := os.Stat(filepath.Join(twoLevelsUp, "migrations")); err == nil {
		return filepath.Clean(twoLevelsUp), nil
	}

	return "", err
}

type Config struct {
	ServerHost     string
	ServerPort     int
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpiryHours int
}

func loadConfig() *Config {
	return &Config{
		ServerHost:     getEnv("SERVER_HOST", "localhost"),
		ServerPort:     getEnvAsInt("SERVER_PORT", 8080),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "blogdb"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key"),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}
