CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE images (
    id BIGSERIAL PRIMARY KEY,
    embedding VECTOR(512) NOT NULL,

    domain_url TEXT NOT NULL,
    product_url TEXT NOT NULL,
    product_name TEXT NOT NULL,
    product_seller TEXT,
    product_price NUMERIC(10, 2)

    file_name TEXT NOT NULL,
    file_size INT,
    height INT NOT NULL,
    width INT NOT NULL,
    format TEXT NOT NULL,

    created_at TIMESTAMPZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Set timestamp to update to current on update queries
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON images
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX ON images USING hsnw (embedding vector_cosine_ops);

