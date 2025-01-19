CREATE TABLE
    IF NOT EXISTS users.roles (id bigserial PRIMARY KEY, title text NOT NULL);

INSERT INTO
    users.roles (title)
VALUES
    ('ADMIN'),
    ('CREATOR'),
    ('RESPONDER');

CREATE TABLE
    IF NOT EXISTS users.users (
        id bigserial PRIMARY KEY,
        role_id integer NOT NULL REFERENCES users.roles (id),
        name text NOT NULL,
        email citext UNIQUE NOT NULL,
        password_hash bytea NOT NULL,
        is_activated bool NOT NULL,
        created_at timestamptz (0) NOT NULL DEFAULT NOW (),
        updated_at timestamptz (0) NOT NULL DEFAULT NOW (),
        version integer NOT NULL DEFAULT 1
    );