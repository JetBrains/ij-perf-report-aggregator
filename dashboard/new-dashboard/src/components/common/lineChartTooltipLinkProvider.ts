import { provide } from "vue"
import { reportInfoProviderKey } from "../../shared/injectionKeys"

export function provideReportUrlProvider(isInstallerExists: boolean = true, isBuildNumberExists: boolean = false): void {
  const infoFields = ["machine", "tc_build_id", "project"]
  if (isInstallerExists) {
    infoFields.push("tc_installer_build_id", "build_c1", "build_c2", "build_c3")
  }
  if (isBuildNumberExists) {
    infoFields.push("build_number")
  }
  provide(reportInfoProviderKey, {
    infoFields,
  })
}
