CREATE TABLE users
(
    id         uuid primary key,
    auth_id    text not null,
    email      varchar,
    password   VARCHAR,
    name       varchar,
    role       varchar,
    image_url  text,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

CREATE TABLE categories
(
    id         uuid primary key,
    name       varchar,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

Create Table news
(
    id          uuid primary key,
    author      text,
    title       text,
    description text,
    content     text,
    url         text,
    image_url   text,
    publish_at  timestamp,
    created_at  timestamp,
    updated_at  timestamp,
    deleted_at  timestamp
);

CREATE TABLE has_categories
(
    news_id     uuid REFERENCES news (id)       NOT NULL,
    category_id uuid references categories (id) NOT NULL,
    PRIMARY KEY (news_id, category_id)
);

CREATE TABLE likes
(
    news_id uuid REFERENCES news (id)  NOT NULL,
    user_id uuid references users (id) NOT NULL,
    PRIMARY KEY (news_id, user_id)
);

CREATE TABLE dislikes
(
    news_id uuid references news (id)  NOT NULL,
    user_id uuid references users (id) NOT NULL,
    PRIMARY KEY (news_id, user_id)
);

CREATE TABLE saves
(
    news_id uuid references news (id)  NOT NULL,
    user_id uuid references users (id) NOT NULL,
    PRIMARY KEY (news_id, user_id)
);

CREATE TABLE comments
(
    news_id uuid references news (id)  NOT NULL,
    user_id uuid references users (id) NOT NULL,
    text text ,
    published_at timestamp
);

CREATE TABLE views
(
    news_id uuid references news (id)  NOT NULL,
    user_id uuid references users (id) NOT NULL,
    PRIMARY KEY (news_id, user_id)
);