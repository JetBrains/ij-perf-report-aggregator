DROP TABLE IF EXISTS analyses;

CREATE TABLE analyses
(
  id                    SERIAL PRIMARY KEY,
  created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  project               VARCHAR(255) NOT NULL,
  metric                VARCHAR(255) NOT NULL,
  current_build_id      VARCHAR(20)  NOT NULL,
  prev_build_id         VARCHAR(20)  NOT NULL,
  current_value         VARCHAR(50),
  previous_value        VARCHAR(50),
  user_name             VARCHAR(100),
  first_commit_revision VARCHAR(40),
  last_commit_revision  VARCHAR(40),
  test_method_name      VARCHAR(255),
  run_build_id          VARCHAR(20),
  yt_issue_id           VARCHAR(20),
  state                 VARCHAR(50)  NOT NULL DEFAULT 'in_progress',
  llm_guilty_commits    VARCHAR(40)[],
  llm_comment           VARCHAR(8000),
  total_cost_usd        NUMERIC(10, 4)
);

CREATE INDEX idx_analysis_lookup
  ON analyses (project, metric, current_build_id);