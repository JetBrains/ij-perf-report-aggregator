CREATE TABLE missing_metrics
(
  id            SERIAL PRIMARY KEY,
  build_type    VARCHAR(255) NOT NULL,
  project       VARCHAR(255) NOT NULL,
  metric        VARCHAR(255) NOT NULL,
  missing_since TIMESTAMP    NOT NULL,
  CONSTRAINT unique_missing_metric UNIQUE (build_type, project, metric, missing_since)
);