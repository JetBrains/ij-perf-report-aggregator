create view metrics as
select machine,
       build_time,
       generated_time,
       project,
       tc_build_id,
       tc_installer_build_id,
       tc_build_properties,
       branch,
       build_c1,
       build_c2,
       build_c3,
       JSONExtractInt(arrayFirst(it -> JSONExtractString(it, 'n') = 'scanning', JSONExtractArrayRaw(raw_report, 'metrics')), 'd') as scanning,
       JSONExtractInt(arrayFirst(it -> JSONExtractString(it, 'n') = 'indexing', JSONExtractArrayRaw(raw_report, 'metrics')), 'd') as indexing
from report