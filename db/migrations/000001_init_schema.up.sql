CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS songs (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     name VARCHAR(100) NOT NULL,
                                     release_date TIMESTAMP NOT NULL,
                                     link TEXT NOT NULL,
                                     text TEXT NOT NULL,
                                     group_name TEXT NOT NULL
);