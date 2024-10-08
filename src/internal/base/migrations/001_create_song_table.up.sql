-- +goose Up
CREATE TABLE if not exists songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255),
    song VARCHAR(255)
);

-- +goose Down
-- drop table songs