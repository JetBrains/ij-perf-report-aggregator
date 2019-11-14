package server

type MachineInfo struct {
  GroupNames []string

  MachineToGroupName map[string]string
}

func GetMachineInfo() MachineInfo {
  // Mac mini Space Gray/3.0 GHz 6C/8GB/256GB; Model No. A1993; Part No. MRTT2RU/A; Serial No. C07XX9PFJYVX; Prod.12/2018, for code-sign (ADM-32069) -> ADM-35488
  const macMini = "macMini 2018"

  // Core i7-3770 16Gb, Intel SSD 535
  const win = "Windows: i7-3770, 16Gb, Intel SSD 535"

  // old RAM	RAM	RAM type	CPU	CPU CLOCK	MotherBoard	HDDs

  // 16384 Mb	16384 Mb	2xDDR3-12800 1600MHz 8Gb(8192Mb)	Core i7-3770	3400 Mhz	Intel DH77EB	240 Gb
  const linux = "Linux: i7-3770, 16Gb (12800 1600MHz), SSD"

  // 16384 Mb	16384 Mb	2xDDR3-10600 1333MHz 8Gb(8192Mb)	Core i7-3770	3400 Mhz	Intel DH77EB	240 Gb
  const linux2 = "Linux: i7-3770, 16Gb (10600 1333MHz), SSD"

  return MachineInfo{
    GroupNames: []string{macMini, linux, linux2, win},
    MachineToGroupName: map[string]string{
      "intellij-macos-hw-unit-1550": macMini,
      "intellij-macos-hw-unit-1551": macMini,

      "intellij-windows-hw-unit-498": win,
      "intellij-windows-hw-unit-499": win,

      "intellij-linux-hw-unit-449": linux,
      "intellij-linux-hw-unit-450": linux,
      "intellij-linux-hw-unit-463": linux2,
      "intellij-linux-hw-unit-484": linux,

      // error in info table - only 16GB ram and not 32
      "intellij-linux-hw-unit-493": linux,

      "intellij-linux-hw-unit-504": linux,
      "intellij-linux-hw-unit-531": linux,
      "intellij-linux-hw-unit-534": linux,
      "intellij-linux-hw-unit-556": linux,
      "intellij-linux-hw-unit-558": linux,
    },
  }
}
