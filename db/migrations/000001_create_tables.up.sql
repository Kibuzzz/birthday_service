BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    birthday DATE NOT NULL
);

CREATE TABLE subscriptions (
    subscriber_id INT NOT NULL,
    birthday_person_id INT NOT NULL,
    notification_time TIMESTAMP NOT NULL,
    PRIMARY KEY (subscriber_id, birthday_person_id),
    FOREIGN KEY (subscriber_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (birthday_person_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMIT;