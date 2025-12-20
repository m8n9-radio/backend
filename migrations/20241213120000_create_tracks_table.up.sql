-- Migration up: Create tracks table
CREATE TABLE tracks (
  id CHAR(32) NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  cover TEXT NOT NULL,
  rotate INTEGER NOT NULL DEFAULT 0,
  likes INTEGER NOT NULL DEFAULT 0,
  listeners INTEGER NOT NULL DEFAULT 0,
  dislikes INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracks_created_at ON tracks (created_at DESC);
