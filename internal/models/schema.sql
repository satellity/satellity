CREATE TABLE IF NOT EXISTS users (
  user_id               VARCHAR(36) PRIMARY KEY,
  public_key            VARCHAR(512),
  email                 VARCHAR(512),
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
CREATE UNIQUE INDEX IF NOT EXISTS users_public_keyx ON users ((LOWER(public_key)));
CREATE INDEX IF NOT EXISTS users_createdx ON users (created_at);


CREATE TABLE IF NOT EXISTS sessions (
  session_id            VARCHAR(36) PRIMARY KEY,
  user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  public_key            VARCHAR(128) NOT NULL UNIQUE,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS sessions_user_createdx ON sessions (user_id, created_at DESC);


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
  liked_at              TIMESTAMP WITH TIME ZONE,
  bookmarked_at         TIMESTAMP WITH TIME ZONE,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (topic_id, user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS topic_users_reversex ON topic_users(user_id, topic_id);
CREATE INDEX IF NOT EXISTS topic_users_likedx ON topic_users(topic_id, liked_at);
CREATE INDEX IF NOT EXISTS topic_users_bookmarkedx ON topic_users(topic_id, bookmarked_at);


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


CREATE TABLE IF NOT EXISTS sources (
  source_id             VARCHAR(36) PRIMARY KEY,
  author                VARCHAR(512) NOT NULL,
  host                  VARCHAR(128) NOT NULL,
  link                  VARCHAR(1024) NOT NULL UNIQUE,
  logo_url              VARCHAR(1024) NOT NULL,
  locality              VARCHAR(128) NOT NULL,
  wreck                 INTEGER NOT NULL DEFAULT 0,
  published_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS sources_updatedx ON sources (updated_at);


CREATE TABLE IF NOT EXISTS gists (
  gist_id               VARCHAR(36) PRIMARY KEY,
  identity              VARCHAR(256) NOT NULL UNIQUE,
  author                VARCHAR(512) NOT NULL,
  title                 VARCHAR(1024) NOT NULL DEFAULT '',
  source_id             VARCHAR(36) NOT NULL,
  genre                 VARCHAR(128) NOT NULL,
  cardinal              BOOL NOT NULL DEFAULT false,
  link                  VARCHAR(256) NOT NULL UNIQUE,
  body                  TEXT NOT NULL,
  publish_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS gists_cardinalx ON gists (cardinal, publish_at DESC);
CREATE INDEX IF NOT EXISTS gists_genrex ON gists (genre, publish_at DESC);
CREATE INDEX IF NOT EXISTS gists_publishx ON gists (publish_at DESC);


CREATE TABLE IF NOT EXISTS statistics (
  statistic_id          VARCHAR(36) PRIMARY KEY,
  name                  VARCHAR(512) NOT NULL,
  count                 BIGINT NOT NULL DEFAULT 0,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
