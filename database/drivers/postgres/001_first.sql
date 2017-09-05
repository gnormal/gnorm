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
  isbn char(32) NOT NULL UNIQUE,
  booktype book_type NOT NULL,
  title text NOT NULL,
  published timestamptz[] NOT NULL,
  years integer[] NOT NULL,
  pages integer NOT NULL,
  summary text,
  available timestamptz NOT NULL DEFAULT 'NOW()',
  tags varchar[] NOT NULL DEFAULT '{}'
);

CREATE INDEX books_title_idx ON books(author_id, title);

---- create above / drop below ----

DROP TABLE IF EXISTS books CASCADE;
DROP TYPE IF EXISTS book_type CASCADE;
DROP TABLE IF EXISTS authors CASCADE;


