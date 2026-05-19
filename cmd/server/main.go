package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"baokaobao/internal/config"
	"baokaobao/internal/migrations"
	"baokaobao/internal/model"
	"baokaobao/internal/pkg/utils"
	"baokaobao/internal/router"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	configPath      = flag.String("config", "config/config.yaml", "config file path")
	createAdminUser = flag.String("create-admin", "", "create admin user: username")
	createAdminPass = flag.String("create-admin-pass", "", "admin password (min 8 chars, must include upper, lower, digit, and special char)")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.InitConfig(*configPath); err != nil {
		return fmt.Errorf("init config failed: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	db, err := initDB()
	if err != nil {
		return fmt.Errorf("init db failed: %w", err)
	}

	// Save *sql.DB reference for graceful shutdown (do NOT call db.DB() in shutdown path)
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get underlying sql.DB failed: %w", err)
	}

	if os.Getenv("BAOKAOBAO_AUTO_MIGRATE") == "true" {
		zap.S().Info("Running AutoMigrate...")
		if err := migrations.AutoMigrate(db); err != nil {
			return fmt.Errorf("auto migrate failed: %w", err)
		}
	} else {
		zap.S().Info("AutoMigrate skipped (set BAOKAOBAO_AUTO_MIGRATE=true to enable)")
	}

	if *createAdminUser != "" {
		if *createAdminPass == "" {
			fmt.Fprintln(os.Stderr, "admin password is required (use -create-admin-pass flag)")
			os.Exit(1)
		}
		if err := utils.ValidateAdminPassword(*createAdminPass); err != nil {
			fmt.Fprintf(os.Stderr, "password too weak: %v\n", err)
			os.Exit(1)
		}
		return createAdmin(db, *createAdminUser, *createAdminPass)
	}

	r := router.SetupRouterWithDB(db)

	addr := fmt.Sprintf("%s:%d", config.GlobalConfig.App.Host, config.GlobalConfig.App.Port)
	zap.S().Infof("Server starting on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdown goroutine
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		zap.S().Info("Shutting down gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			zap.S().Errorf("Server forced to shutdown: %v", err)
		}

		sqlDB.Close()
		zap.S().Sync()
	}()

	zap.S().Infof("Server listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func initDB() (*gorm.DB, error) {
	gormLogger := logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(mysql.Open(config.GlobalConfig.Database.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB failed: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.GlobalConfig.Database.ConnMaxLifetime) * time.Second)

	return db, nil
}

func createAdmin(db *gorm.DB, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}

	admin := &model.AdminUser{
		Username:     username,
		PasswordHash: string(hash),
		Nickname:     username,
		Role:         "admin",
		Status:       1,
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("create admin failed: %w", err)
	}

	fmt.Printf("Admin user '%s' created successfully!\n", username)
	return nil
}
