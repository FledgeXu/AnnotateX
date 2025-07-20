-- +goose Up
-- +goose StatementBegin
SELECT
  'up SQL query';

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(100) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  display_name VARCHAR(100),
  email VARCHAR(255),
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW (),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW ()
);

CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  description TEXT
);

CREATE TABLE user_roles (
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  role_id BIGINT NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, role_id)
);

INSERT INTO
  roles (name, description)
VALUES
  (
    'super_admin',
    'Super administrator with all permissions'
  ),
  (
    'admin',
    'Platform administrator with user and task management privileges'
  ),
  (
    'project_manager',
    'Responsible for project task distribution and monitoring'
  ),
  ('reviewer', 'Reviewer for submitted annotations'),
  (
    'labeler',
    'Labeler responsible for data annotation tasks'
  ),
  (
    'unassigned',
    'User has not been assigned a role yet'
  );

CREATE TABLE organizations (
  id BIGSERIAL PRIMARY KEY, -- Unique organization ID
  type VARCHAR(50) NOT NULL DEFAULT 'vendor', -- Organization type: vendor, internal, etc.
  name VARCHAR(255) NOT NULL UNIQUE, -- Organization name (must be unique)
  code VARCHAR(64) UNIQUE, -- Optional short code or slug
  description TEXT, -- Optional description
  is_active BOOLEAN DEFAULT TRUE, -- Whether the organization is active
  created_at TIMESTAMP
  WITH
    TIME ZONE DEFAULT NOW (), -- Creation timestamp
    updated_at TIMESTAMP
  WITH
    TIME ZONE DEFAULT NOW () -- Last update timestamp
);

CREATE TABLE projects (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  modality TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT now (),
  updated_at TIMESTAMP DEFAULT now ()
);

CREATE TABLE datasets (
  id BIGSERIAL PRIMARY KEY,
  project_id BIGINT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  format_version TEXT NOT NULL,
  status TEXT DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW ()
);

CREATE TABLE dataset_tasks (
  id BIGSERIAL PRIMARY KEY,
  dataset_id BIGINT NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW ()
);

CREATE TABLE task_frames (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT NOT NULL REFERENCES dataset_tasks (id) ON DELETE CASCADE,
  frame_index INTEGER NOT NULL,
  timestamp DOUBLE PRECISION,
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMP DEFAULT NOW (),
  UNIQUE (task_id, frame_index)
);

CREATE TABLE files (
  id BIGSERIAL PRIMARY KEY,
  content_hash TEXT NOT NULL,
  hash_method TEXT NOT NULL DEFAULT 'blake3' CHECK (hash_method IN ('md5', 'sha256', 'blake3')),
  size_bytes BIGINT,
  mime_type TEXT,
  status TEXT DEFAULT 'pending' CHECK (hash_method IN ('pending', 'ready', 'failed')),
  created_at TIMESTAMP DEFAULT NOW (),
  updated_at TIMESTAMP DEFAULT NOW (),
  UNIQUE (content_hash, hash_method) -- 唯一性约束，允许支持不同算法共存
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
  'down SQL query';

DROP TABLE IF EXISTS user_roles;

DROP TABLE IF EXISTS roles;

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
