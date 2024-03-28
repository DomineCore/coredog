{{- define "coredog.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "coredog.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "coredog.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "coredog.labels" -}}
helm.sh/chart: {{ include "coredog.chart" . }}
app.kubernetes.io/name: {{ include "coredog.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "coredog.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coredog.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
watcher Selector labels
*/}}
{{- define "coredog.watcher.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coredog.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.role: watcher
{{- end }}

{{/*
watcher labels
*/}}
{{- define "coredog.watcher.labels" -}}
helm.sh/chart: {{ include "coredog.chart" . }}
app.kubernetes.io/name: {{ include "coredog.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.role: watcher
{{- end }}

{{/*
Custom template function to convert a string to lowercase
*/}}
{{- define "coredog.tolower" -}}
{{- . | lower }}
{{- end }}
