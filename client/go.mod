module client

go 1.15

require (
	google.golang.org/grpc v1.35.0
	iaas-api-server v0.0.0-00010101000000-000000000000
)

replace iaas-api-server => ../
