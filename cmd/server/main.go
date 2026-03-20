package main

import (
	"flag"
	"fmt"
	"os"

	"baokaobao/internal/config"
	"baokaobao/internal/handler"
	"baokaobao/internal/middleware"
	"baokaobao/internal/repository"
	"baokaobao/internal/router"
	"baokaobao/internal/service"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	configPath = flag.String("config", "config/config.yaml", "config file path")
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

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zap.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zap.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zap.DefaultLineEnding,
			EncodeLevel:    zap.LowercaseLevelEncoder,
			EncodeTime:     zap.ISO8601TimeEncoder,
			EncodeDuration: zap.SecondsDurationEncoder,
			EncodeCaller:   zap.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	db, err := initDB()
	if err != nil {
		return fmt.Errorf("init db failed: %w", err)
	}

	middleware.InitJWT()

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	r := router.SetupRouter(h)

	addr := fmt.Sprintf("%s:%d", config.GlobalConfig.App.Host, config.GlobalConfig.App.Port)
	zap.S().Infof("Server starting on %s", addr)

	return r.Run(addr)
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
	sqlDB.SetConnMaxLifetime(0)

	return db, nil
}
