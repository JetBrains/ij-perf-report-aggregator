---
name: add-dashboard-category
description: Add a new top-level dashboard category with database mapping, tests page, and empty dashboard
argument-hint: [CategoryName] [dbName] [tableName]
---

Add a new top-level dashboard category called **$0** with database **$1** and table **$2**.

## Steps

All changes go in `dashboard/new-dashboard/src/routes.ts` plus one new Vue component.

### 1. Add route prefix

In the `ROUTE_PREFIX` enum, add a new entry. Use PascalCase for the key and a lowercase URL path derived from the category name.

Example:

```typescript
enum ROUTE_PREFIX {
  // ... existing entries
  MyCategory = "/myCategory",
}
```

### 2. Add route enum entries

In the `ROUTES` enum, add entries for Tests and Dashboard:

```typescript
enum ROUTES {
  // ... existing entries
  MyCategoryTests = `${ROUTE_PREFIX.MyCategory}/${TEST_ROUTE}`,
  MyCategoryDashboard = `${ROUTE_PREFIX.MyCategory}/${DASHBOARD_ROUTE}`,
}
```

### 3. Add Product definition

Before `export const PRODUCTS`, add the Product object. Tabs are built with the `tab(url, label)` helper:

```typescript
const MY_CATEGORY: Product = {
  url: ROUTE_PREFIX.MyCategory,
  label: "My Category",
  children: [
    {
      url: ROUTE_PREFIX.MyCategory,
      label: "",
      tabs: [tab(ROUTES.MyCategoryDashboard, DASHBOARD_LABEL), tab(ROUTES.MyCategoryTests, TESTS_LABEL)],
    },
  ],
}
```

### 4. Add to PRODUCTS array

Insert the new product in alphabetical order within the `PRODUCTS` array.

### 5. Add route handlers

Before `export function getNewDashboardRoutes()`, add the routes array. Use the route-builder helpers rather than raw object literals:

```typescript
const myCategoryRoutes = [
  perfTests(
    ROUTES.MyCategoryTests,
    { dbName: "<the dbName argument>", table: "<the table argument>", initialMachine: MACHINES.HETZNER, withInstaller: false },
    "My Category Performance tests"
  ),
  dashboard(ROUTES.MyCategoryDashboard, () => import("./components/myCategory/PerformanceDashboard.vue"), "My Category Dashboard"),
]
```

The route-builder helpers (defined near the top of `routes.ts`):

- `perfTests(path, props: PerformanceTestsProps, pageTitle)` — a Performance-tests page using the shared `PerformanceTests.vue`. The `satisfies TypedRouteRecord<PerformanceTestsProps>` annotation lives inside the helper, so you don't write it here.
- `dashboard(path, component, pageTitle, props?)` — a dashboard page; `component` is a lazy `() => import(...)`, `props` is an optional static-props record.
- `startupDashboard(path, { table, defaultProject? }, pageTitle)` — a startup-metrics dashboard.
- `compareBuilds(path, { dbName, table }, pageTitle?)`, `compareBranches(path, { dbName, table, metricsNames? }, pageTitle?)`, `compareModes(path, { dbName, table }, pageTitle?)` — comparison pages; `pageTitle` defaults to the standard label.

The `dbName` and `table` are extracted from the argument: if the argument is `perfintDev_kotlinNotebooks`, then `dbName` is `perfintDev` and `table` is `kotlinNotebooks`. If no underscore separator is present, use the full value as both `dbName` and `table`.

### 6. Register routes

Add `...myCategoryRoutes` to the `children` array inside `getNewDashboardRoutes()`, before the `ReportDegradations` route entry.

### 7. Create empty dashboard Vue component

Create a new file at `dashboard/new-dashboard/src/components/<componentDir>/PerformanceDashboard.vue`:

```vue
<template>
  <DashboardPage
    db-name="<dbName>"
    table="<table>"
    persistent-id="<table>_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  />
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
</script>
```

The component directory name should match the table name (camelCase).

## Important conventions

- Route prefix keys use PascalCase: `KotlinNotebooks`
- URL paths use camelCase: `/kotlinNotebooks`
- Product const names use UPPER_SNAKE_CASE: `KOTLIN_NOTEBOOKS`
- Route array variable names use camelCase: `kotlinNotebooksRoutes`
- PRODUCTS array is sorted alphabetically
- Build tabs with `tab(url, label)`; build routes with the `perfTests` / `dashboard` / `startupDashboard` / `compareBuilds` / `compareBranches` / `compareModes` helpers
- Use existing constants like `DASHBOARD_LABEL`, `TESTS_LABEL`, `TEST_ROUTE`, `DASHBOARD_ROUTE`, and `MACHINES.HETZNER` / `MACHINES.AWS_LINUX`
- Only drop to a raw object literal with `satisfies TypedRouteRecord<...>` for special cases the helpers don't cover (e.g. the `PerformanceUnitTests.vue` component, or a `:subproject?` path param)
