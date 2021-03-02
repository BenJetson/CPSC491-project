#Best practices for PostgreSQL

These notes cover some of the best practices that have been discovered while working with PostgreSQL

- When creating a table, make the name a singular noun. For example, when creating the user table, the table name should be
named "user" since a row of the table corresponds to one user, not multiple.
- When creating new applications, **do not** use *serial* as a type. Instead use this identity column generator:
"PRIMARY KEY GENERATED ALWAYS AS IDENTITY"
- When using the type *text*, **do not** use *varchar(n)*.
    - *varchar(n)* is only used when you need to receive an error message if the inputed text is longer than the limit
    - *text* allows for arbitary length
- If there is a reference to Table B in Table A, Table B should be declared first in order to be able to be referenced to
- When a reference from another table is made, a relationship needs to be generated between the item and the table. In order
to do this, foreign keys need to be used:
`item_name type NOT NULL REFERENCES table_name(item_name) ON DELETE RESTRICT`
     - This makes a reference to the *item_name* column of the *table_name* table of the row. 
     - If there is an attempt to delete a row from *table_name* while it is still referenced in the other table, it will 
        stop the deletion from happening until all references are deleted first.