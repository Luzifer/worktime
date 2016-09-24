==== Summary of {{ .day.DayID }} ====
Tags: {{ range .day.Tags }}{{ . }} {{ else }}n/a{{ end }}

{{ if .day.Times -}}
ID       Start     End       Tags{{ range .day.Times }}
{{ printf "%.7s" .ID }}  {{ .Start }}  {{ .End }}  {{ range .Tags}}{{ . }} {{ else }}n/a{{ end }}{{ end }}

{{ end -}}
Overtime: {{ .overtime }}
