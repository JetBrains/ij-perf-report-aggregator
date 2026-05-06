DROP TABLE IF EXISTS space_uploaded_artifacts;

CREATE TABLE space_uploaded_artifacts
(
  build_id       VARCHAR(20) PRIMARY KEY,
  uploaded_files TEXT[],
  success        BOOLEAN NOT NULL
);
