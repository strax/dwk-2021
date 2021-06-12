-- Add migration script here
CREATE TABLE pings (
    id uuid PRIMARY KEY,
    ts timestamp DEFAULT now()
);