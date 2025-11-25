CREATE TABLE IF NOT EXISTS users
(
    -- id        UUID DEFAULT uuidv7() PRIMARY KEY,
    id        SERIAL PRIMARY KEY,
    is_active BOOL    NOT NULL,
    user_id   VARCHAR NOT NULL UNIQUE,
    username  VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS teams
(
    -- id   UUID DEFAULT uuidv7() PRIMARY KEY,
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS team_member
(
    team_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY (team_id) REFERENCES teams (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (team_id, user_id)
);

-- create a separate table for statuses
-- with this approach, it's easy to change names or add new ones
CREATE TABLE IF NOT EXISTS pr_status
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    -- id                UUID DEFAULT uuidv7() PRIMARY KEY,
    id                SERIAL PRIMARY KEY,
    author_id         INTEGER,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merged_at         TIMESTAMP,
    pull_request_id   VARCHAR   NOT NULL UNIQUE,
    pull_request_name VARCHAR   NOT NULL,
    status            INTEGER   NOT NULL REFERENCES pr_status (id),
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (status) REFERENCES pr_status (id)

);

CREATE TABLE IF NOT EXISTS reviewers
(
    pr_id   INTEGER,
    user_id INTEGER,
    FOREIGN KEY (pr_id) REFERENCES pull_requests (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (pr_id, user_id)
);

CREATE INDEX idx_reviewers_user_id ON reviewers (user_id);
CREATE INDEX idx_users_active ON users (is_active);
CREATE INDEX idx_team_member_team_id ON team_member (team_id);

-- add data
INSERT INTO pr_status (name)
VALUES ('OPEN'),
       ('MERGED');

