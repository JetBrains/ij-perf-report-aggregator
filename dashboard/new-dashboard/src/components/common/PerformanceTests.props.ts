import type { ReleaseType } from "../../configurators/ReleaseNightlyConfigurator"

// Props of PerformanceTests.vue. Kept in a plain .ts module so other modules (e.g. routes.ts) can
// import the type without resolving a .vue file — type-aware lint (tsgolint) can't see .vue exports.
export interface PerformanceTestsProps {
  dbName: string
  table: string
  initialMachine: string | string[] | null
  withInstaller?: boolean
  unit?: "ns" | "ms"
  releaseConfigurator?: ReleaseType
  branch?: string | null
  withoutAccidents?: boolean
  machineGroupFilter?: string | string[] | null
}
