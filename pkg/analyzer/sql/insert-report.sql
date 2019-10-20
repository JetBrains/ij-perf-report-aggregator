REPLACE INTO report (id, machine, product,
                     generated_time, build_time,
                     tc_build_id, tc_installer_build_id, tc_build_properties,
                     build_c1, build_c2, build_c3,
                     raw_report)
VALUES (?, ?, ?,
        ?, ?,
        ?, ?, ?,
        ?, ?, ?,
        ?)