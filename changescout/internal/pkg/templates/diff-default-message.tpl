🌐 {{.Name}} ({{.Mode}})
🔗 {{.URL}} | ⏱ {{.LastChecked}} {{"\n"}}
{{- range .Result.Changes }} ({{.Type }}): {{.Content}}{{"\n"}}{{- end }}