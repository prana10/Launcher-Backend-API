CREATE TABLE IF NOT EXISTS otas (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(255) NOT NULL,
    version_name VARCHAR(255) NOT NULL,
    version_code INTEGER NOT NULL,
    release_notes TEXT,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX idx_otas_app_id ON otas(app_id);
CREATE UNIQUE INDEX idx_otas_app_id_version_code ON otas(app_id, version_code); 