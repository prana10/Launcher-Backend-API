package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"launcherbackend_api/internal/domain/entity"
	repo "launcherbackend_api/internal/domain/repository"
)

type PostgresOTARepository struct {
	db *sql.DB
}

func NewPostgresOTARepository(db *sql.DB) repo.OTARepository {
	return &PostgresOTARepository{
		db: db,
	}
}

func (r *PostgresOTARepository) Create(ctx context.Context, ota entity.OTA) (entity.OTA, error) {
	query := `
		INSERT INTO otas (id, app_id, version_name, version_code, release_notes, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, app_id, version_name, version_code, release_notes, url, created_at, updated_at
	`

	if ota.ID == "" {
		ota.ID = uuid.NewString()
	}

	now := time.Now()
	ota.CreatedAt = now
	ota.UpdatedAt = now

	err := r.db.QueryRowContext(
		ctx,
		query,
		ota.ID,
		ota.AppID,
		ota.VersionName,
		ota.VersionCode,
		ota.ReleaseNotes,
		ota.URL,
		ota.CreatedAt,
		ota.UpdatedAt,
	).Scan(
		&ota.ID,
		&ota.AppID,
		&ota.VersionName,
		&ota.VersionCode,
		&ota.ReleaseNotes,
		&ota.URL,
		&ota.CreatedAt,
		&ota.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return entity.OTA{}, fmt.Errorf("ota with this version already exists: %w", err)
			}
		}
		return entity.OTA{}, fmt.Errorf("failed to create ota: %w", err)
	}

	return ota, nil
}

func (r *PostgresOTARepository) Get(ctx context.Context, id string, appID string) (entity.OTA, string, error) {
	// Check if both id and appID are provided
	if id != "" && appID != "" {
		return nil, "", fmt.Errorf("cannot provide both id and appID, choose one")
	}

	// Check if neither id nor appID are provided
	if id == "" && appID == "" {
		return nil, "", fmt.Errorf("must provide either id or appID")
	}

	// If id is provided, get a single OTA
	if id != "" {
		query := `
			SELECT id, app_id, version_name, version_code, release_notes, url, created_at, updated_at
			FROM otas
			WHERE id = $1
		`

		var ota entity.OTA
		err := r.db.QueryRowContext(ctx, query, id).Scan(
			&ota.ID,
			&ota.AppID,
			&ota.VersionName,
			&ota.VersionCode,
			&ota.ReleaseNotes,
			&ota.URL,
			&ota.CreatedAt,
			&ota.UpdatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, "", fmt.Errorf("ota not found: %w", err)
			}
			return nil, "", fmt.Errorf("failed to get ota: %w", err)
		}

		return ota, "", nil
	}

	// If appID is provided, get multiple OTAs
	query := `
		SELECT id, app_id, version_name, version_code, release_notes, url, created_at, updated_at
		FROM otas
		WHERE app_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, appID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get otas by app id: %w", err)
	}
	defer rows.Close()

	var ota entity.OTA
	for rows.Next() {
		if err := rows.Scan(
			&ota.ID,
			&ota.AppID,
			&ota.VersionName,
			&ota.VersionCode,
			&ota.ReleaseNotes,
			&ota.URL,
			&ota.CreatedAt,
			&ota.UpdatedAt,
		); err != nil {
			return nil, "", fmt.Errorf("failed to scan ota row: %w", err)
		}
	}

	return ota, "", nil
}

func (r *PostgresOTARepository) GetAll(ctx context.Context, cursor string, limit int) ([]entity.OTA, string, int64, error) {
	query := `
		SELECT id, app_id, version_name, version_code, release_notes, url, created_at, updated_at
		FROM otas
	`

	var total int64
	countErr := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM otas").Scan(&total)
	if countErr != nil {
		return nil, "", 0, fmt.Errorf("failed to count otas: %w", countErr)
	}

	params := []interface{}{}
	if cursor != "" {
		query += " WHERE id > $1"
		params = append(params, cursor)
	}

	query += " ORDER BY id ASC LIMIT $" + fmt.Sprintf("%d", len(params)+1)
	params = append(params, limit+1)

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get all otas: %w", err)
	}
	defer rows.Close()

	var otas []entity.OTA
	for rows.Next() {
		var ota entity.OTA
		if err := rows.Scan(
			&ota.ID,
			&ota.AppID,
			&ota.VersionName,
			&ota.VersionCode,
			&ota.ReleaseNotes,
			&ota.URL,
			&ota.CreatedAt,
			&ota.UpdatedAt,
		); err != nil {
			return nil, "", 0, fmt.Errorf("failed to scan ota row: %w", err)
		}
		otas = append(otas, ota)
	}

	var nextCursor string
	if len(otas) > limit {
		nextCursor = otas[limit-1].ID
		otas = otas[:limit]
	}

	return otas, nextCursor, total, nil
}

func (r *PostgresOTARepository) Update(ctx context.Context, ota entity.OTA) (entity.OTA, error) {
	query := `
		UPDATE otas
		SET app_id = $2, version_name = $3, version_code = $4, release_notes = $5, url = $6, updated_at = $7
		WHERE id = $1
		RETURNING id, app_id, version_name, version_code, release_notes, url, created_at, updated_at
	`

	ota.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(
		ctx,
		query,
		ota.ID,
		ota.AppID,
		ota.VersionName,
		ota.VersionCode,
		ota.ReleaseNotes,
		ota.URL,
		ota.UpdatedAt,
	).Scan(
		&ota.ID,
		&ota.AppID,
		&ota.VersionName,
		&ota.VersionCode,
		&ota.ReleaseNotes,
		&ota.URL,
		&ota.CreatedAt,
		&ota.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.OTA{}, fmt.Errorf("ota not found: %w", err)
		}
		return entity.OTA{}, fmt.Errorf("failed to update ota: %w", err)
	}

	return ota, nil
}

func (r *PostgresOTARepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM otas WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ota: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ota not found")
	}

	return nil
} 