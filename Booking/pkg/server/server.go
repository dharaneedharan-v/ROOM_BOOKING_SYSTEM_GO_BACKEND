
package server

import (
	"net/http"
	"BookingSystem/Booking/internal/config"
	"BookingSystem/Booking/internal/handlers"
	utils "BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/pkg/database"
    "go.uber.org/zap"	
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Application struct {
	Logger    *utils.Logger
	Router    *mux.Router
	Config    *config.Config
	DBService database.DBConnector
}

func InitializeApp(env string) (*Application, error) {
	// Load the configuration
	cfg, err := config.LoadConfig("training-service", env)
	if err != nil {
		return nil, err
	}
	dbService := &database.DBService{}
	  // Logger Initialization
    var log=utils.LogConfig{Level: cfg.LogLevel,
        LogDir: cfg.LogDir,
        FileName: cfg.LogFileName,
        ServiceName: cfg.ServiceName,
    }
    // create logger instance
    logger := utils.NewLogger(log,)
	// Initialize logger
	//logger := utils.NewLogger("training-service")
	logger.Info("Logger initialized")

	// Set up the database connection
	logger.Info("Initializing database connection")
	db, err := dbService.EstablishConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}
	logger.Info("Database connection established successfully")

	// // Run auto migration
	// if err := database.AutoMigrate(db, logger); err != nil {
	// 	logger.Error("Failed to run database migrations", zap.Error(err))
	// 	return nil, err
	// }

	// // Seeding 

	// logger.Info("Seeding database...")
	// if err := database.SeedDatabase(db.Gorm); err != nil {
	// 	logger.Error("Failed to seed database",zap.Error (err))
	// } else {
	// 	logger.Info("Database seeding completed successfully")
	// }
	// Migrations
	if err := database.AutoMigrate(db, logger); err != nil {
		logger.Error("Failed to run database migrations", zap.Error(err))
		return nil, err
	}

	// Seeding 
	var userCount, roomCount, bookingCount int64

	// Check if any data exists across your three tables
	_ = db.Gorm.Table("customers").Count(&userCount)
	_ = db.Gorm.Table("rooms").Count(&roomCount)
	_ = db.Gorm.Table("bookings").Count(&bookingCount)

	// Skip seeding if any table already has entries
	if userCount > 0 || roomCount > 0 || bookingCount > 0 {
		logger.Info("Database tables already contain data. Skipping seeding process.",
			zap.Int64("customers", userCount),
			zap.Int64("rooms", roomCount),
			zap.Int64("bookings", bookingCount),
		)
	} else {
		logger.Info("Seeding database...")
		if err := database.SeedDatabase(db.Gorm); err != nil {
			logger.Error("Failed to seed database", zap.Error(err))
		} else {
			logger.Info("Database seeding completed successfully")
		}
	}

	// Create the router and set up routes
	router := mux.NewRouter()
	handlers.SetupRoutes(router, db, logger, cfg)

	return &Application{
		Logger:    logger,
		Router:    router,
		Config:    cfg,
		DBService: dbService,
	}, nil
}

// GetCorsOptions returns the CORS configuration options
func GetCorsOptions() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "accessToken", "deviceIdentifier"},
	}
}

func RunServer(app *Application) error {
	// app.Logger.Info("Starting server on port: %s", app.Config.Port) // Zap doest not support the formating like %s 
	app.Logger.Info("Starting server on port",zap.String("function","RunServer"),zap.String("port",app.Config.Port),)


	corsOpts := GetCorsOptions()
	handler := cors.New(corsOpts).Handler(app.Router)
	return http.ListenAndServe(":"+app.Config.Port, handler)
}

