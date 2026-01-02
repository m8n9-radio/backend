-- Migration up: Create listeners table
CREATE TABLE listeners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id CHAR(255) NOT NULL,
    track_id CHAR(32) NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, track_id)
);

-- Create performance indexes
CREATE INDEX idx_listeners_track_id ON listeners(track_id);
CREATE INDEX idx_listeners_user_id ON listeners(user_id);
CREATE INDEX idx_listeners_created_at ON listeners(created_at DESC);
