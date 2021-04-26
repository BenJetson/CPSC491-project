# Database Privileges

To get the application working, it is necessary to ensure that the database
accounts have the correct level of access.

The following script is useful, pusuant to the existing instructions from
[`server_config.md`](server_config.md).

## For Development Tier

Run this inside the `public` schema for the `dev` database:

```sql
ALTER DEFAULT PRIVILEGES
IN SCHEMA public
GRANT
    INSERT,
    SELECT,
    UPDATE,
    DELETE
ON TABLES
TO cpsc491_dev;

ALTER DEFAULT PRIVILEGES
IN SCHEMA public
GRANT
    EXECUTE
ON FUNCTIONS
TO cpsc491_dev;

GRANT
    SELECT,
    INSERT,
    UPDATE,
    DELETE
ON ALL TABLES
IN SCHEMA public
TO cpsc491_dev;

GRANT
    EXECUTE
ON ALL FUNCTIONS
IN SCHEMA public
TO cpsc491_dev;
```

## For Production Tier

Run this inside the `public` schema for the `prod` database:

```sql
ALTER DEFAULT PRIVILEGES
IN SCHEMA public
GRANT
    INSERT,
    SELECT,
    UPDATE,
    DELETE
ON TABLES
TO cpsc491_prod;

ALTER DEFAULT PRIVILEGES
IN SCHEMA public
GRANT
    EXECUTE
ON FUNCTIONS
TO cpsc491_prod;

GRANT
    SELECT,
    INSERT,
    UPDATE,
    DELETE
ON ALL TABLES
IN SCHEMA public
TO cpsc491_prod;

GRANT
    EXECUTE
ON ALL FUNCTIONS
IN SCHEMA public
TO cpsc491_prod;
```
