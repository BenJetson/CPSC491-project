
CREATE TABLE application (
    application_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    applicant_id int NOT NULL
        REFERENCES person(person_id)
        ON DELETE CASCADE,
    organization_id int NOT NULL
        REFERENCES organization(organization_id)
        ON DELETE CASCADE,
    comment text NOT NULL,
    approved boolean,
    reason text,
    created_at timestamptz NOT NULL,
    approved_at timestamptz
);
