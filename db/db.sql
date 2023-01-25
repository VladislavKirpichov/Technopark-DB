DROP SCHEMA IF EXISTS public CASCADE;
CREATE SCHEMA public;

SET client_encoding = 'UTF8';
SET client_min_messages = warning;
SET lock_timeout = 0;
SET statement_timeout = 0;

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE users
(
    nickname CITEXT NOT NULL PRIMARY KEY,
    email    CITEXT NOT NULL UNIQUE,
    fullname TEXT NOT NULL,
    about    TEXT NOT NULL
);

CREATE UNLOGGED TABLE forums
(
    id      SERIAL NOT NULL UNIQUE,
    slug    CITEXT NOT NULL PRIMARY KEY,
    title   TEXT NOT NULL,
    "user"  CITEXT NOT NULL REFERENCES users (nickname),
    posts   INTEGER DEFAULT 0,
    threads INTEGER DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS forum_users
(
    forum    CITEXT NOT NULL REFERENCES forums (slug),
    nickname CITEXT NOT NULL REFERENCES users (nickname),

    PRIMARY KEY (forum, nickname)
);

CREATE UNLOGGED TABLE threads
(
    id      SERIAL NOT NULL PRIMARY KEY,
    slug    CITEXT UNIQUE,
    author  CITEXT NOT NULL REFERENCES users (nickname),
    forum   CITEXT NOT NULL REFERENCES forums (slug),
    title   TEXT NOT NULL,
    message TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    votes   INTEGER DEFAULT 0
);

CREATE UNLOGGED TABLE votes
(
    nickname CITEXT NOT NULL REFERENCES users(nickname),
    thread   INTEGER NOT NULL REFERENCES threads(id),
    value    INTEGER NOT NULL,

    PRIMARY KEY (thread, nickname)
);

CREATE UNLOGGED TABLE posts
(
    id       SERIAL NOT NULL PRIMARY KEY,
    parent   INTEGER DEFAULT NULL,
    author   CITEXT NOT NULL REFERENCES users (nickname),
    forum    CITEXT NOT NULL REFERENCES forums (slug),
    thread   INTEGER NOT NULL REFERENCES threads (id),
    created  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    isEdited BOOLEAN NOT NULL DEFAULT false,
    message  TEXT NOT NULL,
    path     INTEGER[] NOT NULL
);

CREATE OR REPLACE FUNCTION makeVote() RETURNS TRIGGER LANGUAGE plpgsql AS
$$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.value
    WHERE id = NEW.thread;

    RETURN NULL;
END;
$$;

CREATE TRIGGER insertVote
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE makeVote();

CREATE OR REPLACE FUNCTION updateVote() RETURNS TRIGGER LANGUAGE plpgsql AS
$$
BEGIN
    UPDATE threads
    SET votes = votes - OLD.value + NEW.value
    WHERE id = NEW.thread;

    RETURN NULL;
END;
$$;

CREATE TRIGGER updateVote
    AFTER UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE updateVote();

CREATE OR REPLACE FUNCTION newPath() RETURNS TRIGGER LANGUAGE plpgsql AS
$$
DECLARE
    newPath INTEGER[];
BEGIN
    IF NEW.parent IS NULL THEN
        NEW.path := NEW.path || NEW.id;
    ELSE
        SELECT INTO newPath path FROM posts WHERE id = NEW.parent AND thread = NEW.thread;

        IF (newPath[1] IS NULL) THEN
            RAISE EXCEPTION 'parent empty';
        END IF;

        NEW.path := NEW.path || newPath || NEW.id;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER newPath
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE newPath();

CREATE OR REPLACE FUNCTION threadsCounter() RETURNS TRIGGER LANGUAGE plpgsql AS
$$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;

    RETURN NULL;
END;
$$;

CREATE TRIGGER threadsCounter
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE threadsCounter();

CREATE OR REPLACE FUNCTION postsCounter() RETURNS TRIGGER LANGUAGE plpgsql AS
$$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;

    RETURN NEW;
END;
$$;

CREATE TRIGGER postsCounter
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE postsCounter();

CREATE OR REPLACE FUNCTION addUser() RETURNS TRIGGER LANGUAGE plpgsql AS
$$
BEGIN
    INSERT INTO forum_users (forum, nickname)
    VALUES (NEW.forum, NEW.author)
    ON CONFLICT do nothing;
    RETURN NULL;
END;
$$;

CREATE TRIGGER addUserByForum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE addUser();

CREATE TRIGGER addUserByPosts
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE addUser();

CREATE INDEX IF NOT EXISTS sortUsers ON forum_users (nickname);
CREATE INDEX IF NOT EXISTS sortForumsAndTime ON threads (forum, created);
CREATE INDEX IF NOT EXISTS sortUsers ON users (nickname, email);

CREATE INDEX IF NOT EXISTS sortThreadsAndId ON posts (thread, id);
CREATE INDEX IF NOT EXISTS sortThreadsAndPath ON posts (thread, path);
CREATE INDEX IF NOT EXISTS sortThreadsAndParent ON posts (thread, (path[1]));

VACUUM ANALYZE;
