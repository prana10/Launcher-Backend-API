# Script PowerShell untuk membuat dan menginisialisasi database setelah docker-compose up

Write-Host "Menunggu database PostgreSQL siap..." -ForegroundColor Yellow
do {
    $isReady = $false
    try {
        $result = docker exec launcher_db pg_isready -U postgres 2>&1
        if ($result -match "accepting connections") {
            $isReady = $true
        }
    } catch {
        Write-Host "Menunggu PostgreSQL..." -ForegroundColor Cyan
    }
    if (-not $isReady) {
        Start-Sleep -Seconds 2
    }
} until ($isReady)

Write-Host "PostgreSQL siap! Membuat tables..." -ForegroundColor Green

# Membuat file SQL sementara
@"
-- Database initialization script for Example Launcher Backend API

-- OTA (Over-the-Air) updates table
CREATE TABLE IF NOT EXISTS otas (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(255) NOT NULL,
    version_name VARCHAR(100) NOT NULL,
    version_code INTEGER NOT NULL,
    release_notes TEXT,
    url VARCHAR(2048) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(app_id, version_code)
);

-- Index untuk pencarian berdasarkan app_id
CREATE INDEX IF NOT EXISTS idx_otas_app_id ON otas(app_id);

-- Fungsi untuk mengupdate kolom updated_at
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger untuk mengupdate kolom updated_at
CREATE TRIGGER update_otas_updated_at
BEFORE UPDATE ON otas
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Tambahkan data sample (opsional)
INSERT INTO otas (app_id, version_name, version_code, release_notes, url)
VALUES 
('com.example.launcher', '1.0.0', 100, 'Initial release', 'https://storage.example.com/apps/launcher-1.0.0.apk'),
('com.example.launcher', '1.0.1', 101, 'Bug fixes', 'https://storage.example.com/apps/launcher-1.0.1.apk')
ON CONFLICT (app_id, version_code) DO NOTHING;
"@ | Out-File -FilePath .\init.sql -Encoding utf8

# Jalankan script SQL di container PostgreSQL
Get-Content .\init.sql | docker exec -i launcher_db psql -U postgres -d launcher_db

# Hapus file sementara
Remove-Item .\init.sql

Write-Host "Database berhasil diinisialisasi!" -ForegroundColor Green 