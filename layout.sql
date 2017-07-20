-- user

PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE user (
    id integer primary key autoincrement,
    username varchar(100),
    password varchar(1000),
    email varchar(100)
);
COMMIT;

-- exercise

BEGIN TRANSACTION;
CREATE TABLE exercise (
    id integer primary key autoincrement,
    name varchar(1000),
    description varchar(1000)
);

COMMIT;
-- workout
BEGIN TRANSACTION;
CREATE TABLE workout(
    id integer primary key autoincrement,
    date datetime
);

COMMIT;
-- sets

BEGIN TRANSACTION;
CREATE TABLE sets (
    id integer primary key autoincrement,
    workout references workout(id),
    exercise references exercise(id),
    reps integer,
    weight real,
    seconds integer
);

COMMIT;


