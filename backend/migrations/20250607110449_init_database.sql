-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE users (
    id SERIAL PRIMARY KEY,

    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_name VARCHAR(100),
    email VARCHAR(255),

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

INSERT INTO roles (name, description) VALUES
  ('super_admin', 'Super administrator with all permissions'),
  ('admin', 'Platform administrator with user and task management privileges'),
  ('project_manager', 'Responsible for project task distribution and monitoring'),
  ('reviewer', 'Reviewer for submitted annotations'),
  ('labeler', 'Labeler responsible for data annotation tasks'),
  ('unassigned', 'User has not been assigned a role yet');

CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,                             -- Unique organization ID
    type VARCHAR(50) NOT NULL DEFAULT 'vendor',        -- Organization type: vendor, internal, etc.
    name VARCHAR(255) NOT NULL UNIQUE,                 -- Organization name (must be unique)
    code VARCHAR(64) UNIQUE,                           -- Optional short code or slug
    description TEXT,                                  -- Optional description
    is_active BOOLEAN DEFAULT TRUE,                    -- Whether the organization is active
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Creation timestamp
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()  -- Last update timestamp
);

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    modality TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
