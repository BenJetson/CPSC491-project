CREATE TABLE product (
    product_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title text NOT NULL,
    product_description text NOT NULL,
    price decimal(7,2)
);