CREATE TABLE IF NOT EXISTS users (
    id text(36) not null,
    username text,
    email text,
    password text,
    is_email_verified boolean default false,
    role enum('manager','tech') default 'tech',
    signature text,
    manager_id text(36),
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    PRIMARY KEY (id(36))
);

CREATE TABLE IF NOT EXISTS tasks (
    id text(36) not null,
    summary text(2500),
    user_id text(36),
    performed_at timestamp,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    PRIMARY KEY (id(36))
);

-- MYSQL doesn't likes uuids
-- ALTER TABLE tasks ADD CONSTRAINT user_tasks_fk FOREIGN KEY (user_id(36)) REFERENCES users(id);