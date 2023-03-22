{{/* vim: set filetype=mustache: */}}
{{/*
Create chart name and version as used by the chart label.
*/}}

{{- define "release-name" -}}
{{- if hasPrefix .Chart.Name .Release.Name -}}
{{- .Release.Name -}}
{{- else -}}
{{- .Release.Name -}}
{{- end -}}
{{- end -}}

{{- define "namespace" -}}
{{- if ne .Release.Namespace "default" -}}
{{- .Release.Namespace -}}
{{- else -}}
{{- "_default_namespace_is_forbidden_" -}}
{{- end -}}
{{- end -}}

{{- define "timestamp" -}}
{{- date "2006-01-02T15.04.05" .Release.Time -}}
{{- end -}}

{{- define "chart-version" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}
