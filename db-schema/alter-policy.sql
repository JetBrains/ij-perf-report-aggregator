ALTER TABLE ij.report MODIFY SETTING storage_policy='s3';
ALTER TABLE ij.installer MODIFY SETTING storage_policy='s3';
ALTER TABLE ij.collector_state MODIFY SETTING storage_policy='s3';

ALTER TABLE fleet.report MODIFY SETTING storage_policy='s3';
ALTER TABLE fleet.collector_state MODIFY SETTING storage_policy='s3';
ALTER TABLE fleet.measure MODIFY SETTING storage_policy='s3';
ALTER TABLE fleet.installer MODIFY SETTING storage_policy='s3';

ALTER TABLE perfint.collector_state MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.datagrip MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.goland MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.idea MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.ideaSharedIndices MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.installer MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.kotlin MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.phpstorm MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.phpstormWithPlugins MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.ruby MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.rust MODIFY SETTING storage_policy='s3';
ALTER TABLE perfint.scala MODIFY SETTING storage_policy='s3';

ALTER TABLE perfintDev.idea MODIFY SETTING storage_policy='s3';
ALTER TABLE perfintDev.collector_state MODIFY SETTING storage_policy='s3';