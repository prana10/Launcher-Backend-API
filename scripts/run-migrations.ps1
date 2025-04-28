# Script PowerShell untuk menjalankan migrasi database tanpa Docker

# Variabel konfigurasi database
$DbHost = "localhost"
$DbPort = "5432"
$DbUser = "postgres"
$DbPass = "postgres"
$DbName = "launcher_db"
$MigrationsPath = "./migrations"

Write-Host "Catatan: Pastikan database '$DbName' sudah dibuat sebelum menjalankan migrasi" -ForegroundColor Yellow
Write-Host "Jika menggunakan Docker, database biasanya sudah dibuat otomatis oleh docker-compose" -ForegroundColor Yellow
Write-Host "Jika menggunakan PostgreSQL lokal, buat database terlebih dahulu dengan pgAdmin atau psql" -ForegroundColor Yellow

# Periksa apakah golang-migrate sudah terinstal
$migrateExists = $null
try {
    $migrateExists = Get-Command migrate -ErrorAction SilentlyContinue
} catch {
    # Tidak ada masalah jika command tidak ditemukan
}

# Jika golang-migrate belum terinstal, install dengan go
if ($null -eq $migrateExists) {
    Write-Host "Golang-migrate belum terinstal. Menginstal..." -ForegroundColor Yellow
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    
    # Tambahkan ke PATH jika belum ada
    $goPath = go env GOPATH
    $migratePath = Join-Path $goPath "bin"
    
    if ($env:Path -notlike "*$migratePath*") {
        Write-Host "Menambahkan golang-migrate ke PATH untuk sesi ini..." -ForegroundColor Yellow
        $env:Path += ";$migratePath"
    }
}

# Buat connection string - perbaikan format
$connectionString = "postgresql://$DbUser`:$DbPass@$DbHost`:$DbPort/$DbName`?sslmode=disable"

Write-Host "Menggunakan connection string: $connectionString" -ForegroundColor Cyan

# Jalankan migrasi
Write-Host "Menjalankan migrasi database..." -ForegroundColor Green
try {
    migrate -path $MigrationsPath -database $connectionString up
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Migrasi berhasil dijalankan!" -ForegroundColor Green
    } else {
        Write-Host "Migrasi gagal dengan kode: $LASTEXITCODE" -ForegroundColor Red
        Write-Host "Periksa apakah PostgreSQL berjalan dan database '$DbName' sudah dibuat." -ForegroundColor Red
        
        # Menampilkan informasi tambahan untuk debug
        Write-Host "`nInformasi tambahan untuk memecahkan masalah:" -ForegroundColor Cyan
        Write-Host "1. Pastikan PostgreSQL sudah berjalan di $DbHost pada port $DbPort" -ForegroundColor Cyan
        Write-Host "2. Pastikan user '$DbUser' dengan password yang sesuai memiliki akses" -ForegroundColor Cyan
        Write-Host "3. Pastikan database '$DbName' sudah dibuat di PostgreSQL" -ForegroundColor Cyan
        Write-Host "4. Jika menggunakan Docker, pastikan container PostgreSQL berjalan dengan: docker ps" -ForegroundColor Cyan
        Write-Host "5. Jika menggunakan Docker, coba gunakan: docker-compose exec postgres psql -U postgres -c 'CREATE DATABASE launcher_db;'" -ForegroundColor Cyan
    }
} catch {
    Write-Host "Error menjalankan migrasi: $_" -ForegroundColor Red
} 