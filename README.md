# About
GUFO - General Universal FramewOrk. In Italian "Gufo" means Owl. This projects starts as RESTfull API Gateway for plugins. Now Gufo supports GRPC connections. So you can create GRPC/plugins microservices.

# GUFO API Gateway

With Gufo you can create any API server you want. Just need to write a plugin or GRPC microservice with your features and connect it to Gufo.

Gufo can work as single and independent server in case if you do not need authorisation. For example you can connect to Gufo any gufo-compatible  microserivice for personal use. In other case, if you need to be youtrised you should use next microservices:

- Masterservice - it is main microservice. It is content host and ports of all microservices you connect to Gufo. With such microservice you can check in and check out any other microservices
- Auth - Autentication microservice. It check login and password and ask Session microservice to generate OAuth2.0 tokens
- Session - it holds access and refresh tokens. Every time Gufo asks Session microservice for check access token
- Rights - it holds OTP hashes and any other hashes. Also it holsd APP tokens for API request without authorisation
- Notifications - microservice who response for email and chat notifications. Auth microservice send OTP, password and confirmation emails with Notification microservices
- User -  content all inforamation about Users: name, Surname, birthdate etc.
- Admin - need for check in and check out microservices, grant rights to users, create and manipulate with users
- Reg - Sign Up new users. We decide to use Sign UP microservice independently in case if you need to collect and extra data from users with no need to rewrite Auth microservice

## Generate GRPC connection files with proto

go to /proto folder
```docker
docker run -v $PWD:/defs namely/protoc-all -f microservice.proto -o go/ -l go  #or ruby, csharp, etc
```

## Build Gufo

```docker
docker build --no-cache -t amyerp/gufo-api-gateway:latest -f Dockerfile .
```
