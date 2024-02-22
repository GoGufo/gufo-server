# Changelog

## 1.11.6

### Bug Fixing
- Fixed issue with wrong token initialisation

## 1.11.5

### Gufo functions

- Add ErrorReturn Function as general Error Handler
- Add CheckForSign and Gufo Sign. Gufo Sign is necessary for safety connection between GRPC microservices. By this sign, Microservice can understand that request was made from right GUFO instance

### proto file (GRPC Request)

- Add Sign data with GUFO Sign
- Add Requestor's IP address and UserAgent Data
