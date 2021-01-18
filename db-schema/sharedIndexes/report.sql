create table report
(
  `machine`               String CODEC (ZSTD(20)),
  `build_time`            DateTime CODEC (Delta(4), ZSTD(20)),
  `generated_time`        DateTime CODEC (Delta(4), ZSTD(20)),
  `project`               String CODEC (ZSTD(20)),
  `tc_build_id`           UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `tc_installer_build_id` UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `tc_build_properties`   String CODEC (ZSTD(20)),
  `branch`                Enum8('master' = 0) CODEC (ZSTD(20)),
  `raw_report`            String CODEC (ZSTD(20)),

  `build_c1`              UInt8 CODEC (DoubleDelta, ZSTD(20)),
  `build_c2`              UInt16 CODEC (DoubleDelta, ZSTD(20)),
  `build_c3`              UInt16 CODEC (DoubleDelta, ZSTD(20))
)
  engine = MergeTree
    partition by (toYYYYMM(generated_time))
    order by (machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time)
    settings old_parts_lifetime = 10;

-- View for metric descriptors.
--
-- <metric_name>         <metric_type>
-- indexing                  d
-- scanning                  d
-- numberOfIndexedFiles      c
-- etc
create or replace view metric_descriptors as
select tuple.1 as metric_name, tuple.2 as metric_type
from (
      select distinct arrayJoin(
                        arrayMap(
                          o -> (JSONExtractString(o, 'n'),
                                if(JSONHas(o, 'd'), 'd',
                                   if(JSONHas(o, 'c'), 'c',
                                      if(JSONHas(o, 'i'), 'i', 'd'))
                                  )
                            ),
                          JSONExtractArrayRaw(raw_report, 'metrics')
                          )
                        ) as tuple
      from report
       );

-- Metric values table
-- generated_time metric_name metric_type value
--
-- generated_time is used as the primary key for the 'report' table
create or replace view metric_values as
select generated_time,
       metric_name,
       metric_type,
       arrayExists(o -> JSONExtractString(o, 'n') = metric_name and JSONHas(o, metric_type), metric_array)   as exists,
       JSONExtractInt(arrayFirst(it -> JSONExtractString(it, 'n') = metric_name, metric_array), metric_type) as metric_value
from metric_descriptors
       cross join (select generated_time, JSONExtractArrayRaw(raw_report, 'metrics') as metric_array from report) as all_metrics
where exists;

-- Metrics table inheriting all fields of the 'report' table.
create or replace view metrics as
select *
from metric_values
       join report
            on metric_values.generated_time = report.generated_time;