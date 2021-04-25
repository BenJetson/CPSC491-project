
ALTER TABLE product
    ADD COLUMN vendor_id int NOT NULL,
    ADD COLUMN organization_id int NOT NULL
        REFERENCES organization (organization_id)
        ON DELETE CASCADE,
    ADD COLUMN image_url text,
    ADD COLUMN is_available boolean NOT NULL DEFAULT TRUE,

    -- https://pganalyze.com/blog/full-text-search-ruby-rails-postgres
    ADD COLUMN searchable tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description,'')), 'B')
    ) STORED,

    ADD UNIQUE(vendor_id, organization_id)
;
