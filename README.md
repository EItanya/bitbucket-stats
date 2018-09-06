## Bitbucket Stats Pet Project

Doesn't do much yet :)

To use this tool follow these steps:

1. Clone this repo to your `$GOPATH/src`
2. cd into `bitbucket` directory
3. run `go run *.go` to run program. `-h` for help
4. Three main functions are access by appending either `stats`, `update`, or `get`
5. Setup your config in config.json of the form:

```json
{
  "url": "https://<bitbucket.instance.url>/rest/api/1.0",
  "username": "<username>",
  "password": "<password>"
}
```

6. Download docker and the basic redis docker conatiner
7. Set the Access Url of the docker container to `localhost:6379`
8. run `go run *.go update`
9. Now the data should be synced and you can run `go run *.go stats -h` for more info on how to get stats about languages in our bitbucke instance.
