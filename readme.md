
## Create a new config instance
```go
config := goshen.NewSshConfig("~/.ssh/config")
```

## Matching a host
Once you've created an instance of an SSH config, you can search for hosts using the `MatchHost` method.

### sample config
```
Host server-1
  HostName server-1.example.com
  User user
  Port 22
  IdentityFile ~/.ssh/example_rsa
  LogLevel FATAL

Host server-2*
  HostName server-2.example.com
  User user
  Port 22
  IdentityFile ~/.ssh/example_rsa
  LogLevel FATAL
```

```go
// return the config for server-2 with searching for server-2-web
hostEntry := config.MatchHost("server-2-web")

fmt.Println(hostEntry.HostName) // server-2.example.com
```

## Saving the config
Calling save on the config instance will write the configuration to disk. If the `outputPath` argument is an empty string, the path given
at the time `goshen.NewSshConfig` was called is used.

```go
outputPath := ""
config.Save(outputPath)
```