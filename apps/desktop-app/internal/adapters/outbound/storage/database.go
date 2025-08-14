package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/oLenador/mulltbost/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
    // Diretório onde o banco vai ficar
    dbDir := filepath.Join(xdg.DataHome, config.UserDir)

    if err := os.MkdirAll(dbDir, 0o755); err != nil {
        return nil, fmt.Errorf("criar diretório do banco: %w", err)
    }

    // Caminho completo do arquivo
    dbPath := filepath.Join(dbDir, config.UserDbName)

    gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("abrir sqlite com gorm: %w", err)
    }

    return gormDB, nil
}


func AutoMigrateModels(db *gorm.DB, models ...interface{}) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	return db.AutoMigrate(models...)
}