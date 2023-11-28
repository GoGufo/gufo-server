# GUFO API Gateway

With Gufo you can create any API server you want. Just need to write a plugin or GRPC microservice with your features and connect it to Gufo.

## Generate GRPC connection files with proto
```docker
docker run -v $PWD:/defs namely/protoc-all -f microservice.proto -l go  #or ruby, csharp, etc
```
