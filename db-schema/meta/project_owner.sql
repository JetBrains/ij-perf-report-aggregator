DROP TABLE IF EXISTS project_owner;

CREATE TABLE project_owner
(
  project    VARCHAR(255) NOT NULL,
  owner      VARCHAR(255) NOT NULL,
  db_name    VARCHAR(255) NOT NULL,
  table_name VARCHAR(255) NOT NULL,
  PRIMARY KEY (project)
);

DROP INDEX IF EXISTS idx_projectOwner_owner;
CREATE INDEX idx_projectOwner_owner ON project_owner (owner);
