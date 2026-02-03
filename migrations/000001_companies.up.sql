CREATE TABLE companies (
    id UUID PRIMARY KEY,
    name VARCHAR(15) NOT NULL UNIQUE,
    description VARCHAR(3000) NOT NULL,
    employees_count INTEGER NOT NULL,
    registered BOOLEAN NOT NULL DEFAULT false,
    type SMALLINT NOT NULL
);