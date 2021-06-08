CREATE TABLE todo(
    id uuid PRIMARY KEY,
    text text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
