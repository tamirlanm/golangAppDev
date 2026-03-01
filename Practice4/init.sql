

CREATE TABLE IF NOT EXISTS movies (
    id          VARCHAR(50)  PRIMARY KEY,
    genre       VARCHAR(100) NOT NULL,
    budget      INTEGER      NOT NULL,
    title       VARCHAR(200) NOT NULL,
    actors JSONB        DEFAULT '[]'
);


INSERT INTO movies (id, genre, budget, title, actors)
VALUES
    ('1', 'Crime drama',  500000,  'The Godfather',  '["Al Pacino", "Marlon Brando"]'),
    ('2', 'Prison drama', 1000000, 'The Shawshank Redemption', '["Tim Robbins", "Morgan Freeman"]')
ON CONFLICT (id) DO NOTHING;