---
name: add-dashboard-page
description: Add a new page (tab) to an existing dashboard category
argument-hint: [CategoryName] [PageType] [PageLabel]
---

Add a new page to the existing **$0** category. Page type: **$1**. Label: **$2**.

All changes go in `dashboard/new-dashboard/src/routes.ts` and optionally a new Vue component.

## Steps

### 1. Add route enum entry

In the `ROUTES` enum, add a new entry using the category's `ROUTE_PREFIX`:

```typescript
MyCategoryNewPage = `${ROUTE_PREFIX.MyCategory}/newPage`,
```

### 2. Add tab to the Product definition

Find the existing Product const for the category and add a `tab(url, label)` entry to its `tabs` array:

```typescript
tab(ROUTES.MyCategoryNewPage, "New Page"),
```

Use existing label constants when applicable: `TESTS_LABEL`, `DASHBOARD_LABEL`, `STARTUP_LABEL`, `PRODUCT_METRICS_LABEL`, `COMPARE_BUILDS_LABEL`, `COMPARE_BRANCHES_LABEL`, `COMPARE_MODES_LABEL`.

### 3. Add route handler

Add a new entry to the category's existing routes array using the route-builder helpers (defined near the top of `routes.ts`).

**For a Tests page** (uses the shared `PerformanceTests.vue` component via `perfTests(path, props, pageTitle)`):

```typescript
perfTests(
  ROUTES.MyCategoryNewPage,
  { dbName: "<dbName>", table: "<table>", initialMachine: MACHINES.HETZNER, withInstaller: false },
  "My Category Performance tests"
),
```

`props` is typed as `PerformanceTestsProps` (`dbName`, `table`, `initialMachine`, and optional `withInstaller`, `unit`, `releaseConfigurator`, `branch`, `withoutAccidents`, `machineGroupFilter`). The `satisfies TypedRouteRecord<PerformanceTestsProps>` annotation lives inside the helper.

**For a Dashboard page** (custom Vue component via `dashboard(path, component, pageTitle, props?)`):

```typescript
dashboard(ROUTES.MyCategoryNewPage, () => import("./components/myCategory/NewPageDashboard.vue"), "My Category New Page"),
```

**For a Startup dashboard** (`startupDashboard(path, { table, defaultProject? }, pageTitle)`):

```typescript
startupDashboard(ROUTES.MyCategoryStartup, { table: "<table>", defaultProject: "<project>" }, "My Category Startup dashboard"),
```

**For a Compare Builds page** (`compareBuilds(path, { dbName, table }, pageTitle?)`):

```typescript
compareBuilds(ROUTES.MyCategoryCompare, { dbName: "<dbName>", table: "<table>" }),
```

**For a Compare Branches page** (`compareBranches(path, { dbName, table, metricsNames? }, pageTitle?)`):

```typescript
compareBranches(ROUTES.MyCategoryCompareBranches, { dbName: "<dbName>", table: "<table>" }),
```

**For a Compare Modes page** (`compareModes(path, { dbName, table }, pageTitle?)`):

```typescript
compareModes(ROUTES.MyCategoryCompareModes, { dbName: "<dbName>", table: "<table>" }),
```

The `pageTitle` argument of the three compare helpers defaults to the standard label (`Compare Builds` / `Compare Branches` / `Compare Modes`), so pass it only to override.

### 4. Create Vue component (if needed)

Only needed for custom dashboard pages (not for Tests/Compare/Startup pages, which use shared components).

Create `dashboard/new-dashboard/src/components/<category>/NewPageDashboard.vue`:

```vue
<template>
  <DashboardPage
    db-name="<dbName>"
    table="<table>"
    persistent-id="<category>_<page>_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  />
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
</script>
```

## Conventions

- Route enum keys use PascalCase: `KotlinNotebooksTests`
- Route URL paths use camelCase: `/kotlinNotebooks/tests`
- Use standard route constants (`TEST_ROUTE`, `DEV_TEST_ROUTE`, `DASHBOARD_ROUTE`, `STARTUP_ROUTE`, `PRODUCT_METRICS_ROUTE`, `COMPARE_ROUTE`, `COMPARE_BRANCHES_ROUTE`, `COMPARE_MODES_ROUTE`) when the page type matches
- Build tabs with `tab(url, label)`; build routes with the `perfTests` / `dashboard` / `startupDashboard` / `compareBuilds` / `compareBranches` / `compareModes` helpers
- Reference machines via `MACHINES.HETZNER` / `MACHINES.AWS_LINUX` instead of raw strings
- Only drop to a raw object literal with `satisfies TypedRouteRecord<...>` for special cases the helpers don't cover (e.g. the `PerformanceUnitTests.vue` component, or a `:subproject?` path param)
