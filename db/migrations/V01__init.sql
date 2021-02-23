CREATE TABLE users(
    user_id serial PRIMARY KEY,
    last_name varchar(50) NOT NULL,
    first_name varchar(50) NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    role_id integer NOT NULL
)

CREATE TABLE roles(
    role_id serial PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    title varchar(50) NOT NULL
)
