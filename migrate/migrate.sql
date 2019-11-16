/*

1. update table sql if needed (see actual sql in clickhouse metadata)
2. backup old database (the whole clickhouse directory).
3. use clickhouse-client because other client read timeout maybe not enough (copy data is a long operation)

  docker run -it --rm --volume=$HOME/ij-perf-db:/data:delegated --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server

docker run -it --rm --volume=$HOME/ij-perf-db:/data:delegated --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client:19.16.2.2 --host clickhouse-server

SELECT
    cast(product, 'UInt8') AS product,
    cast(machine, 'UInt8') AS machine,
    cast(branch, 'UInt8') AS branch,
    generated_time,
    build_time,
    raw_report,
    tc_build_id,
    tc_installer_build_id,
    tc_build_properties,
    build_c1,
    build_c2,
    build_c3
FROM report
ORDER BY
    product ASC,
    machine ASC,
    branch ASC,
    build_c1 ASC,
    build_c2 ASC,
    build_c3 ASC,
    build_time ASC,
    generated_time ASC
INTO OUTFILE '/data/f'
FORMAT Parquet
*/

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

