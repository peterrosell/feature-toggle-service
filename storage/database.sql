CREATE TABLE public.feature_properties (
  name TEXT NOT NULL,
  property TEXT NOT NULL,
  value TEXT NOT NULL,
  created TIMESTAMP,
  expires TIMESTAMP,
  enabled BOOLEAN,
  PRIMARY KEY (name, property)

);
CREATE UNIQUE INDEX feature_unique ON feature_properties_old USING BTREE (name, property);
