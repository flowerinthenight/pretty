## Overview

`pretty` is a simple wrapper tool that prettifies a process's stdout/stderr's JSON output line by line. If output line is not a valid JSON, or does not contain a valid JSON, it will just print the line as is. It also works in "follow" mode, such as `tail -f ...` or `kubectl logs -f ...`.

It also works with https://github.com/wercker/stern (which I use a lot), that is, the JSON log has a text prefix (pod name in `stern`'s case). It will print the text prefix first, then append the prettified JSON.

Example: (normal JSON logs using logrus, redacted)

```bash
$ kubectl logs -f svc-699544fd4d-zzlcf
...
{"action":"describe-deployments","http_method":"GET","level":"info","service":"stack","time":"..."...}
{"context":"metrics-middleware","http_method": "GET","level":"info","service":"stack","time":"..."...}
{"action": "encode-response","http_method": "GET","level": "info","msg": "marshal response as is"...}
...
```

Example: (wrapped by pretty, redacted)

```bash
$ pretty -- kubectl logs -f svc-699544fd4d-zzlcf
...
2018/05/11 14:58:45 [stdout] {
  "action": "describe-deployments",
  "http_method": "GET",
  "level": "info",
  "msg": "describe={Type:compute.v1.instance Zone:asia-northeast1-a Name:...}",
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
  "msg": "marshal response as is",
  "request": "92b86ab1-7c8e-4f77-871a-7caf132b421e",
  "time": "2018-05-11T05:58:45Z"
}
...
```

## Installation

```bash
$ go get -u -v github.com/flowerinthenight/pretty
```

## Usage

You can prepend the command that you want to prettify with `pretty --`. The double-dash will ensure that any succeeding flags belong to the wrapped command, not to `pretty`.

```bash
$ pretty -- kubectl logs -f svc-699544fd4d-zzlcf
```
