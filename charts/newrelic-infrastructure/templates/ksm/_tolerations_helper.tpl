{{- /*
As this chart deploys what it should be three charts to maintain the transition to v3 as smooth as possible.
This means that this chart has 3 tolerations so a helper should be done per scraper.
*/ -}}
{{- define "nriKubernetes.ksm.tolerations" -}}
{{- if .Values.ksm.tolerations -}}
    {{- toYaml .Values.ksm.tolerations -}}
{{- else if include "common.tolerations" . -}}
    {{- include "common.tolerations" . -}}
{{- end -}}
{{- end -}}