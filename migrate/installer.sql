create table installer
(
  `id`      UInt32 CODEC(DoubleDelta, ZSTD(19)),
  `changes` Array(FixedString(27)) CODEC(ZSTD(19))
)
  engine = MergeTree
    order by id
    settings old_parts_lifetime = 10, index_granularity = 8192
