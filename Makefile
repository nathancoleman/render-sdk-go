
gen: FORCE
	oapi-codegen --generate client,types --package apigen  -o gen/oapi.go --old-config-style oapi.yaml

FORCE:
