## Bitbucket Stats Pet Project

Doesn't do much yet :)

To use this tool follow these steps:

1. Clone this repo to your `$GOPATH/src`
2. cd into `bitbucket` directory
3. run `go run *.go` to run program. `-h` for help
4. Three main functions are access by appending either `stats`, `update`, or `get`
5. Copy `https://***REMOVED***/rest/api/1.0` into a file called `.env.url`
6. Copy your bitbucket creds into a file called `.env.creds` in the form `<username>:<password>`
7. run `go run *.go update`
8. Now the data should be synced and you can run `go run *.go stats -h` for more info on how to get stats about languages in our bitbucke instance.
