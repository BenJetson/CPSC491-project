
INSERT INTO person (
    first_name, last_name, email, role_id,
    pass_hash
) VALUES
( -- ID: 1
    'Ben', 'Admin', 'bfgodfr+admin@g.clemson.edu', 1,
    '$2a$10$y9Ia6Oh/tyZbWGkLwQxb2u3gOJ0SdmvpvAMz5PlE.3fgbNYy6GBmu'
),
( -- ID: 2
    'Ben', 'Sponsor', 'bfgodfr+sponsor@g.clemson.edu', 2,
    '$2a$10$y9Ia6Oh/tyZbWGkLwQxb2u3gOJ0SdmvpvAMz5PlE.3fgbNYy6GBmu'
),
( -- ID: 3
    'Ben', 'Driver', 'bfgodfr+driver@g.clemson.edu', 4,
    '$2a$10$y9Ia6Oh/tyZbWGkLwQxb2u3gOJ0SdmvpvAMz5PlE.3fgbNYy6GBmu'
),
( -- ID: 4
    'Cynthia', 'Admin', 'rbrazil+admin@g.clemson.edu', 1,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 5
    'Cynthia', 'Sponsor', 'rbrazil+sponsor@g.clemson.edu', 2,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 6
    'Cynthia', 'Driver', 'rbrazil+driver@g.clemson.edu', 4,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 7
    'Chloe', 'Admin', 'ccaples+admin@g.clemson.edu', 1,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 8
    'Chloe', 'Sponsor', 'ccaples+sponsor@g.clemson.edu', 2,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 9
    'Chloe', 'Driver', 'ccaples+driver@g.clemson.edu', 4,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 10
    'Cameron', 'Admin', 'sharpe9+admin@g.clemson.edu', 1,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 11
    'Cameron', 'Sponsor', 'sharpe9+sponsor@g.clemson.edu', 2,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
),
( -- ID: 12
    'Cameron', 'Driver', 'sharpe9+driver@g.clemson.edu', 4,
    '$2a$10$XnZvNenRRA00ufEtl2fs8Ol5Jh0cntDJkxsu9o.oDqXT78pd6UFrG'
);

INSERT INTO organization (
    name, point_value
) VALUES
( -- ID: 1
    'XIV LLC', 1
);

INSERT INTO affiliation (
    person_id, organization_id, points
) VALUES
(
    2, 1, NULL
),
(
    3, 1, 0
),
(
    5, 1, NULL
),
(
    6, 1, 0
),
(
    8, 1, NULL
),
(
    9, 1, 0
),
(
    11, 1, NULL
),
(
    12, 1, 0
);
