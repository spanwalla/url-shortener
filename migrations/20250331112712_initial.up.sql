CREATE TABLE links(
    alias VARCHAR(10) PRIMARY KEY,
    uri TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_links_uri_hash ON links USING HASH (uri);