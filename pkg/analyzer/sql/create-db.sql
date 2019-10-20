begin;

create table machine
(
  name id string not null primary key
);

create table report
(
  id                    string not null primary key,
  machine               int    not null,

  generated_time        int    not null,
  build_time            int default 0 not null,

  tc_build_id           int default 0 not null,
  tc_installer_build_id int default 0 not null,
  tc_build_properties   string,

  product               string not null,
  build_c1              int    not null,
  build_c2              int    not null,
  build_c3              int    not null,

  raw_report            string not null,

  foreign key (machine) references machine (ROWID) on delete restrict
);

create index generated_time on report (generated_time);
create index machine_index on report (machine);
create index product_index on report (product);
create index build_major_index on report (build_c1);
create index build_minor_index on report (build_c2);
create index build_patch_index on report (build_c3);

pragma user_version=6;

commit;