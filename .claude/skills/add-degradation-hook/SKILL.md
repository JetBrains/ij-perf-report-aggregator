---
name: add-degradation-hook
description: Send the new degradations of a dashboard to a webhook (Air Automation, a bot, a service) besides Slack
argument-hint: [HookName] [Db] [Table] [ProjectPrefix]
---

Add a degradation hook **$0** for the tests of **$1**.**$2** whose project starts with **$3**.

A hook gets each new degradation of the matching tests as a JSON POST, after the detector stores it as an accident.
The detector sends a degradation to a hook once. The code is in `pkg/degradation-detector/hooks.go`.

## Steps

### 1. The tests must have degradation settings

The hook sees only the degradations that a settings generator in `pkg/degradation-detector/setting/` finds. Check
that a generator has `PerformanceSettings` with this `Db`, `Table` and project.

### 2. Add the hook

Add an entry to `Hooks()` in `pkg/degradation-detector/setting/hooks.go`:

```go
{
	Name:          "my-hook",
	Db:            "perfintDev",
	Table:         "mytable",
	ProjectPrefix: "my/",
	URLEnv:        "MY_HOOK_URL",
	TokenEnv:      "MY_HOOK_TOKEN",
	AuthScheme:    "ApiKey", // "Bearer" when empty
},
```

Use a project prefix that only these tests have: an empty prefix is rejected by `TestHooksAreComplete`.

### 3. Add the URL and the token to the secret

The degradation-detector CronJob reads every key of the optional Kubernetes secret `degradation-detector-hooks` as an
environment variable. Ask the IJ Perf maintainers to add the keys `URLEnv` and `TokenEnv` to it. Until the keys exist,
the hook logs `hook is not configured` and sends nothing, so the entry can be merged first.

### 4. Check

```sh
go test ./pkg/degradation-detector/...
golangci-lint run ./pkg/degradation-detector/...
```

## The request

`POST <URL>` with `Authorization: <AuthScheme> <token>` and `Content-Type: application/json`. The body is a
`HookEvent`:

```json
{
  "eventId": "air/long-session/XL/open.XL.coverRemovedMs:1090164265",
  "source": "ij-perf",
  "kind": "degradation",
  "db": "perfintDev",
  "table": "air",
  "project": "air/long-session/XL",
  "metric": "open.XL.coverRemovedMs",
  "mode": "default",
  "branch": "master",
  "machine": "intellij-linux-performance-aws-%",
  "build": "1090164265",
  "date": "2026-10-14T03:00:00Z",
  "previousValue": 3000,
  "newValue": 4500,
  "changePercent": 50,
  "reason": "Degradation detected. Median changed by: 50.00%. ...",
  "chartUrl": "https://ij-perf.labs.jb.gg/air/testsDev?...",
  "reportUrl": "https://ij-perf.labs.jb.gg/degradations/report?..."
}
```

`eventId` is the same for every delivery of one degradation. A network error, a 5xx and a 429 are retried 3 times;
another answer that is not 2xx is not. A delivery that fails after its retries is logged and not sent again.
