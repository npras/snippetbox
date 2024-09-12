/* psql commands

-- connet as user
psql -U postgres

-- list all dbs
\l

-- connect to a db
\c db_name

-- list tables of schema "public"
\dt

-- list tables of schema "any_other_schema"
\dt any_other_schema.*

-- change current schema to "any_other_schema"
SET search_path TO any_other_schema;

-- show table details (indexes, foreign keys, references)
\d table_name

 */

-- create db
CREATE DATABASE snippetbox
WITH 
ENCODING 'UTF8' 
LC_COLLATE='en_US.UTF-8' 
LC_CTYPE='en_US.UTF-8'
TEMPLATE template0;

-- create table
CREATE TABLE snippets (
  id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL
);

-- create index
CREATE INDEX idx_snippets_created_at ON snippets(created_at);

-- insert records
INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
  'An old silent pond',
  E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
  CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
  (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
  'Over the wintry forest',
  E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
  CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
  (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
  'First autumn morning',
  E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
  CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
  (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + INTERVAL '7 days'
);

-- unrun, untested
-- Create a new user
CREATE USER web WITH PASSWORD 'golanger123456';

-- Grant privileges on all tables in the schema
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO web;

-- Grant privileges on future tables in the schema
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO web;

-- If you need to change the password later
ALTER USER web WITH PASSWORD 'golanger1234567';
