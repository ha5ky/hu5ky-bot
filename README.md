# hu5ky-bot

## build

```shell
go build -ldflags "-X 'github.com/ka5ky/hu5ky-bot/pkg/config.GitCommit=`git logger --pretty=format:%H -1`'"
```