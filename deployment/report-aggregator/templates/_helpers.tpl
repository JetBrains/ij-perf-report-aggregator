{{- define "common.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{ include "common.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "common.selectorLabels" -}}
app.kubernetes.io/name: report-aggregator
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "s3Env" -}}
- name: S3_BUCKET
  valueFrom:
    secretKeyRef:
      name: ij-perf-data-minio-bucket
      key: AWS_ACCESS_KEY_ID
{{- end }}

{{- define "clickhouseVolumeMounts" -}}
volumeMounts:
  - name: clickhouse-data
    mountPath: /var/lib/clickhouse
  - name: s3
    mountPath: /etc/s3
    readOnly: true
{{- end }}
