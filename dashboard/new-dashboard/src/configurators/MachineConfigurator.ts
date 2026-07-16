import { deepEqual } from "fast-equals"
import { combineLatest, distinctUntilChanged, Observable, shareReplay, switchMap } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, ServerConfigurator, toArray } from "../components/common/dataQuery"
import { loadDimension } from "./DimensionConfigurator"
import { createComponentState, updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"
// Single source of truth for machine hardware-class grouping, shared with the Go backend.
import machineGroupsJson from "@machineGroups"

// todo what is it?
const macLarge = "mac large"

export class MachineConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly selected = shallowRef<string[]>([])
  readonly values = shallowRef<GroupedDimensionValue[]>([])

  private readonly observable: Observable<unknown>
  readonly state = createComponentState()
  private readonly groupNameToItem = new Map<string, GroupedDimensionValue>()

  private static readonly valueToGroup: Record<string, string> = getValueToGroup()

  readonly filters = shallowRef<FilterConfigurator[]>([])

  constructor(
    serverConfigurator: ServerConfigurator,
    persistentStateManager?: PersistentStateManager,
    initialFilters: FilterConfigurator[] = [],
    readonly multiple: boolean = true,
    predefinedMachines?: string[]
  ) {
    const name = "machine"
    persistentStateManager?.add(name, this.selected, (it) => toArray(it as never))
    if (predefinedMachines) {
      this.selected.value = predefinedMachines
    }

    this.filters.value = initialFilters
    const filterObservable = refToObservable(this.filters).pipe(
      switchMap((currentFilters) => {
        return createFilterObservable(serverConfigurator, currentFilters)
      }),
      shareReplay(1)
    )

    const listObservable = filterObservable.pipe(
      switchMap(() => {
        return loadDimension(name, serverConfigurator, this.filters.value, this.state)
      }),
      updateComponentState(this.state),
      shareReplay(1)
    )

    listObservable.subscribe((data) => {
      if (data == null) {
        return
      }

      this.groupNameToItem.clear()
      this.values.value = this.groupValues(data)
    })

    // selected value may be a group name, so, we must re-execute query on machine list update
    this.observable = combineLatest([refToObservable(this.selected, true), listObservable]).pipe(distinctUntilChanged(deepEqual))
    // init groupNameToItem - if actual machine list is not yet loaded, but there is stored value for filter, use it to draw chart
    this.groupValues(Object.keys(MachineConfigurator.valueToGroup))
  }

  updateFilters(newFilters: FilterConfigurator[]) {
    this.filters.value = newFilters
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  private groupValues(values: string[]): GroupedDimensionValue[] {
    const grouped: GroupedDimensionValue[] = []
    for (const value of values) {
      let groupName: string | null = getMachineGroupName(value)
      const machineFromMap = MachineConfigurator.valueToGroup[value] as string | null
      if (groupName == "Unknown" && machineFromMap != null) {
        groupName = machineFromMap
      }

      let item = this.groupNameToItem.get(groupName)
      if (item == null) {
        item = {
          value: groupName,
          children: [],
          icon: this.getIcons(groupName),
        }
        grouped.push(item)
        this.groupNameToItem.set(groupName, item)
      }

      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      item.children!.push({ value })
    }
    grouped.sort((a, b) => a.value.localeCompare(b.value))
    return grouped
  }

  private getIcons(groupName: string): string {
    if (groupName.toLowerCase().startsWith("linux")) {
      return "pi icon-linux"
    } else if (groupName.toLowerCase().startsWith("mac")) {
      return "pi pi-apple"
    }
    return "pi pi-microsoft"
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const selected = this.selected.value
    if (selected.length === 0) {
      console.debug("machine is not configured")
      return false
    }

    if (!this.multiple) {
      this.configureQueryAsFilter(selected, query)
      return true
    }

    const groupNameToItem = this.groupNameToItem

    const values: string[] = []
    const filter: DataQueryFilter = { f: "machine", v: values }
    query.addFilter(filter)
    configuration.queryProducers.push({
      size(): number {
        return selected.length
      },

      mutate(index: number): void {
        const value = selected[index]
        const groupItem = groupNameToItem.get(value)
        values.length = 0
        if (groupItem == null) {
          values.push(value)
        } else {
          // it's group
          if (groupItem.children != null) {
            if (groupItem.children.length > 1) {
              filter.v = prefix(groupItem.children.map((it) => it.value)) + "%"
              filter.o = "like"
              return
            }
            for (const child of groupItem.children) {
              values.push(child.value)
            }
          }
          values.sort()
        }
      },

      getSeriesName(index: number): string {
        return selected.length > 1 ? selected[index] : ""
      },

      getMeasureName(index: number): string {
        return selected[index]
      },
    })
    return true
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value.length === 0) {
      return false
    }

    this.configureQueryAsFilter(value, query)
    return true
  }

  getMergedValue(): string {
    const values: string[] = this.getSelectedValues(this.selected.value)
    return prefix(values) + "%"
  }

  private configureQueryAsFilter(selected: string[], query: DataQuery) {
    const values = this.getSelectedValues(selected)

    if (values.length > 0) {
      // stable order of fields in query (caching)
      values.sort()
      if (values.length > 50) {
        query.addFilter({ f: "machine", v: prefix(values) + "%", o: "like" })
      } else {
        query.addFilter({ f: "machine", v: values })
      }
    }
  }

  private getSelectedValues(selected: string[]) {
    const values: string[] = []
    for (const value of selected) {
      const groupItem = this.groupNameToItem.get(value)
      if (groupItem == null) {
        values.push(value)
      } else {
        // it's group
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        for (const child of groupItem.children!) {
          values.push(child.value)
        }
      }
    }
    return values
  }
}

export interface GroupedDimensionValue {
  value: string
  children?: GroupedDimensionValue[]
  icon?: string
}

function prefix(words: string[]): string {
  if (!words[0] || words.length == 1) return words[0] || ""
  let i = 0
  while (words[0][i] && words.every((w) => w[i] === words[0][i])) i++
  return words[0].slice(0, Math.max(0, i))
}

function getValueToGroup() {
  // Mac mini Space Gray/3.0 GHz 6C/8GB/256GB
  const macMini = "macMini 2018"
  const macMiniIntel = "macMini Intel 3.2, 16GB"
  // Mac Mini M1 Chip with 8‑Core CPU und 8‑Core GPU, SSD 256Gb, RAM 16Gb
  const macMiniM1 = "macMini M1 2020"
  const macMiniM1_16 = "macMini M1, 16 Gb"
  // Mac Mini M2 Pro with 12-Core CPU, RAM 32 GB
  const macMiniM2 = "macMini M2 Pro, 32 GB"

  // Core i7-3770 16Gb, Intel SSD 535
  const win = "Windows Space i7-3770, 16Gb"

  // old RAM	RAM	RAM type	CPU	CPU CLOCK	MotherBoard	HDDs

  // 16384 Mb	16384 Mb	2xDDR3-12800 1600MHz 8Gb(8192Mb)	Core i7-3770	3400 Mhz	Intel DH77EB	240 Gb
  const linux = "Linux Space i7-3770, 16Gb"

  const blade = "linux-blade"

  return {
    "intellij-macos-hw-unit-1550": macMini,
    "intellij-macos-hw-unit-1551": macMini,
    "intellij-macos-hw-unit-1772": macMini,
    "intellij-macos-hw-unit-1773": macMini,

    "intellij-macos-hw-munit-716": macMiniIntel,
    "intellij-macos-hw-munit-721": macMiniIntel,
    "intellij-macos-hw-munit-722": macMiniIntel,
    "intellij-macos-hw-munit-723": macMiniIntel,
    "intellij-macos-hw-munit-724": macMiniIntel,

    "intellij-macos-hw-munit-608": macMiniM1_16,
    "intellij-macos-hw-munit-689": macMiniM1_16,
    "intellij-macos-hw-munit-690": macMiniM1_16,
    "intellij-macos-hw-munit-691": macMiniM1_16,
    "intellij-macos-hw-munit-692": macMiniM1_16,
    "intellij-macos-hw-munit-693": macMiniM1_16,
    "intellij-macos-hw-munit-694": macMiniM1_16,
    "intellij-macos-hw-munit-695": macMiniM1_16,
    "intellij-macos-hw-munit-696": macMiniM1_16,
    "intellij-macos-hw-munit-697": macMiniM1_16,
    "intellij-macos-hw-munit-698": macMiniM1_16,

    "intellij-macos-hw-unit-2204": macMiniM1,
    "intellij-macos-hw-unit-2205": macMiniM1,
    "intellij-macos-hw-unit-2206": macMiniM1,
    "intellij-macos-hw-unit-2207": macMiniM1,

    "intellij-macos-unit-2200-large-10298": macLarge,

    "intellij-macos-docker-hw-de-unit-1500": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1501": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1502": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1503": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1574": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1575": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1576": macMiniM2,
    "intellij-macos-docker-hw-de-unit-1577": macMiniM2,

    "intellij-windows-hw-unit-498": win,
    "intellij-windows-hw-unit-499": win,
    "intellij-windows-hw-unit-449": win,
    "intellij-windows-hw-unit-463": win,
    "intellij-windows-hw-unit-493": win,
    "intellij-windows-hw-unit-504": win,

    "intellij-linux-hw-unit-449": linux,
    "intellij-linux-hw-unit-499": linux,
    "intellij-linux-hw-unit-450": linux,
    "intellij-linux-hw-unit-484": linux,

    // error in info table - only 16GB ram and not 32
    "intellij-linux-hw-unit-493": linux,

    "intellij-linux-hw-unit-504": linux,
    "intellij-linux-hw-unit-531": linux,
    "intellij-linux-hw-unit-534": linux,
    "intellij-linux-hw-unit-556": linux,
    "intellij-linux-hw-unit-558": linux,

    "intellij-linux-hw-blade-023": blade,
    "intellij-linux-hw-blade-024": blade,
    "intellij-linux-hw-blade-025": blade,
    "intellij-linux-hw-blade-026": blade,
    "intellij-linux-hw-blade-027": blade,
    "intellij-linux-hw-blade-028": blade,
    "intellij-linux-hw-blade-029": blade,
    "intellij-linux-hw-blade-030": blade,
    "intellij-linux-hw-blade-031": blade,
    "intellij-linux-hw-blade-032": blade,
    "intellij-linux-hw-blade-033": blade,
    "intellij-linux-hw-blade-034": blade,
    "intellij-linux-hw-blade-035": blade,
    "intellij-linux-hw-blade-036": blade,
    "intellij-linux-hw-blade-037": blade,
    "intellij-linux-hw-blade-038": blade,
    "intellij-linux-hw-blade-039": blade,
    "intellij-linux-hw-blade-040": blade,
    "intellij-linux-hw-blade-041": blade,
    "intellij-linux-hw-blade-042": blade,
    "intellij-linux-hw-blade-043": blade,
    "intellij-linux-hw-blade-044": blade,
    "intellij-linux-hw-blade-045": blade,
    "intellij-linux-hw-blade-046": blade,
    "intellij-linux-hw-blade-047": blade,
    "intellij-linux-hw-blade-048": blade,
    "intellij-linux-hw-blade-049": blade,
    "intellij-linux-test-hw-blade-002": blade,

    "intellij-linux-hw-de-unit-1705": "Linux Munich i7-13700, 64 Gb",
    "intellij-linux-hw-de-unit-1716": "Linux Munich i7-13700, 64 Gb",
    "intellij-linux-hw-de-unit-1715": "Linux Munich i7-13700, 64 Gb",
    "intellij-windows-hw-de-unit-1702": "Windows Munich i7-13700, 64 Gb",
    "intellij-windows-hw-de-unit-1703": "Windows Munich i7-13700, 64 Gb",
    "intellij-windows-hw-de-unit-1704": "Windows Munich i7-13700, 64 Gb",
  }
}

// Hardware-class rules shared with the backend — the single source of truth is
// pkg/machine/groups.json (see pkg/machine/machine.go). First matching rule wins.
interface MachineGroupRule {
  group: string
  prefixes?: string[]
  regex?: string
}
const machineGroupRules: { group: string; prefixes?: string[]; regex?: RegExp }[] = (machineGroupsJson as MachineGroupRule[]).map((rule) => ({
  group: rule.group,
  prefixes: rule.prefixes,
  regex: rule.regex == null ? undefined : new RegExp(rule.regex),
}))

export function getMachineGroupName(machine: string): string {
  for (const rule of machineGroupRules) {
    if (rule.regex != null) {
      if (rule.regex.test(machine)) {
        return rule.group
      }
    } else if (rule.prefixes?.some((prefix) => machine.startsWith(prefix))) {
      return rule.group
    }
  }
  return "Unknown"
}
