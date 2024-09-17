CREATE TABLE
  public.users (
    id serial NOT NULL,
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    hash character varying(255) NOT NULL,
    num_reports integer NOT NULL DEFAULT 0
  );

ALTER TABLE
  public.users
ADD
  CONSTRAINT users_pkey PRIMARY KEY (id)