DROP TABLE IF EXISTS llm_run_accidents;
DROP TABLE IF EXISTS llm_run_yt_issues;
DROP TABLE IF EXISTS llm_analysis_runs;

CREATE TABLE llm_analysis_runs
(
-- fields that inserted on creation and not subject to change:
  id                    SERIAL PRIMARY KEY,
  created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  date                  DATE         NOT NULL,
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
-- fields that inserted on update and subject to change:
  run_build_id          VARCHAR(20),
  -- valid values: not_started | queued | in_progress | success | failed | cancelled
  state                 VARCHAR(50)  NOT NULL DEFAULT 'not_started',
-- fields that inserted on update and not a subject to change:
  llm_guilty_commits    VARCHAR(40)[],
  llm_comment           TEXT,
-- fields that inserted on update and subject to change:
  user_rate             BOOLEAN,
  user_comment          TEXT
);

CREATE INDEX idx_llm_analysis_runs_lookup
  ON llm_analysis_runs (date, project, metric, current_build_id, prev_build_id);

-- Junction table: LLM runs <-> accidents (many-to-many)
CREATE TABLE llm_run_accidents
(
  llm_run_id  INTEGER NOT NULL REFERENCES llm_analysis_runs (id) ON DELETE CASCADE,
  accident_id INTEGER NOT NULL REFERENCES accidents (id) ON DELETE CASCADE,
  PRIMARY KEY (llm_run_id, accident_id)
);

-- Junction table: LLM runs <-> YouTrack issues (many-to-many)
CREATE TABLE llm_run_yt_issues
(
  llm_run_id  INTEGER     NOT NULL REFERENCES llm_analysis_runs (id) ON DELETE CASCADE,
  yt_issue_id VARCHAR(20) NOT NULL,
  PRIMARY KEY (llm_run_id, yt_issue_id)
);