
build-api-gateway:
	docker build -t api_gateway:10.0.2 -f api_gateway/Dockerfile .

build-identity-service:
	docker build -t identity_service:1.0.4 -f identity_service/Dockerfile .

build-file_storage:
	docker build -t file_storage:1.0.5 -f file_storage/Dockerfile .

build-file_metadata:
	docker build -t file_metadata:1.0.3 -f file_metadata/Dockerfile .
