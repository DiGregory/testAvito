CREATE TABLE IF NOT EXISTS adverts (
  id serial PRIMARY KEY,
  name  text,
  description text,
  images text[],
  price real,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  deleted_at timestamptz
);



