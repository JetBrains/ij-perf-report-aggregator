DROP TABLE IF EXISTS space_uploaded_artifacts;

CREATE TABLE space_uploaded_artifacts
(
  build_id       VARCHAR(20)  NOT NULL,
  project        VARCHAR(255) NOT NULL,
  uploaded_files TEXT[],
  success        BOOLEAN      NOT NULL,
  PRIMARY KEY (build_id, project)
);
