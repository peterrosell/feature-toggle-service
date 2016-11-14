CREATE TABLE public.feature (
  id          TEXT    NOT NULL PRIMARY KEY,
  name        TEXT    NOT NULL UNIQUE,
  enabled     BOOLEAN NOT NULL,
  description TEXT    NOT NULL
);

CREATE TABLE public.property (
  name        TEXT NOT NULL PRIMARY KEY,
  description TEXT NOT NULL
);

CREATE TABLE public.toggle_rule (
  id        TEXT      NOT NULL,
  featureId TEXT      NOT NULL,
  property  TEXT      NOT NULL,
  value     TEXT      NOT NULL,
  created   TIMESTAMP NOT NULL,
  expires   TIMESTAMP,
  enabled   BOOLEAN   NOT NULL,
  PRIMARY KEY (id, property),
  CONSTRAINT fk_property
  FOREIGN KEY (property)
  REFERENCES public.property (name),
  CONSTRAINT fk_feature
  FOREIGN KEY (featureId)
  REFERENCES feature (id)
);
