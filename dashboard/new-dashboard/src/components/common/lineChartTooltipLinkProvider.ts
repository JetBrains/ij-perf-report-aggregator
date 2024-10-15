import { provide } from "vue"
import { reportInfoProviderKey } from "../../shared/injectionKeys"
import { dbTypeStore } from "../../shared/dbTypes"

/**
 * If you change this method, you should adjust:
 * getInfoDataFrom in dashboard/new-dashboard/src/components/common/sideBar/InfoSidebarPerformance.ts
 * configureChart in dashboard/new-dashboard/src/configurators/MeasureConfigurator.ts:306
 */
export function provideReportUrlProvider(isInstallerExists: boolean = true, isBuildNumberExists: boolean = false): void {
  console.log("provideReportUrlProvider")
  const infoFields = ["machine", "tc_build_id", "project"]
  if (isInstallerExists) {
    infoFields.push("tc_installer_build_id", "build_c1", "build_c2", "build_c3")
  }
  if (isBuildNumberExists) {
    infoFields.push("build_number")
  }
  infoFields.push("branch")
  console.log(dbTypeStore().isModeSupported())
  if (dbTypeStore().isModeSupported()) {
    infoFields.push("mode")
  }
  provide(reportInfoProviderKey, {
    infoFields,
  })
}
