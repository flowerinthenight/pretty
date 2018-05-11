## Overview

Pretty is a simple wrapper tool that prettifies a process's stdout/stderr's JSON output line by line. If output line is not a valid JSON, it will just print the line as is.

```bash
# normal JSON logs sample (using logrus), redacted
$ kubectl logs -f svc-699544fd4d-zzlcf
...
{"action":"describe-deployments","http_method":"GET","level":"info","service":"stack","time":"2018-05-11T05:58:45Z"...}
{"context":"metrics-middleware","http_method": "GET","level":"info","service":"stack","time":"2018-05-11T05:58:45Z"...}
{"action": "encode-response","http_method": "GET","level": "info","msg": "marshal main response as is"...}
...

# wrapped by pretty (redacted)
$ pretty -- kubectl logs -f svc-699544fd4d-zzlcf
...
2018/05/11 14:58:45 [stdout] {
  "action": "describe-deployments",
  "http_method":"GET",
  "level":"info",
  "msg": "describe={Type:compute.v1.instance Zone:asia-northeast1-a Name:gcp-5837e6c9ef3bd-000000-vm-yrukryuk}",
  "request": "92b86ab1-7c8e-4f77-871a-7caf132b421e",
  "service": "stack",
  "time": "2018-05-11T05:58:45Z"
}
2018/05/11 14:58:45 [stdout] {
  "action": "describe-deployments",
  "http_method": "GET",
  "level": "info",
  "msg": "describe={Type:compute.v1.instance Zone:asia-northeast1-a Name:gcp-5837e6c9ef3bd-000000-vm-123123123}",
  "request": "92b86ab1-7c8e-4f77-871a-7caf132b421e",
  "service": "stack",
  "time": "2018-05-11T05:58:45Z"
}
2018/05/11 14:58:45 [stdout] {
  "action": "describe-deployments",
  "http_method": "GET",
  "level": "info",
  "msg": "describe={Type:compute.v1.instance Zone:asia-northeast1-a Name:gcp-5837e6c9ef3bd-000000-vm-asdas}",
  "request": "92b86ab1-7c8e-4f77-871a-7caf132b421e",
  "service": "stack",
  "time": "2018-05-11T05:58:45Z"
}
2018/05/11 14:58:45 [stdout] {
  "context": "metrics-middleware",
  "http_method": "GET",
  "level": "info",
  "msg": "fn=DescribeDeployment, duration=1.561291986s",
  "request": "92b86ab1-7c8e-4f77-871a-7caf132b421e",
  "service": "stack",
  "time": "2018-05-11T05:58:45Z"
}
2018/05/11 14:58:45 [stdout] {
  "action": "encode-response",
  "http_method": "GET",
  "level": "info",
  "msg": "marshal main response as is",
  "request": "92b86ab1-7c8e-4f77-871a-7caf132b421e",
  "time": "2018-05-11T05:58:45Z"
}
...
```
