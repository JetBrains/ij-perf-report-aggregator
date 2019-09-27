create table report
(
    id              string not null primary key,
    machine         string not null,
    generated_time  int    not null,
    metrics_version int    not null,
    metrics         string not null,
    raw_report      string not null
);

create index machine_index on report (machine);

pragma user_version=1