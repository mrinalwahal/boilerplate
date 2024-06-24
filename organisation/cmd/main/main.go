package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mrinalwahal/boilerplate/organisation/cmd/main/router"
	"github.com/mrinalwahal/boilerplate/organisation/config"
	"github.com/mrinalwahal/boilerplate/organisation/db/organisation"
	"github.com/mrinalwahal/boilerplate/organisation/service"
	"github.com/mrinalwahal/boilerplate/pkg/middleware"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
)

func main() {

	// Load the configuration.
	config := config.Get()

	//	Setup the logger.
	addSource := false
	level := config.Logs.Level.ToSlogLevel()
	if level == slog.LevelDebug {
		addSource = true
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	}))
	logger = logger.
		With("service", "organisation").
		With("environment", config.Env)

	//	Setup the gorm logger.
	handler := logger.With("layer", "database").Handler()
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(handler),                        // since v1.3.0
		slogGorm.WithTraceAll(),                              // trace all messages
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, level), // set log level (default: slog.LevelInfo)
	)

	// Open a database connection.
	conn, err := gorm.Open(config.Database.Dialector(), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		panic(err)
	}

	// Configure connection pooling.
	//
	// Link: https://gorm.io/docs/generic_interface.html#Connection-Pool
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	// GORM provides Prometheus plugin to collect DBStats or user-defined metrics
	// https://gorm.io/docs/prometheus.html
	// https://github.com/go-gorm/prometheus
	//
	// db.Use(prometheus.New(prometheus.Config{
	// 	DBName:          "db1",                       // use `DBName` as metrics label
	// 	RefreshInterval: 15,                          // Refresh metrics interval (default 15 seconds)
	// 	PushAddr:        "prometheus pusher address", // push metrics if `PushAddr` configured
	// 	StartServer:     true,                        // start http server to expose metrics
	// 	HTTPServerPort:  8080,                        // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
	// 	MetricsCollector: []prometheus.MetricsCollector{
	// 		&prometheus.MySQL{
	// 			VariableNames: []string{"Threads_running"},
	// 		},
	// 	}, // user defined metrics
	// }))

	// Get the database layer.
	db := organisation.NewDB(&organisation.DBConfig{
		Conn: conn,
	})

	// Get the service layer.
	service := service.NewService(&service.Config{
		DB:     db,
		Logger: logger,
	})

	//	Initialize the router.
	router := router.NewHTTPRouter(&router.HTTPRouterConfig{
		Service: service,
		Logger:  logger,
	})

	// Prepare the middleware chain.
	// The order of the middlewares is important.
	// Recommended order: Request ID -> RateLimit -> CORS -> Logging -> Recover -> Auth -> Cache -> Compression
	middlewareLogger := logger.With("protocol", "HTTP/1.0")
	chain := middleware.Chain(
		middleware.RequestID,
		middleware.TraceID,
		middleware.CorrelationID,
		// TODO: middleware.RateLimit,
		middleware.CORS(nil),
		middleware.Recover(&middleware.RecoverConfig{
			Logger: middlewareLogger,
		}),
		middleware.Logging(&middleware.LoggingConfig{
			Logger: middlewareLogger,
		}),
		middleware.JWT(&middleware.JWTConfig{
			Key: config.Authentication.Key.Key,
			ExceptionalRoutes: []string{
				"/login",
				"/records/healthz",
			},
		}),
	)

	// Prepare the base router.
	// baseRouter := http.NewServeMux()
	// baseRouter.Handle("/organisations/", http.StripPrefix("/organisations", router))

	//	Configure and start the server.
	server := http.Server{
		Addr:     fmt.Sprintf(":%s", *config.Port),
		Handler:  chain(router),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	fmt.Println("Server is running on port 8080")
	server.ListenAndServe()

	// Close the database connection.
	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}
