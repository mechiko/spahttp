-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- Create users table
CREATE TABLE if not exists updates (
  id integer NOT null PRIMARY KEY,
  message    text NOT null default(''),
  sender_id integer NOT null default(0),
  chat_id integer NOT null default(0),
  recepient string NOT null default(''),
  'update'    text NOT null default(''),
  created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M','now', 'localtime'))
);

CREATE TABLE if not exists telebotusers (
  id integer NOT null PRIMARY KEY,
  first_name    text NOT null default(''),
  last_name     text NOT null default(''),
  username     text NOT null default(''),
  language_code text NOT null default(''),
  is_bot        integer NOT null default(0),
  is_premium    integer NOT null default(0),
  added_to_menu  integer NOT null default(0),
  -- Returns only in getMe
  can_join_groups   integer NOT null default(0),
  can_read_messages integer NOT null default(0),
  supports_inline  integer NOT null default(0),
  is_admin integer NOT null default(0),
  created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M','now', 'localtime')),
  masters text NOT null default('')
);

CREATE TABLE if not exists chats (
  id integer NOT null  PRIMARY KEY,
  type text NOT null default(''),
  title text NOT null default(''),
  first_name    text NOT null default(''),
  last_name     text NOT null default(''),
  username     text NOT null default(''),
  bio              text NOT null default(''),
  photo            text NOT null default(''),
  description      text NOT null default(''),
  invite_link       text NOT null default(''),
  pinned_message    text NOT null default(''),
  permissions      text NOT null default(''),
  slow_mode         integer NOT null default(0),
  sticker_set       text NOT null default(''),
  can_set_sticker_set integer NOT null default(0),
  linked_chat_id     integer NOT null default(0),
  chat_location     text NOT null default(''),
  private          integer NOT null default(0),
  protected        integer NOT null default(0),
  no_voice_and_video  integer NOT null default(0),
  created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M','now', 'localtime'))
);

-- активна та которая в боте создана и манипулируется
-- при каждом обращении считывания активной проверяется срок начала и окончания
-- при завершении командировки она снимается с активной
-- может существовать только одна активная для recepient_id
-- masters список начальников куда отправлять стату
CREATE TABLE if not exists missions (
  id INTEGER NOT NULL PRIMARY KEY,
  recepient_id integer NOT null,
  start text NOT null default(''),
  end text NOT null default(''),
  place text NOT null default(''),
  departament text NOT null default(''),
  target text NOT null default(''),
  rem text NOT null default(''),
  reported integer NOT null default(0),
  active integer NOT null default(0),
  created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M','now', 'localtime'))
);

-- модуль app по умолчанию на все приложение
CREATE TABLE if not exists app_state (
  module TEXT NOT NULL DEFAULT 'app',
  key TEXT NOT NULL DEFAULT '',
  value TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (module, key)
);

-- экзамены пользователя по сути мап[key]value
-- где ключ ОТ ППБ ОАЭ_АС_ДИ ФНП ПРБ ЭБ Медосмотр
CREATE TABLE if not exists user_states (
  user_id INTEGER NOT NULL DEFAULT (0),
  key TEXT NOT NULL DEFAULT '',
  value TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (user_id, key)
);

-- где ключ ОТ ППБ ОАЭ_АС_ДИ ФНП ПРБ ЭБ Медосмотр
-- telephone fio is_admin email
CREATE TABLE if not exists state_key (
  key TEXT NOT NULL DEFAULT '' PRIMARY KEY,
  is_examen integer NOT null default(0),
  is_intro integer NOT null default(0),
  description TEXT NOT NULL DEFAULT ''
);

-- таблица логинов клиентов на всякий пожарный чтобы связь с клиентом была
CREATE TABLE if not exists users (
  login TEXT NOT NULL DEFAULT ('') ,
  passwd TEXT DEFAULT(''),
  name TEXT DEFAULT(''),
  email TEXT DEFAULT(''),
  active integer not null DEFAULT 0,
  is_admin integer not null DEFAULT 0,
  rem TEXT DEFAULT(''),
  unique(email)
);

CREATE TABLE if not exists examen_ended (
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL DEFAULT (0),
  key TEXT NOT NULL DEFAULT '',
  date TEXT NOT NULL DEFAULT '',
  created TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M','now', 'localtime'))
);

INSERT OR REPLACE INTO state_key(key, is_examen, is_intro, description) VALUES
('ОТ', 1, 0 ,'ОТ'),
('ППБ', 1, 0 ,'ППБ'),
('ОПЭ АС, ДИ', 1, 0 ,'ОПЭ АС, ДИ'),
('ФНП', 1, 0 ,'ФНП'),
('ПРБ', 1, 0 ,'ПРБ'),
('ЭБ', 1, 0 ,'ЭБ'),
('Медосмотр', 1, 0 ,'Медосмотр'),
('ИДЕНТ', 0, 1 ,'идентификатор');

INSERT OR REPLACE INTO users(login, passwd, name, email, active, is_admin, rem) VALUES 
	('kbprime@mail.ru', '$2a$08$PjVaxr7RxjrHWpW0laJwcuKISPNPSOxCkbHCxZvHJkQKXsQUHDKA2', 'mikl', 'kbprime@mail.ru', 1, 1, 'author'),
	('a.kuleshov.m@gmail.com', '$2a$08$PjVaxr7RxjrHWpW0laJwcuKISPNPSOxCkbHCxZvHJkQKXsQUHDKA2', 'Admin', 'a.kuleshov.m@gmail.com', 1, 1, 'na4alnik'),
  ('n91n91@mail.ru', '$2a$08$PjVaxr7RxjrHWpW0laJwcuKISPNPSOxCkbHCxZvHJkQKXsQUHDKA2', 'nastya', 'n91n91@mail.ru', 1, 1, 'info');

-- Index to speed up expiry-based look-ups / sweeps

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
-- Drop indexes
-- Drop tables
DROP TABLE IF EXISTS updates;
DROP TABLE IF EXISTS telebotusers;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS missions;
DROP TABLE IF EXISTS app_state;
DROP TABLE IF EXISTS user_states;
DROP TABLE IF EXISTS state_key;
DROP TABLE IF EXISTS examen_ended;
DROP TABLE IF EXISTS users;
