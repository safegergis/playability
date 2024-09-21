CREATE TABLE
  public.reports (
    id serial NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    game_id integer NOT NULL,
    user_id integer NOT NULL,
    closed_captions feature_support NOT NULL,
    color_blind feature_support NOT NULL,
    full_controller_support feature_support NOT NULL,
    controller_remapping feature_support NOT NULL,
    score character varying(255) NOT NULL,
    report character varying(255) NULL
  );

ALTER TABLE
  public.reports
ADD
  CONSTRAINT reports_pkey PRIMARY KEY (id)