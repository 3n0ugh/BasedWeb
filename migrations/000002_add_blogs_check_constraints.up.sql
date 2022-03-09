ALTER TABLE blogs ADD CONSTRAINT category_length_check CHECK
    (array_length(category, 1) BETWEEN 1 AND 5);
ALTER TABLE blogs ADD CONSTRAINT title_length_check CHECK
    (char_length(title) BETWEEN 3 AND 70);