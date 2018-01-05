# redirsrv

A simple redirection server built on [rocket.rs](https://rocket.rs)

## Usage

Write a `linkfile.json` it should look something like this:

```
{
    [
        { "short": "gh", "to": "https://github.com/robertgzr" }
    ]
}
```

Deploy using Docker:

```
$ make docker
$ docker run -d --name redirsrv -v "$(pwd)/linkfile.json:/etc/redirsrv/linkfile.json" -p 8080:80 robertgzr/redirsrv
```

Find the API token in the log and you can get started...

```
$ http GET :8080/adm Authorization:Bearer:<api_token>
```
