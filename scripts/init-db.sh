#!/bin/bash

# Script untuk membuat dan menginisialisasi database setelah docker-compose up

echo "Menunggu database PostgreSQL siap..."
until docker exec launcher_db pg_isready -U postgres > /dev/null 2>&1; do
  echo "Menunggu PostgreSQL..."
  sleep 2
done

echo "PostgreSQL siap! Membuat tables..."

# Membuat file SQL sementara
cat > init.sql << EOF
-- Database initialization script for example Launcher Backend API

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
EOF

# Jalankan script SQL di container PostgreSQL
docker exec -i launcher_db psql -U postgres -d launcher_db < init.sql

# Hapus file sementara
rm init.sql

echo "Database berhasil diinisialisasi!" 