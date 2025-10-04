-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- 1. Rename the original table
ALTER TABLE users RENAME TO users_old;

-- 2. Create a new table with the desired column type
CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY,
  login TEXT NOT NULL DEFAULT (''),
  passwd TEXT DEFAULT(''),
  name TEXT DEFAULT(''),
  email TEXT DEFAULT(''),
  active integer not null DEFAULT 0,
  is_admin integer not null DEFAULT 0,
  rem TEXT DEFAULT(''),
  unique(email)
);

-- 3. Copy data from the temporary table to the new table
INSERT INTO users (login, passwd, name, email, active, is_admin, rem)
SELECT login, passwd, name, email, active, is_admin, rem FROM users_old;

-- 4. Drop the temporary table
DROP TABLE users_old;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users RENAME TO users_old;

-- 2. Create a new table with the desired column type
CREATE TABLE users (
  login TEXT NOT NULL DEFAULT ('') PRIMARY KEY,
  passwd TEXT DEFAULT(''),
  name TEXT DEFAULT(''),
  email TEXT DEFAULT(''),
  active integer not null DEFAULT 0,
  is_admin integer not null DEFAULT 0,
  rem TEXT DEFAULT(''),
  unique(email)
);

-- 3. Copy data from the temporary table to the new table
INSERT INTO users (login, passwd, name, email, active, is_admin, rem)
SELECT login, passwd, name, email, active, is_admin, rem FROM users_old;

-- 4. Drop the temporary table
DROP TABLE users_old;
-- +goose StatementEnd
