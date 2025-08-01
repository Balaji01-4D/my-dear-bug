CREATE TABLE IF NOT EXISTS bugs (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    language TEXT,
    snippet TEXT,
    upvotes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE upvotes (
    id SERIAL PRIMARY KEY,
    bug_id INT REFERENCES bugs(id),
    ip_hash TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
