# planet

```shell
$ cd mars && go run server.go
2025/03/31 09:11:10 connect to http://localhost:4001/ for GraphQL playground
```

```shell
$ cd venus && go run server.go
2025/03/31 09:11:17 connect to http://localhost:4002/ for GraphQL playground
```

```shell
$ cd earth && pnpm run start
🚀  Server ready at: http://localhost:4000/
```

```shell
$ curl -s http://localhost:4000/ --json '{"query":"query GetrPost { getPost(postId:\"1\") { id title body author { id name email } } }"}'
{"data":{"getPost":{"id":"1","title":"1: Hello World","body":"1: Hello World","author":{"id":"1","name":"John","email":"<EMAIL>"}}}}

mars> go run server.go
2025-03-31T09:14:40.116+0900	INFO	mars/server.go:82	access log	{"http_method": "POST", "path": "/query", "remote_addr": "[::1]:53831", "user_agent": "minipass-fetch/3.0.5 (+https://github.com/isaacs/minipass-fetch)", "request": "{\"query\":\"query GetrPost__mars__0{getPost(postId:1){id title body author{__typename id}}}\",\"variables\":{},\"operationName\":\"GetrPost__mars__0\"}", "response": "{\"data\":{\"getPost\":{\"id\":\"1\",\"title\":\"1: Hello World\",\"body\":\"1: Hello World\",\"author\":{\"__typename\":\"User\",\"id\":\"1\"}}}}", "response_duration_second": "0.001"}

venus> go run server.go
2025-03-31T09:14:40.120+0900	INFO	venus/server.go:82	access log	{"http_method": "POST", "path": "/query", "remote_addr": "[::1]:53832", "user_agent": "minipass-fetch/3.0.5 (+https://github.com/isaacs/minipass-fetch)", "request": "{\"query\":\"query GetrPost__venus__1($representations:[_Any!]!){_entities(representations:$representations){...on User{name email}}}\",\"variables\":{\"representations\":[{\"__typename\":\"User\",\"id\":\"1\"}]},\"operationName\":\"GetrPost__venus__1\"}", "response": "{\"data\":{\"_entities\":[{\"name\":\"John\",\"email\":\"\\u003cEMAIL\\u003e\"}]}}", "response_duration_second": "0.001"}
```
