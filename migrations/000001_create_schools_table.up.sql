CREATE TABLE schools (
    school_id bigserial PRIMARY KEY,
    name text NOT NULL,
    contact text NOT NULL,
    phone serial NOT NULL,
    email citext UNIQUE NOT NULL,
    level text NOT NULL,
    website text NOT NULL,
    address serial NOT NULL,
    mode text NOT NULL,
    version serial,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

