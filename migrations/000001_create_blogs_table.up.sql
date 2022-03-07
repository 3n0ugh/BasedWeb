CREATE TABLE IF NOT EXISTS blogs (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    body text NOT NULL,
    category text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);