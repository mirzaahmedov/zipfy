CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4(),
  name text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  created timestamp NOT NULL DEFAULT now()
);
ALTER TABLE users ADD CONSTRAINT user_id_pkey PRIMARY KEY (id);
ALTER TABLE users ADD CONSTRAINT unique_user_email UNIQUE (email);

CREATE TABLE urls (
  id SERIAL,
  url text NOT NULL,
  user_id uuid NOT NULL,
  created timestamp NOT NULL DEFAULT now()
);
ALTER TABLE urls ADD CONSTRAINT url_id_pkey PRIMARY KEY (id);
ALTER TABLE urls ADD CONSTRAINT url_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);
