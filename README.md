# redirsrv

a simple redirection server (rewritten in Go)

## How to use it?

1. Build and start the container

```
$ docker build -t robertgzr/redirsrv .
$ docker run robertgzr/redirsrv
DBUG| config initialized
INFO| token for API auth        bearer="<token>"
INFO| started listening         host="localhost" port=8080
```

2. Using the bearer token inspect the API (using [httpie](http://httpie.org/))

```
$ http GET :8080/adm bucket==redirs Authorization:Bearer:<token>
...

[
    "example"
]
```

result is that there is nothing

3. Add the route `:8080/test` that redirects to `https://example.com` via the API

```
$ echo "https://example.com" | http POST :8080/adm bucket==redirs key==test Authorization:Bearer:<token>
...

{
    "message": "success",
    "status": 201
}
```

4. Test redirection works

```
$ http -h :8080/test
HTTP/1.1 302 Found
...
```
