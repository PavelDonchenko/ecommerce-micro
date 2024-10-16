CREATE SCHEMA users;

CREATE TABLE IF NOT EXISTS users.users
(
    id            UUID        NOT NULL PRIMARY KEY,
    name          TEXT        NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP   NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users.claims
(
    id    SERIAL PRIMARY KEY,
    type  TEXT NOT NULL,
    value TEXT NOT NULL
);

CREATE TABLE users.users_claims (
    user_id uuid references users.users(id) ON DELETE CASCADE,
    claim_id int references users.claims(id) ON DELETE CASCADE
);

INSERT INTO users.claims(type, value)
VALUES ('admin', 'read');
INSERT INTO users.claims(type, value)
VALUES ('admin', 'create');
INSERT INTO users.claims(type, value)
VALUES ('admin', 'update');
INSERT INTO users.claims(type, value)
VALUES ('user', 'read');
INSERT INTO users.claims(type, value)
VALUES ('user', 'create');
INSERT INTO users.claims(type, value)
VALUES ('user', 'update');

