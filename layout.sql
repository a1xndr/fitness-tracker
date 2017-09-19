-- user

PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE user (
    id integer primary key,
    username varchar(100),
    password_hashed varchar(1000),
    email varchar(100),
    created_at datetime,
    disabled bool
);
COMMIT;

-- exercise

BEGIN TRANSACTION;
CREATE TABLE exercise (
    id integer primary key,
    name varchar(1000),
    type references exercisemetric(id),
    description varchar(1000),
    reps integer,
    weight integer,
    time integer,
    speed integer,
    grade integer
);
COMMIT;
-- exercisetype

BEGIN TRANSACTION;
CREATE TABLE exercisemetric (
    id integer primary key,
    name varchar(1000),
    reps integer,
    weight integer,
    seconds integer,
    speed integer,
    grade integer,
    note integer,
    author references user(id)
);
COMMIT;

-- workout
BEGIN TRANSACTION;
CREATE TABLE workout (
    id integer primary key,
    date datetime,
    user references user(id)
);

COMMIT;
-- sets

BEGIN TRANSACTION;
CREATE TABLE sets (
    id integer primary key,
    workout references workout(id),
    exercise references exercise(id),
    reps integer,
    weight real,
    seconds integer
    note varchar(1000)
);

COMMIT;


