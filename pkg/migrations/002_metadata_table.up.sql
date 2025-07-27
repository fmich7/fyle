CREATE TABLE IF NOT EXISTS file_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    filename TEXT NOT NULL,
    location TEXT NOT NULL,
    description TEXT,
    tags TEXT[],
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);