# Project

IJ Perf collects performance reports from TeamCity, stores them in ClickHouse, detects degradations and visualizes the results.

- `cmd/`, `pkg/` — Go services and libraries (`backend`, `tc-collector`, `degradation-detector`, …).
- `dashboard/new-dashboard/` — Vue 3 + TypeScript dashboard (pnpm, Vite); pages live in `src/components`, routes in `src/routes.ts`.
- `db-schema/`, `deployment/` — ClickHouse schemas and deployment configs.
- `.agents/skills/` — step-by-step guides for recurring tasks (adding a dashboard category or page, updating ClickHouse).

# Formatting

Format TypeScript/Vue/JS/Markdown with oxfmt (config: `.oxfmtrc.json`); never run prettier — it is not a dependency, and `npx prettier` pulls a version that rewraps whole files.

```sh
node_modules/.bin/oxfmt ./dashboard/new-dashboard/src/routes.ts
```

Pass paths with a `./` prefix: bare relative paths can fail with "Expected at least one target file".

# Checks

Dashboard:

```sh
pnpm vue-tsc --noEmit -p dashboard/new-dashboard
pnpm lint
pnpm test
```

Go:

```sh
go test ./cmd/... ./pkg/...
golangci-lint run cmd/... pkg/...
```

`make lint` runs the full set (also `oxfmt --check`, `knip` and `deadcode`).
