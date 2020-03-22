create table collector_state
(
  build_type_id String CODEC (ZSTD(19)),
  last_time     DateTime CODEC (Delta(4), ZSTD(19))
)
  engine = ReplacingMergeTree(last_time)
    order by build_type_id
    settings old_parts_lifetime = 10

