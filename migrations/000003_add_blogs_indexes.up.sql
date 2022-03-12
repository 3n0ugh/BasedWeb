CREATE INDEX IF NOT EXISTS blogs_title_idx ON blogs USING GIN
    (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS blogs_category_idx ON blogs USING GIN (category);
