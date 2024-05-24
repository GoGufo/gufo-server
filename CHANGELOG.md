# Changelog

## 1.14.1
- Fixed bugs with Proto
- Other minor improvements

## 1.14.0
- Modify Proto for all internal requests. Add Arguments in Internal
- Add GRPC connection for Internal Requests
- Other small modifications

## 1.13.1
- Delete session functions from gufodao
- Modify Proto for masterservice ext

## 1.13.0
### New Features
- Huge modifications related to Microservice architectire. Move session check to session microservice.
- Master microservice implementation
- Config modifications. All microservices will be registered and hosted by Master microservice

## 1.12.4
### Improvemetns
- Code Refactoring
- Add PATCH and TRACE Methods
- Modify JSON Requests Functions to allow use freedom Authorisation tokens and Headers
- Add GRPC requests with PATCH and TRACE Methods
- Add Sign as argument in GRPC request. It necessary for access to internal microservices
- Create extantion to microservices types. Now it can be internal with no possible for connections outside of Gufo

## 1.12.3
### Improvemetns
- Modify GRPC internal request URL key in config file to avoid from conflicts
- Add JSONReq for make simple server-server requests


## 1.12.2
### New Features
- Add GET, POST, PUT, DELETE requests between microservices

## 1.12.1

### Bugs Fixes
- Issue with wrong path Request
- Issue with sessions
- other small fixes

### Improvemetns
- Fixed issue with file upload
- Remove Request struct

## 1.12.0

### Core Improvemetns
- Remove plugins support
- Change Error response structure
- Remove API v2 suport
- Proto modifications. Preparations for streaming connection

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
