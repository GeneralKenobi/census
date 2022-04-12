CREATE SCHEMA censusdb;
SET search_path TO censusdb;


CREATE TABLE person
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL CHECK (name <> ''),
    surname          VARCHAR(255) NOT NULL CHECK (surname <> ''),
    email            VARCHAR(255) NOT NULL CHECK (email <> ''),
    date_of_birth    DATE         NOT NULL,
    hobby            VARCHAR(255) NOT NULL,
    created_at       TIMESTAMP    NOT NULL,
    last_modified_at TIMESTAMP    NOT NULL,

    CONSTRAINT unique_email UNIQUE (email)
);
CREATE INDEX person_email ON person (email);
