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

Find the existing Product const for the category and add a new tab entry to its `tabs` array:

```typescript
{
  url: ROUTES.MyCategoryNewPage,
  label: "New Page",
},
```

Use existing label constants when applicable: `TESTS_LABEL`, `DASHBOARD_LABEL`, `STARTUP_LABEL`, `PRODUCT_METRICS_LABEL`, `COMPARE_BUILDS_LABEL`, `COMPARE_BRANCHES_LABEL`, `COMPARE_MODES_LABEL`.

### 3. Add route handler

Add a new entry to the category's existing routes array.

**For a Tests page** (uses the shared `PerformanceTests.vue` component):
```typescript
{
  path: ROUTES.MyCategoryNewPage,
  component: COMPONENTS.perfTests,
  props: {
    dbName: "<dbName>",
    table: "<table>",
    initialMachine: MACHINES.HETZNER,
    withInstaller: false,
  },
  meta: { pageTitle: "My Category Performance tests" },
} satisfies TypedRouteRecord<PerformanceTestsProps>,
```

**For a Dashboard page** (custom Vue component):
```typescript
{
  path: ROUTES.MyCategoryNewPage,
  component: () => import("./components/myCategory/NewPageDashboard.vue"),
  meta: { pageTitle: "My Category New Page" },
},
```

**For a Compare Builds page**:
```typescript
{
  path: ROUTES.MyCategoryCompare,
  component: COMPONENTS.compareBuilds,
  meta: { pageTitle: "My Category Compare Builds" },
},
```

**For a Compare Branches page**:
```typescript
{
  path: ROUTES.MyCategoryCompareBranches,
  component: COMPONENTS.compareBranches,
  meta: { pageTitle: "My Category Compare Branches" },
},
```

### 4. Create Vue component (if needed)

Only needed for custom dashboard pages (not for Tests/Compare pages which use shared components).

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
- Use standard route constants (`TEST_ROUTE`, `DEV_TEST_ROUTE`, `DASHBOARD_ROUTE`, `COMPARE_ROUTE`, `COMPARE_BRANCHES_ROUTE`, `COMPARE_MODES_ROUTE`) when the page type matches
- The `satisfies TypedRouteRecord<PerformanceTestsProps>` annotation is required on routes using `COMPONENTS.perfTests`