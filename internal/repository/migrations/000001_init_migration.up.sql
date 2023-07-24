CREATE TABLE film(
    id uuid PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    genre VARCHAR(255),
    release_year INT NOT NULL,
    revenue FLOAT8 NOT NULL,
    hashcode text NOT NUll
);

CREATE index idx_film_title ON film(title);

CREATE TABLE producer(
    id uuid PRIMARY KEY ,
    name VARCHAR(255) NOT NULL,
    birthdate DATE not null,
    gender VARCHAR(1) NOT NULL
);

CREATE index idx_producer_name on producer(name);

CREATE TABLE actor(
    id uuid PRIMARY KEY ,
    name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    gender VARCHAR(1) NOT NULL,
    role VARCHAR(255) NOT NULL
);

CREATE index idx_actor_name on actor(name);


CREATE TABLE crew(
    actor_id uuid,
    producer_id uuid,
    film_id uuid,
    FOREIGN KEY (actor_id) REFERENCES actor(id),
    FOREIGN KEY (producer_id) REFERENCES producer(id),
    FOREIGN KEY (film_id) REFERENCES film(id)
);

