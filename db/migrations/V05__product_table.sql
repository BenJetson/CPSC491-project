CREATE TABLE product (
    product_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title text NOT NULL,
    description text NOT NULL,
    price int NOT NULL
);