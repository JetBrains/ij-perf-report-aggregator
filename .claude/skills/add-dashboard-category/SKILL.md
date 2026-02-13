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
const enum ROUTE_PREFIX {
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

Before `export const PRODUCTS`, add the Product object:

```typescript
const MY_CATEGORY: Product = {
  url: ROUTE_PREFIX.MyCategory,
  label: "My Category",
  children: [
    {
      url: ROUTE_PREFIX.MyCategory,
      label: "",
      tabs: [
        {
          url: ROUTES.MyCategoryDashboard,
          label: DASHBOARD_LABEL,
        },
        {
          url: ROUTES.MyCategoryTests,
          label: TESTS_LABEL,
        },
      ],
    },
  ],
}
```

### 4. Add to PRODUCTS array

Insert the new product in alphabetical order within the `PRODUCTS` array.

### 5. Add route handlers

Before `export function getNewDashboardRoutes()`, add the routes array:

```typescript
const myCategoryRoutes = [
  {
    path: ROUTES.MyCategoryTests,
    component: COMPONENTS.perfTests,
    props: {
      dbName: "<the dbName argument>",
      table: "<the table argument>",
      initialMachine: MACHINES.HETZNER,
      withInstaller: false,
    },
    meta: { pageTitle: "My Category Performance tests" },
  } satisfies TypedRouteRecord<PerformanceTestsProps>,
  {
    path: ROUTES.MyCategoryDashboard,
    component: () => import("./components/myCategory/PerformanceDashboard.vue"),
    meta: { pageTitle: "My Category Dashboard" },
  },
]
```

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
- Use existing constants like `DASHBOARD_LABEL`, `TESTS_LABEL`, `TEST_ROUTE`, `DASHBOARD_ROUTE`
- The `satisfies TypedRouteRecord<PerformanceTestsProps>` annotation is required on test routes
