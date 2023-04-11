```

`service.duration` Array(Int32) CODEC (ZSTD(20)),
`measure.duration` Array(UInt32) CODEC(ZSTD(20)),

```

For arbitrary metrics that store as a nested array using Gorilla codec for `duration` doesn't help but only increases size. ZSTD as the only code works better. Confirmed by experimenting on IJ database.
But for `start` situation is different — Gorilla should be used as first codec before ZSTD.

Performance tips:
1. Use partitioning by months not years! To avoid scanning unnecessary parts.
2. As keys (order by) use columns that are used in WHERE conditions AND order they by cardinality - the lower goes first.