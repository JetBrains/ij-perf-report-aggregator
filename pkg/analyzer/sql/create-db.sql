BEGIN;

create table report
(
    id               string not null primary key,
    machine          string not null,
    generated_time   int    not null,
    metrics_version  int    not null,

    product          string not null,
    build_c1         int    not null,
    build_c2         int    not null,
    build_c3         int    not null,

    duration_metrics string not null,
    instant_metrics  string not null,
    raw_report       string not null
);

create index machine_index on report (machine);
create index product_index on report (product);
create index build_major_index on report (build_c1);
create index build_minor_index on report (build_c2);

pragma user_version=3;

COMMIT;