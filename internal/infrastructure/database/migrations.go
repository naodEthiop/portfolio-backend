package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type sqlMigration struct {
	Version  int64
	Filename string
	Path     string
}

func ApplySQLMigrations(db *gorm.DB, migrationsDir string) error {
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	migrations, err := discoverSQLMigrations(migrationsDir)
	if err != nil {
		return err
	}

	applied, err := appliedMigrationVersions(db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}

		sqlBytes, err := os.ReadFile(m.Path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", m.Filename, err)
		}
		sqlText := strings.TrimSpace(string(sqlBytes))
		if sqlText == "" {
			continue
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(`LOCK TABLE app_migrations IN EXCLUSIVE MODE`).Error; err != nil {
				return err
			}

			var count int64
			if err := tx.Raw(`SELECT COUNT(1) FROM app_migrations WHERE version = ?`, m.Version).Scan(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return nil
			}

			if err := tx.Exec(sqlText).Error; err != nil {
				return fmt.Errorf("exec %s: %w", m.Filename, err)
			}

			if err := tx.Exec(`INSERT INTO app_migrations (version, filename) VALUES (?, ?)`, m.Version, m.Filename).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.Filename, err)
		}

		applied[m.Version] = true
	}

	return nil
}

func ensureMigrationsTable(db *gorm.DB) error {
	return db.Exec(`
CREATE TABLE IF NOT EXISTS app_migrations (
  version BIGINT PRIMARY KEY,
  filename TEXT NOT NULL,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`).Error
}

func appliedMigrationVersions(db *gorm.DB) (map[int64]bool, error) {
	type row struct {
		Version int64
	}
	var rows []row
	if err := db.Raw(`SELECT version FROM app_migrations`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[int64]bool, len(rows))
	for _, r := range rows {
		out[r.Version] = true
	}
	return out, nil
}

func discoverSQLMigrations(dir string) ([]sqlMigration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir %s: %w", dir, err)
	}

	var out []sqlMigration
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		versionPart := strings.SplitN(name, "_", 2)[0]
		version, err := strconv.ParseInt(versionPart, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid migration filename %s (bad version prefix)", name)
		}

		out = append(out, sqlMigration{
			Version:  version,
			Filename: name,
			Path:     filepath.Join(dir, name),
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out, nil
}

