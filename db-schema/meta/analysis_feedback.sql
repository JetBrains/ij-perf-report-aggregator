DROP TABLE IF EXISTS analysis_feedback;

CREATE TABLE analysis_feedback
(
  id          SERIAL PRIMARY KEY,
  analysis_id INTEGER      NOT NULL REFERENCES analyses (id) ON DELETE CASCADE,
  rate        SMALLINT     NOT NULL CHECK (rate BETWEEN 1 AND 5),
  feedback    TEXT,
  user_email  VARCHAR(255) NOT NULL,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  UNIQUE (analysis_id, user_email)
);

CREATE INDEX idx_analysis_feedback_analysis_id
  ON analysis_feedback (analysis_id);