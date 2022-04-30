AWS Account:
===

## Instances Describe :

| Instance ID | State | Type | Public IPv4 | Key Name |
| ---         | ---   | ---  | ---         | ---      |
{{- range .}}
| {{- .InstanceID -}} | {{- .InstanceState -}} | {{- .InstanceType -}} | {{- .PublicIPv4 -}} | {{- .KeyName -}} |
{{- end}}