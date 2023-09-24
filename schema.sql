
CREATE TABLE IF NOT EXISTS tracks(
    id SERIAL,
    start_at TIMESTAMP NOT NULL,
    finish_at TIMESTAMP DEFAULT NULL,
    "type" INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    company_type INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);
