CREATE TABLE if not exists users
(
    id       SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    coins    INT DEFAULT 1000
);

CREATE TABLE if not exists inventory
(
    user_id  INT REFERENCES users (id),
    item     TEXT NOT NULL,
    quantity INT DEFAULT 0,
    PRIMARY KEY (user_id, item)
);

CREATE TABLE if not exists transactions
(
    id        SERIAL PRIMARY KEY,
    from_user INT,
    to_user   INT,
    amount    INT,
    time      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user) REFERENCES users (id),
    FOREIGN KEY (to_user) REFERENCES users (id)
);