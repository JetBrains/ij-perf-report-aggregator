create table collector_state
(
  build_type_id String CODEC (ZSTD(20)),
  last_time     DateTime CODEC (Gorilla, ZSTD(20))
)
  engine = ReplacingMergeTree(last_time)
    order by build_type_id
    settings old_parts_lifetime = 10

