import { provide } from "vue"
import { reportInfoProviderKey } from "../../shared/injectionKeys"

export function provideReportUrlProviderForStartup(): void {
  provide(reportInfoProviderKey, {
    infoFields: ["machine", "tc_build_id", "project", "tc_installer_build_id", "build_c1", "build_c2", "build_c3"],
  })
}

export function provideReportUrlProviderWithDataFetching(isInstallerExists: boolean = false, isBuildNumberExists: boolean = false): void {
  const infoFields = ["tc_build_id", "project"]
  if (isInstallerExists) {
    infoFields.push("build_c1", "build_c2", "build_c3")
  }
  if (isBuildNumberExists) {
    infoFields.push("build_number")
  }
  provide(reportInfoProviderKey, {
    infoFields,
  })
}
