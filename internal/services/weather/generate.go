package weather

//go:generate oapi-codegen --old-config-style -o adapters/openweather_client/api_types.gen.go -include-tags=Weather -package=openweather_client -generate types ../../../doc/openweather-api.yaml
//go:generate oapi-codegen --old-config-style -o adapters/openweather_client/client.gen.go -include-tags=Weather -package=openweather_client -generate client ../../../doc/openweather-api.yaml
