# go-healthcheck
The thing to check that your service is working

```
           .,----,.
        .:'        `:.
      .'              `.
     .'                `.
     :                  :
     `    .'`':'`'`/    '
      `.   \  |   /   ,'
        \   \ |  /   /
         `\_..,,.._/'
          {`'-,_`'-}
          {`'-,_`'-}
          {`'-,_`'-}
           `YXXXXY'
             ~^^~

```

### Getting started

```bash
go get github.com/Sagleft/go-healthcheck
```

First we create a handler on the top level of the service:

```go
var Healthchecker *gohealth.Handler = gohealth.NewHandler(gohealth.HandlerTask{})
```

By default, the listening will open at: `GET 127.0.0.1:8080/healthcheck`.
Returns the code 200 if all checks are passed, and 500 if there is an error and returns its text.

Then write checks in the necessary places in your service:

```go
Healthchecker.AddCheckpoint(gohealth.CheckpointData{
  Name: "LastError-Checker",
  CheckCallback: func() gohealth.Signal {
    if app.LastError != nil { // example
      return gohealth.SignalError("last error: " + app.LastError.Error())
    }
    return gohealth.SignalNormal()
  },
})
```

Or we can put the check in a separate method:

```go
Healthchecker.AddCheckpoint(gohealth.CheckpointData{
  Name: "LastError-Checker",
  CheckCallback: app.LastErrorCheck,
})
```
