CREATE TABLE grafana_views
(
  dashboard_uid VARCHAR(255) NOT NULL,
  user_login    VARCHAR(255) NOT NULL,
  day           DATE         NOT NULL DEFAULT CURRENT_DATE,
  views         INTEGER      NOT NULL DEFAULT 1,
  PRIMARY KEY (dashboard_uid, user_login, day)
);

CREATE INDEX idx_grafana_views_day
  ON grafana_views (day);
