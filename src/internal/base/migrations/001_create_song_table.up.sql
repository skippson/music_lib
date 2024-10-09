-- +goose Up
create table if not exists songs (
    id serial PRIMARY KEY,
    group_name varchar(255),
    song varchar(255),
    release_date varchar(255),
    lyrics varchar(255)
);