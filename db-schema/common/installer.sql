create table installer
(
  `id`      UInt32 CODEC(Gorilla, ZSTD(20)),
  `changes` Array(FixedString(27)) CODEC(ZSTD(20))
)
  engine = MergeTree
    order by id
    settings old_parts_lifetime = 10
