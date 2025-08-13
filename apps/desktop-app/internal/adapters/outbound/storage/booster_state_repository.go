package storage

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    "github.com/adrg/xdg"
    "github.com/oLenador/mulltbost/internal/core/domain/entities"
    _ "github.com/mattn/go-sqlite3"
)

type BoosterStateRepository struct {
    db *sql.DB
}

func NewBoosterStateRepository(appName string) (*BoosterStateRepository, error) {
    // Pega o caminho padrão do XDG_DATA_HOME
    dbDir := filepath.Join(xdg.DataHome, appName)
    if err := os.MkdirAll(dbDir, 0o755); err != nil {
        return nil, fmt.Errorf("erro criando diretório do banco: %w", err)
    }

    dbPath := filepath.Join(dbDir, "states.db")

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    repo := &BoosterStateRepository{db: db}

    createTable := `
    CREATE TABLE IF NOT EXISTS Booster_states (
        id TEXT PRIMARY KEY,
        data TEXT NOT NULL
    );`
    if _, err := db.Exec(createTable); err != nil {
        return nil, err
    }

    return repo, nil
}

func (r *BoosterStateRepository) Save(ctx context.Context, state *entities.BoosterState) error {
    jsonData, err := json.Marshal(state)
    if err != nil {
        return err
    }

    _, err = r.db.ExecContext(ctx,
        `INSERT INTO Booster_states (id, data) VALUES (?, ?)
         ON CONFLICT(id) DO UPDATE SET data=excluded.data`,
        state.ID, string(jsonData),
    )
    return err
}

func (r *BoosterStateRepository) GetByID(ctx context.Context, id string) (*entities.BoosterState, error) {
    row := r.db.QueryRowContext(ctx, `SELECT data FROM Booster_states WHERE id = ?`, id)

    var jsonData string
    if err := row.Scan(&jsonData); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    var state entities.BoosterState
    if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
        return nil, err
    }

    return &state, nil
}

func (r *BoosterStateRepository) GetAll(ctx context.Context) ([]*entities.BoosterState, error) {
    rows, err := r.db.QueryContext(ctx, `SELECT data FROM Booster_states`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var states []*entities.BoosterState
    for rows.Next() {
        var jsonData string
        if err := rows.Scan(&jsonData); err != nil {
            return nil, err
        }
        var state entities.BoosterState
        if err := json.Unmarshal([]byte(jsonData), &state); err != nil {
            return nil, err
        }
        states = append(states, &state)
    }

    return states, nil
}

func (r *BoosterStateRepository) Delete(ctx context.Context, id string) error {
    _, err := r.db.ExecContext(ctx, `DELETE FROM Booster_states WHERE id = ?`, id)
    return err
}
