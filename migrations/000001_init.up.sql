CREATE TABLE pull_requests
(
    id                UUID DEFAULT uuidv7() PRIMARY KEY,
    -- id                SERIAL PRIMARY KEY,
    author_id         UUID,
    -- created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    merged_at         TIMESTAMP,
    -- pull_request_id   VARCHAR   NOT NULL,
    pull_request_name VARCHAR NOT NULL,
    status            INTEGER,
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (status) REFERENCES pr_status (id)

);

-- создаем отдельную таблицу для статусов
-- при таком подходе легко изменять названия или добавлять новые
CREATE TABLE pr_status
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE reviewers
(
    pr_id   UUID,
    user_id UUID,
    FOREIGN KEY (pr_id) REFERENCES pull_requests (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (pr_id, user_id)
);

CREATE TABLE teams
(
    id   UUID DEFAULT uuidv7() PRIMARY KEY,
    -- id   SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE users
(
    id        UUID DEFAULT uuidv7() PRIMARY KEY,
    -- id        SERIAL PRIMARY KEY,
    is_active BOOL    NOT NULL,
    user_id   VARCHAR NOT NULL UNIQUE,
    username  VARCHAR NOT NULL UNIQUE
);

CREATE TABLE team_member
(
    team_id UUID,
    user_id UUID,
    FOREIGN KEY (team_id) REFERENCES teams (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    PRIMARY KEY (team_id, user_id)
)

