# API Structure

##1. POST /
link for all POST API requests

##1.1. Signup
### Request data
```
{
"module":"reg",
"param": "signup",
"args": {
          "user":"user",
          "pass": "xxx",
          "email":"user@example.com"
         },
"timestamp":111111,
"lang":"english"
}
```
### Response
```
{
    "success": 1,
    "data": {
        "response": "100200"
    },
    "timestamp": 1598211881,
    "lang": "eng"
}
```

##1.2. Sign In

### Request data
```
{
"module":"auth",
"param": "signin",
"args": {
          "user":"user", // or "user@example.com"
          "pass": "xxx",
         },
"timestamp":111111,
"lang":"english"
}
```

### Response
```
{
    "success": 1,
    "data": {
        "email": "user@example.com",
        "email_confirmed": true,
        "isadmin": false,
        "session_expired": 1909252348,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVC....",
        "uid": "lthuc5do",
        "username": "user"
    },
    "timestamp": 1598212348,
    "lang": "eng"
}
```

##1.3 Resend Confirmation email

### Request data


```
{
"module":"reg",
"param": "Resesendconfemail",
"timestamp":111111,
"lang":"english"
}
```

### Response

```

    "success": 1,
    "data": {
        "isadmin": 0,
        "session_expired": 1909315640,
        "uid": "lthuc5do"
    },
    "timestamp": 1598275640,
    "lang": "eng"
}
```



2. GET /info
Initial Information about Gufo server.

### Request data
no data

### Response
```
{
    "success": 1,
    "data": {
        "registration": true,
        "version": "v0.1.0"
    },
    "timestamp": 1597507704,
    "lang": "eng"
}
```

3. GET /confirmemail
This link send by email for confirm it
