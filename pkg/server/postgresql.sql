DROP table if exists accidents;

CREATE TABLE accidents
(
  id            SERIAL PRIMARY KEY,
  date          DATE NOT NULL,
  affected_test VARCHAR(255) NOT NULL,
  reason        TEXT NOT NULL,
  branch        VARCHAR(255) NOT NULL,
  db_table      VARCHAR(50) NOT NULL,
  build_number  VARCHAR(20),
  kind          VARCHAR(50) NOT NULL default 'regression'
);

CREATE INDEX idx_accidents_table_branch ON accidents(db_table, branch);
CREATE INDEX idx_accidents_date ON accidents(date);
CREATE INDEX idx_accidents_affected_test ON accidents(affected_test);

INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-01-13 20:00:00', 'community/go-to-class/EditorImpl', 'IDEA-310587: 1ba879cb4472', 'master', 'perfint_idea', '231.5967');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-01-17 00:00:00', 'community/go-to-class/EditorImpl', 'IDEA-311498: 1e4a7cd1c77d', 'master', 'perfint_idea', '231.5362');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-01-25 23:00:00', 'intellij_commit/findUsages', 'memory allocation was added to async profiling: 14de0ae', 'master', 'perfintDev_idea', '');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-02-01 05:00:00', 'community/go-to-action/SharedIndex', 'cleaning workspace cache has been added and later removed', 'master', 'perfint_idea', '231.6544');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-02-11 02:00:00', 'community/go-to-class/EditorImpl', 'changed logic of collecting SE metrics: 83bd9b9bc8bf', 'master', 'perfint_idea', '231.7336');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-02-14 23:00:00', 'intellij_commit/findUsages', 'head commit of the tested project was updated: 7e2691a', 'master', 'perfintDev_idea', '');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-02-16 16:00:00', 'intellij_commit/localInspection', 'IDEA-313677: e2bba17dbe42', 'master', 'perfintDev_idea', '');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-03-04 07:00:00', 'mediawiki/inspection', 'JSHint became initialized', 'master', 'perfint_phpstorm', '232.1409');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-03-08 23:00:00', 'intellij_commit/completion/java_file', 'a38c21124132', 'master', 'perfintDev_idea', '277675176');
INSERT INTO accidents (date, affected_test, reason, branch, db_table, build_number) VALUES ('2023-03-13 07:00:00', 'kotlin_coroutines/highlight', 'file which is opened in the test was updated', 'master', 'perfint_idea', '232.1926');