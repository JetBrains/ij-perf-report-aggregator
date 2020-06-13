/*

1. update table sql if needed (see actual sql in clickhouse metadata)
2. backup old database (the whole clickhouse directory).
3. use clickhouse-client because other client read timeout maybe not enough (copy data is a long operation)

docker run -it --rm --volume=$HOME/ij-perf-db:/data:delegated --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client:20.3.4.10 --host clickhouse-server
*/

alter table report modify column machine Enum8('intellij-macos-hw-unit-1550' = 1, 'intellij-macos-hw-unit-1551' = 2,
  'intellij-windows-hw-unit-499' = 3, 'intellij-windows-hw-unit-498' = 4,
  'intellij-linux-hw-unit-558' = 5, 'intellij-linux-hw-unit-449' = 6, 'intellij-linux-hw-unit-450' = 7, 'intellij-linux-hw-unit-463' = 8, 'intellij-linux-hw-unit-504' = 9, 'intellij-linux-hw-unit-493' = 10, 'intellij-linux-hw-unit-556' = 11, 'intellij-linux-hw-unit-531' = 12, 'intellij-linux-hw-unit-484' = 13, 'intellij-linux-hw-unit-534' = 14,
  'Dead agent' = 15);

alter table report modify column product Enum8('IU' = 1, 'WS' = 2, 'PS' = 3, 'DB' = 4, 'GO' = 5, 'RM' = 6);

/* use clickhouse-client and not IDEA executor - to get progress and proper read-timeout  */
insert into report2 select * from report;

drop table report;
rename table report2 to report

/*
// migrate changes as new-line delimited to array
insert into installer2 select id, if(changes = '', emptyArrayString(), splitByChar('\n', changes)) from installer;
drop table installer;
rename table installer2 to installer;
 */

