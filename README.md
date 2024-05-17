# About
GUFO - General Universal FramewOrk. In Italian "Gufo" means Owl. This projects starts as RESTfull API Gateway for plugins. Now Gufo supports GRPC connections. So you can create GRPC/plugins microservices.

# GUFO API Gateway

With Gufo you can create any API server you want. Just need to write a plugin or GRPC microservice with your features and connect it to Gufo.

## Generate GRPC connection files with proto

go to /proto folder
```docker
docker run -v $PWD:/defs namely/protoc-all -f microservice.proto -o go/ -l go  #or ruby, csharp, etc
```

## Build Gufo

```docker
docker build --no-cache -t amyerp/gufo-api-gateway:latest -f Dockerfile .
```
