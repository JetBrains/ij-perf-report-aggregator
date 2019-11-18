create table installer2
(
  `name` String CODEC (DoubleDelta, ZSTD(19))
)
  engine = MergeTree
    order by name
    settings old_parts_lifetime = 10, index_granularity = 8192;

create view m_mac_mini as
select array('intellij-macos-hw-unit-1550', 'intellij-macos-hw-unit-1551')

