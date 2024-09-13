DROP table if exists accidents;

CREATE TABLE accidents
(
  id            SERIAL PRIMARY KEY,
  date          DATE         NOT NULL,
  affected_test VARCHAR(255) NOT NULL,
  build_number  VARCHAR(20)  NOT NULL
    CONSTRAINT build_number_not_empty CHECK (build_number <> ''),
  reason        TEXT         NOT NULL,
  kind          VARCHAR(50)  NOT NULL default 'regression',
  externalId    VARCHAR(20)  NOT NULL default '',
  stacktrace    TEXT         NOT NULL,
  user_name     VARCHAR(100) NOT NULL default ''
);

CREATE INDEX idx_accidents_date ON accidents (date);
CREATE INDEX idx_accidents_affected_test ON accidents (affected_test);
CREATE UNIQUE INDEX unique_accidents_inferred ON accidents (date, affected_test, build_number, kind) WHERE kind <> 'EXCEPTION';