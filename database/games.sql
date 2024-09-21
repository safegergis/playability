CREATE TABLE
  public.games (
    id serial NOT NULL,
    name character varying(255) NOT NULL,
    summary text NOT NULL,
    cover_art character varying(255) NOT NULL,
    platforms jsonb NOT NULL,
    closed_captions feature_support NOT NULL,
    color_blind feature_support NOT NULL,
    full_controller_support feature_support NOT NULL,
    controller_remapping feature_support NOT NULL
  );

ALTER TABLE
  public.games
ADD
  CONSTRAINT games2_pkey PRIMARY KEY (id)