REPLACE INTO report (id, machine, product,
                     generated_time, tc_build_id,
                     build_c1, build_c2, build_c3,
                     metrics_version, duration_metrics, instant_metrics,
                     raw_report)
VALUES (?, ?, ?,
        ?, ?,
        ?, ?, ?,
        ?, ?, ?,
        ?)