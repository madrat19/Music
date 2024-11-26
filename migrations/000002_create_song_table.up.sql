CREATE TABLE IF NOT EXISTS "Song" (
    song_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    text TEXT,
    link VARCHAR(2048),
    group_id INT,
    CONSTRAINT fk_group FOREIGN KEY (group_id)
        REFERENCES "Group" (group_id)
        ON DELETE SET NULL
);
