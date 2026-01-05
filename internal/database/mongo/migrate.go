package mongo

import (
	"errors"
	"event-collector/internal/config"
	"event-collector/internal/pkg/logger"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigration(cfg *configs.Config, logger logger.Logger) error {
	m, err := migrate.New("file://"+cfg.Database.Migrate.Path, cfg.Database.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch cfg.Database.Migrate.Type {
	case "up":
		return up(m, logger)
	case "down":
		return down(m, logger)
	case "step":
		return step(m, logger, cfg.Database.Migrate.Step)
	default:
		logger.Error("unsupported migration type", zap.Error(errors.New("unsupported migration type"+cfg.Database.Migrate.Type)))
		return errors.New("unsupported migration type" + cfg.Database.Migrate.Type)
	}

}

func up(m *migrate.Migrate, logger logger.Logger) error {
	err := m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Database migration succeeded: nothing to change ✅")
			return nil
		}
		logger.Error("Migrate up failed", zap.Error(err))
		return err
	}

	log.Println("Database migration succeeded ✅")
	return nil
}

func down(m *migrate.Migrate, logger logger.Logger) error {
	err := m.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Database migration rolled back succeeded: nothing to change ✅")
			return nil
		}
		logger.Error("Migrate down failed", zap.Error(err))
		return err
	}

	log.Println("Database migration rolled back successfully ✅")
	return nil
}

func step(m *migrate.Migrate, logger logger.Logger, steps int) error {
	err := m.Steps(steps)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Database migration step succeeded: No migrations to apply/rollback ✅")
			return nil
		}
		logger.Error("Migration step failed", zap.Int("steps", steps), zap.Error(err))
		return err
	}

	direction := "applied"
	if steps < 0 {
		direction = "rolled back"
	}

	log.Printf("Migration successful ✅, steps: %v, direction: %v", steps, direction)
	return nil
}
