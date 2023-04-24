DROP table if exists project_description;

CREATE TABLE project_description
(
  project     VARCHAR(255) NOT NULL,
  branch      VARCHAR(255) NOT NULL,
  url         VARCHAR(255),
  methodName  VARCHAR(255),
  description VARCHAR(255),
  PRIMARY KEY (project, branch)
);

DROP INDEX if exists idx_projectDescription_project;
CREATE INDEX idx_projectDescription_project ON project_description (project);
