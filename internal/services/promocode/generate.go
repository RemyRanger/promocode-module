package promocode

//go:generate oapi-codegen --old-config-style -o ports/handler_types.gen.go -package=ports -include-tags=Promocodes -generate types ../../../doc/promocode-api.yaml
//go:generate oapi-codegen --old-config-style -o ports/handler_api.gen.go -package=ports -include-tags=Promocodes -generate chi-server ../../../doc/promocode-api.yaml
