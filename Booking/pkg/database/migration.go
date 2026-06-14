
package database

import (
	"BookingSystem/Booking/internal/models"
	"BookingSystem/Booking/internal/loggers"
	"go.uber.org/zap"
)

func AutoMigrate(db *Db, logger *loggers.Logger) error {
	logger.Info("Starting auto migration for database tables")

	err := db.Gorm.AutoMigrate(
		// &models.User{},
		&models.Customer{},
		&models.Room{},
		&models.Booking{},

	)

	if err != nil {
		logger.Error("Failed to run auto migration", zap.Error(err))
		return err
	}

	logger.Info("Auto migration completed successfully")
	return nil
}