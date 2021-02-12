CREATE TABLE IF NOT EXISTS users (
  user_id               VARCHAR(36) PRIMARY KEY,
  email                 VARCHAR(512),
  username              VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
  nickname              VARCHAR(64) NOT NULL DEFAULT '',
  biography             VARCHAR(2048) NOT NULL DEFAULT '',
  avatar_url            VARCHAR(512) NOT NULL DEFAULT '',
  encrypted_password    VARCHAR(1024),
  github_id             VARCHAR(1024) UNIQUE,
  role                  VARCHAR(128) NOT NULL DEFAULT '',
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_emailx ON users ((LOWER(email)));
CREATE UNIQUE INDEX IF NOT EXISTS users_usernamex ON users ((LOWER(username)));
CREATE INDEX IF NOT EXISTS users_createdx ON users (created_at);


CREATE TABLE IF NOT EXISTS sessions (
  session_id            VARCHAR(36) PRIMARY KEY,
  user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  secret                VARCHAR(1024) NOT NULL,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS sessions_userx ON sessions (user_id);


CREATE TABLE IF NOT EXISTS email_verifications (
  verification_id       VARCHAR(36) PRIMARY KEY,
  email                 VARCHAR(512),
  code                  VARCHAR(512),
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS email_verifications_email_code_createdx ON email_verifications (email, code, created_at DESC);
CREATE INDEX IF NOT EXISTS email_verifications_createdx ON email_verifications (created_at DESC);


CREATE TABLE IF NOT EXISTS categories (
  category_id           VARCHAR(36) PRIMARY KEY,
  name                  VARCHAR(36) NOT NULL,
  alias                 VARCHAR(128) NOT NULL,
  description           VARCHAR(512) NOT NULL,
  topics_count          BIGINT NOT NULL DEFAULT 0,
  last_topic_id         VARCHAR(36),
  position              INTEGER NOT NULL DEFAULT 0,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS categories_positionx ON categories (position);
CREATE UNIQUE INDEX IF NOT EXISTS categories_namex ON categories (name);


CREATE TABLE IF NOT EXISTS topics (
  topic_id              VARCHAR(36) PRIMARY KEY,
  short_id              VARCHAR(256) NOT NULL UNIQUE,
  title                 VARCHAR(512) NOT NULL,
  body                  TEXT NOT NULL,
  topic_type            VARCHAR(256) NOT NULL,
  comments_count        BIGINT NOT NULL DEFAULT 0,
  bookmarks_count       BIGINT NOT NULL DEFAULT 0,
  likes_count           BIGINT NOT NULL DEFAULT 0,
  views_count           BIGINT NOT NULL DEFAULT 0,
  category_id           VARCHAR(36) NOT NULL,
  user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  score                 INTEGER NOT NULL DEFAULT 0,
  draft                 BOOL NOT NULL DEFAULT false,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS topics_draft_createdx ON topics(draft, created_at DESC);
CREATE INDEX IF NOT EXISTS topics_user_draft_createdx ON topics(user_id, draft, created_at DESC);
CREATE INDEX IF NOT EXISTS topics_category_draft_createdx ON topics(category_id, draft, created_at DESC);
CREATE INDEX IF NOT EXISTS topics_score_draft_createdx ON topics(score DESC, draft, created_at DESC);


CREATE TABLE IF NOT EXISTS topic_users (
  topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
  user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  liked                 BOOL NOT NULL DEFAULT false,
  bookmarked            BOOL NOT NULL DEFAULT false,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (topic_id, user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS topic_users_reversex ON topic_users(user_id, topic_id);
CREATE INDEX IF NOT EXISTS topic_users_likedx ON topic_users(topic_id, liked);
CREATE INDEX IF NOT EXISTS topic_users_bookmarkedx ON topic_users(topic_id, bookmarked);


CREATE TABLE IF NOT EXISTS comments (
  comment_id            VARCHAR(36) PRIMARY KEY,
  body                  TEXT NOT NULL,
  topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
  user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  score                 INTEGER NOT NULL DEFAULT 0,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS comments_topic_createdx ON comments (topic_id, created_at);
CREATE INDEX IF NOT EXISTS comments_user_createdx ON comments (user_id, created_at);
CREATE INDEX IF NOT EXISTS comments_score_createdx ON comments (score DESC, created_at);


CREATE TABLE IF NOT EXISTS products (
  product_id            VARCHAR(36) PRIMARY KEY,
  name                  VARCHAR(512) NOT NULL,
  body                  TEXT NOT NULL,
  cover_url             VARCHAR(512) NOT NULL,
  source                VARCHAR(512) NOT NULL,
  tags                  VARCHAR[] DEFAULT '{}',
  views_count           BIGINT NOT NULL DEFAULT 0,
  score                 INTEGER NOT NULL DEFAULT 0,
  user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS products_score_createdx ON products(score DESC, updated_at DESC);
CREATE INDEX IF NOT EXISTS products_tags_score_createdx ON products(tags, score DESC, updated_at DESC);


CREATE TABLE IF NOT EXISTS statistics (
  statistic_id          VARCHAR(36) PRIMARY KEY,
  name                  VARCHAR(512) NOT NULL,
  count                 BIGINT NOT NULL DEFAULT 0,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
