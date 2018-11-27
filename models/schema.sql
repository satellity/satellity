CREATE TABLE IF NOT EXISTS users (
	user_id               VARCHAR(36) PRIMARY KEY,
	email                 VARCHAR(512),
	username              VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
	nickname              VARCHAR(64) NOT NULL DEFAULT '',
	encrypted_password    VARCHAR(1024),
	github_id             VARCHAR(1024) UNIQUE,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users ((LOWER(email)));
CREATE UNIQUE INDEX ON users ((LOWER(username)));
CREATE INDEX ON users (created_at);


CREATE TABLE IF NOT EXISTS sessions (
	session_id            VARCHAR(36) PRIMARY KEY,
	user_id               VARCHAR(36) NOT NULL,
	secret                VARCHAR(1024) NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON sessions (user_id);
CREATE INDEX ON sessions (created_at);


CREATE TABLE IF NOT EXISTS categories (
	category_id           VARCHAR(36) PRIMARY KEY,
	name                  VARCHAR(36) NOT NULL,
	description           VARCHAR(512) NOT NULL,
	topics_count          INTEGER NOT NULL DEFAULT 0,
	last_topic_id         VARCHAR(36),
	position              INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON categories (position);


CREATE TABLE IF NOT EXISTS topics (
	topic_id              VARCHAR(36) PRIMARY KEY,
	title                 VARCHAR(512) NOT NULL,
	body                  TEXT NOT NULL,
	comments_count        INTEGER NOT NULL DEFAULT 0,
	category_id           VARCHAR(36) NOT NULL,
	user_id               VARCHAR(36) NOT NULL,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON topics (user_id);
CREATE INDEX ON topics (category_id);
CREATE INDEX ON topics (created_at DESC);
CREATE INDEX ON topics (user_id, created_at DESC);
CREATE INDEX ON topics (score, created_at DESC);


CREATE TABLE IF NOT EXISTS comments (
  comment_id            VARCHAR(36) PRIMARY KEY,
  body                  TEXT NOT NULL,
  topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
  user_id               VARCHAR(36) NOT NULL,
  score                 INTEGER NOT NULL DEFAULT 0,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON comments (topic_id, created_at);
CREATE INDEX ON comments (user_id, created_at);
CREATE INDEX ON comments (score, created_at);
