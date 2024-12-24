{{/*
Expand the name of the chart.
*/}}
{{- define "grey-svc-user.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "grey-svc-user.fullname" -}}
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
{{- define "grey-svc-user.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "grey-svc-user.labels" -}}
helm.sh/chart: {{ include "grey-svc-user.chart" . }}
{{ include "grey-svc-user.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "grey-svc-user.selectorLabels" -}}
app.kubernetes.io/name: {{ include "grey-svc-user.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "grey-svc-user.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "grey-svc-user.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generate a checksum for ConfigMap
*/}}
{{- define "grey-svc-user.configmapHash" -}}
{{- toYaml .Values.appConfig | sha256sum }}
{{- end }}

{{/*
Generate a checksum for External Secret
*/}}
{{- define "grey-svc-user.externalSecretHash" -}}
{{- toYaml .Values.externalSecrets.data | sha256sum }}
{{- end }}

{{/*
Combine both checksums for external config secret
*/}}
{{- define "grey-svc-user.externalConfigSecretChecksum" -}}
{{ include "grey-svc-user.configmapHash" . }}-{{ include "grey-svc-user.externalSecretHash" . }}
{{- end }}