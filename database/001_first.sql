DROP TABLE IF EXISTS books CASCADE;
DROP TYPE IF EXISTS book_type CASCADE;
DROP TABLE IF EXISTS authors CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE authors (
  id uuid DEFAULT uuid_generate_v4() NOT NULL primary key,
  name text NOT NULL
);

CREATE INDEX authors_name_idx ON authors(name);

CREATE TYPE book_type AS ENUM (
  'FICTION',
  'NONFICTION'
);

CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  author_id uuid NOT NULL REFERENCES authors(id),
  isbn text NOT NULL UNIQUE,
  booktype book_type NOT NULL,
  title text NOT NULL,
  years integer[] NOT NULL,
  available timestamp with time zone NOT NULL DEFAULT 'NOW()',
  tags varchar[] NOT NULL DEFAULT '{}'
);

CREATE INDEX books_title_idx ON books(author_id, title);

---- create above / drop below ----

DROP TABLE IF EXISTS books CASCADE;
DROP TYPE IF EXISTS book_type CASCADE;
DROP TABLE IF EXISTS authors CASCADE;


