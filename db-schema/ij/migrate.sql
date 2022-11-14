/*

1. update table sql if needed (see actual sql in clickhouse metadata)
2. backup old database (the whole clickhouse directory).
3. use transform.go to re-analyze data or use clickhouse-client because other client read timeout maybe not enough (copy data is a long operation) to perform SQL operation.
*/

drop table report;
rename table report2 to report

-- insert into report2 select * from report;

/*
// migrate changes as new-line delimited to array
insert into installer2 select id, if(changes = '', emptyArrayString(), splitByChar('\n', changes)) from installer;
drop table installer;
rename table installer2 to installer;
 */


