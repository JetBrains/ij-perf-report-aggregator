# Formatting

Format TypeScript/Vue/JS/Markdown with oxfmt (config: `.oxfmtrc.json`); never run prettier — it is not a dependency, and `npx prettier` pulls a version that rewraps whole files.

```sh
node_modules/.bin/oxfmt ./dashboard/new-dashboard/src/routes.ts
```

Pass paths with a `./` prefix: bare relative paths can fail with "Expected at least one target file".

# Checks

```sh
pnpm vue-tsc --noEmit -p dashboard/new-dashboard
pnpm lint
pnpm test
```
