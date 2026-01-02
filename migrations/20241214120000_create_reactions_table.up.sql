-- Migration up: Create reactions table
CREATE TABLE reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id CHAR(255) NOT NULL,
    track_id CHAR(32) NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    reaction VARCHAR(10) NOT NULL CHECK (reaction IN ('like', 'dislike')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, track_id)
);

CREATE INDEX idx_reactions_user_track ON reactions(user_id, track_id);
CREATE INDEX idx_reactions_track_id ON reactions(track_id);
