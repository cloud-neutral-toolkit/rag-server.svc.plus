package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rag-server"
	"rag-server/api"
	"rag-server/config"
	"rag-server/internal/auth"
	"rag-server/internal/cache"
	rconfig "rag-server/internal/rag/config"
	"rag-server/proxy"
)

var (
	configPath string
	logLevel   string
)

var rootCmd = &cobra.Command{
	Use:   "xcontrol-server",
	Short: "Start the xcontrol server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load(configPath)
		if err != nil {
			slog.Warn("load config", "err", err)
			cfg = &config.Config{}
		}
		if logLevel != "" {
			cfg.Log.Level = logLevel
		}
		if configPath != "" {
			api.ConfigPath = configPath
			rconfig.ServerConfigPath = configPath
		}
		proxy.Set(cfg.Global.Proxy)

		level := slog.LevelInfo
		switch strings.ToLower(cfg.Log.Level) {
		case "debug":
			level = slog.LevelDebug
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
		slog.SetDefault(logger)

		// Environment variable overrides for Cloud Run / Container support
		if pgURL := os.Getenv("DATABASE_URL"); pgURL != "" {
			cfg.Global.VectorDB.PGURL = pgURL
		} else if pgURL := os.Getenv("PG_URL"); pgURL != "" {
			cfg.Global.VectorDB.PGURL = pgURL
		}

		api.ConfigureServiceDB(nil)
		dsn := cfg.Global.VectorDB.DSN()
		var (
			conn  *pgx.Conn
			sqlDB *sql.DB
		)
		if dsn != "" {
			logger.Debug("connecting to postgres", "dsn", dsn)
			connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)
			conn, err = pgx.Connect(connectCtx, dsn)
			connectCancel()
			if err != nil {
				logger.Error("postgres connect error", "err", err)
			} else {
				logger.Info("postgres connected")
			}

			gormDB, gormErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if gormErr != nil {
				logger.Error("gorm postgres connect error", "err", gormErr)
			} else {
				api.ConfigureServiceDB(gormDB)
				sqlDB, err = gormDB.DB()
				if err != nil {
					logger.Error("postgres db handle error", "err", err)
				}
			}
		} else {
			logger.Warn("postgres dsn not provided")
		}

		if sqlDB != nil {
			defer func() {
				if cerr := sqlDB.Close(); cerr != nil {
					logger.Error("close postgres db", "err", cerr)
				}
			}()
		}

		var tokenCache *auth.TokenCache
		if conn != nil {
			cacheStore := cache.NewPostgresStore(conn, cache.Options{
				Table:      cfg.Global.Cache.Table,
				DefaultTTL: cfg.Global.Cache.DefaultTTL.Duration,
			})
			if err := cacheStore.EnsureSchema(context.Background()); err != nil {
				logger.Error("cache schema init failed", "err", err)
			} else {
				tokenCache = auth.NewTokenCache(cacheStore)
				cacheTable := cfg.Global.Cache.Table
				if cacheTable == "" {
					cacheTable = "cache_kv"
				}
				logger.Info("postgres cache ready", "table", cacheTable)
			}
		} else {
			logger.Warn("postgres cache disabled; no database connection")
		}

		r := server.New(
			api.RegisterRoutes(conn, cfg.Sync.Repo.Proxy),
		)

		// 启用认证中间件
		if cfg.Auth.Enable {
			logger.Info("enabling authentication middleware")

			// 创建认证客户端
			authConfig := auth.DefaultConfig()
			authConfig.AuthURL = cfg.Auth.AuthURL
			authConfig.PublicToken = cfg.Auth.PublicToken

			authClient := auth.NewAuthClient(authConfig)

			// 创建中间件配置
			middlewareConfig := auth.DefaultMiddlewareConfig(authClient)
			if cfg.Global.Cache.DefaultTTL.Duration > 0 {
				middlewareConfig.CacheTTL = cfg.Global.Cache.DefaultTTL.Duration
			}
			middlewareConfig.TokenCache = tokenCache

			// 添加健康检查跳过路径
			middlewareConfig.SkipPaths = append(middlewareConfig.SkipPaths, "/health", "/healthz", "/ping")

			// 应用中间件（全局）
			r.Use(auth.VerifyTokenMiddleware(middlewareConfig))

			// 添加健康检查路由
			r.GET("/health", auth.HealthCheckHandler(authClient))
			r.GET("/healthz", auth.HealthCheckHandler(authClient))
			r.GET("/ping", auth.HealthCheckHandler(authClient))

			logger.Info("authentication middleware enabled",
				"auth_url", cfg.Auth.AuthURL,
				"cache_ttl", middlewareConfig.CacheTTL.String(),
			)
		} else {
			logger.Warn("authentication is disabled")
			healthHandler := func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "ok",
					"auth":   "disabled",
				})
			}
			r.GET("/health", healthHandler)
			r.GET("/healthz", healthHandler)
		}

		server.UseCORS(r, logger, cfg.Server)

		addr := cfg.Server.Addr
		if port := os.Getenv("PORT"); port != "" {
			addr = ":" + port
		} else if addr == "" {
			addr = ":8080"
		}

		srv := &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout.Duration,
			WriteTimeout: cfg.Server.WriteTimeout.Duration,
		}

		logger.Info("starting http server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server shutdown", "err", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&configPath, "config", "", "path to server configuration file")
	rootCmd.Flags().StringVar(&logLevel, "log-level", "", "log level (debug, info, warn, error)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
