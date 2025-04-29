package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	
	"launcherbackend_api/internal/config"
	"launcherbackend_api/internal/domain/entity"
)

// Sample OTA data for seeding
var otaSamples = []entity.OTA{
	{
		AppID:        "com.yapindo.launcher",
		VersionName:  "1.0.0",
		VersionCode:  100,
		ReleaseNotes: "Initial release of the Yapindo Launcher app",
		URL:          "https://storage.example.com/yapindo/launcher/v1.0.0/app.apk",
	},
	{
		AppID:        "com.yapindo.launcher",
		VersionName:  "1.1.0",
		VersionCode:  110,
		ReleaseNotes: "Bug fixes and performance improvements",
		URL:          "https://storage.example.com/yapindo/launcher/v1.1.0/app.apk",
	},
	{
		AppID:        "com.yapindo.launcher",
		VersionName:  "1.2.0",
		VersionCode:  120,
		ReleaseNotes: "Added new features:\n- Dark mode\n- Push notifications\n- Improved UI",
		URL:          "https://storage.example.com/yapindo/launcher/v1.2.0/app.apk",
	},
	{
		AppID:        "com.yapindo.launcher.pro",
		VersionName:  "1.0.0",
		VersionCode:  100,
		ReleaseNotes: "Initial release of the Yapindo Launcher Pro app",
		URL:          "https://storage.example.com/yapindo/launcher-pro/v1.0.0/app.apk",
	},
	{
		AppID:        "com.yapindo.launcher.pro",
		VersionName:  "1.0.1",
		VersionCode:  101,
		ReleaseNotes: "Hotfix for authentication issues",
		URL:          "https://storage.example.com/yapindo/launcher-pro/v1.0.1/app.apk",
	},
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database, starting seeding...")

	// Insert sample OTA data
	for _, ota := range otaSamples {
		insertOTA(db, ota)
	}

	fmt.Println("Seeding completed!")
}

func insertOTA(db *sql.DB, ota entity.OTA) {
	query := `
		INSERT INTO otas (id, app_id, version_name, version_code, release_notes, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (app_id, version_code) DO UPDATE SET
			version_name = EXCLUDED.version_name,
			release_notes = EXCLUDED.release_notes,
			url = EXCLUDED.url,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	if ota.ID == "" {
		ota.ID = uuid.NewString()
	}

	now := time.Now()
	ota.CreatedAt = now
	ota.UpdatedAt = now

	var id string
	err := db.QueryRowContext(
		context.Background(),
		query,
		ota.ID,
		ota.AppID,
		ota.VersionName,
		ota.VersionCode,
		ota.ReleaseNotes,
		ota.URL,
		ota.CreatedAt,
		ota.UpdatedAt,
	).Scan(&id)

	if err != nil {
		log.Printf("Failed to insert OTA record for %s (v%s): %v", ota.AppID, ota.VersionName, err)
		return
	}

	fmt.Printf("Inserted OTA for %s (v%s) with ID: %s\n", ota.AppID, ota.VersionName, id)
} 