package database

import (
	"errors"
	"event-collector/internal/firstapp/config"
	"event-collector/internal/firstapp/logger"
	"event-collector/pkg/parse"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"log"
)

func RunMigration(cfg *config.Config, logger logger.Logger) error {
	m, err := migrate.New("file://"+cfg.GetEnv(cfg.Database.Migrate.PATH), cfg.GetEnv(cfg.Database.URI))
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	migrateType := cfg.GetEnv(cfg.Database.Migrate.Type)
	switch migrateType {
	case "up":
		return up(m, logger)
	case "down":
		return down(m, logger)
	case "step":
		return step(m, logger, cfg.GetEnv(cfg.Database.Migrate.Step))
	default:
		logger.Error("unsupported migration type", zap.Error(errors.New("unsupported migration type"+migrateType)))
		return errors.New("unsupported migration type" + migrateType)
	}

}

func up(m *migrate.Migrate, logger logger.Logger) error {
	err := m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Database migration succeeded: nothing to change ✅")
			return nil
		}
		logger.Error("Migrate up failed", zap.Error(err))
		return err
	}

	logger.Info("Database migration succeeded ✅")
	return nil
}

func down(m *migrate.Migrate, logger logger.Logger) error {
	err := m.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Database migration rolled back succeeded: nothing to change ✅")
			return nil
		}
		logger.Error("Migrate down failed", zap.Error(err))
		return err
	}

	logger.Info("Database migration rolled back successfully ✅")
	return nil
}

func step(m *migrate.Migrate, logger logger.Logger, stepsStr string) error {
	steps, _ := parse.ToPrimary[int](stepsStr)

	err := m.Steps(steps)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("Database migration step succeeded: No migrations to apply/rollback ✅")
			return nil
		}
		logger.Error("Migration step failed", zap.Int("steps", steps), zap.Error(err))
		return err
	}

	direction := "applied"
	if steps < 0 {
		direction = "rolled back"
	}
	logger.Info("Migration successful ✅", zap.Int("steps", steps), zap.String("direction", direction))
	return nil
}
