create table accident
(
  id            INTEGER primary key,
  date          TEXT not null,
  affected_test TEXT not null,
  reason        TEXT not null,
  branch        TEXT not null,
  db_table      TEXT not null,
  build_number  TEXT
);
INSERT INTO accident (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-03-04 07:00:00', 'mediawiki/inspection', 'JSHint became initialized', 'master', 'phpstorm', '232.1409');