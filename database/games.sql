CREATE TABLE
  public.games (
    id serial NOT NULL,
    name character varying(255) NOT NULL,
    summary text NOT NULL,
    cover_art character varying(255) NOT NULL,
    platforms jsonb NOT NULL,
    closed_captions character varying(255) NOT NULL,
    color_blind character varying(255) NOT NULL,
    full_controller_support character varying(255) NOT NULL,
    controller_remapping character varying(255) NOT NULL
  );

ALTER TABLE
  public.games
ADD
  CONSTRAINT games2_pkey PRIMARY KEY (id)