import { deepEqual } from "fast-equals"
import { distinctUntilChanged, Observable, switchMap,  withLatestFrom } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, toArray } from "../dataQuery"
import { loadDimension } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { createComponentState, updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

// todo what is it?
const macLarge = "mac large"

export class MachineConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly selected = shallowRef<Array<string>>([])
  readonly values = shallowRef<Array<GroupedDimensionValue>>([])

  private readonly observable: Observable<unknown>
  readonly state = createComponentState()
  private readonly groupNameToItem = new Map<string, GroupedDimensionValue>()

  constructor(serverConfigurator: ServerConfigurator, persistentStateManager: PersistentStateManager, filters: Array<FilterConfigurator> = [], readonly multiple: boolean = true) {
    const name = "machine"
    persistentStateManager.add(name, this.selected, it => toArray(it as never))
    const listObservable = createFilterObservable(serverConfigurator, filters)
      .pipe(
        switchMap(() => loadDimension(name, serverConfigurator, filters, this.state)),
        updateComponentState(this.state)
      )
    listObservable.subscribe(data => {
      if (data == null) {
        return
      }

      this.groupNameToItem.clear()
      this.values.value = this.groupValues(data)
    })

    // selected value may be a group name, so, we must re-execute query on machine list update
    this.observable = listObservable.pipe(
      withLatestFrom(refToObservable(this.selected, true)),
      distinctUntilChanged(deepEqual)
    )
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  private groupValues(values: Array<string>): Array<GroupedDimensionValue> {
    const grouped: Array<GroupedDimensionValue> = []
    for (const value of values) {
      let groupName = ""
      if (value.startsWith("intellij-linux-hw-blade-")) {
        groupName = "linux-blade"
      }
      else if (value.startsWith("intellij-windows-hw-blade-")) {
        groupName = "windows-blade"
      }
      else if (value.startsWith("intellij-windows-hw-munit-")) {
        groupName = "Windows Munich i7-3770, 32Gb"
      }
      else {
        if (value.startsWith("intellij-macos-unit-2200-large-")) {
          groupName = macLarge
        }
        else if (value.startsWith("intellij-linux-aws-m-i") || value.startsWith("intellij-linux-aws-3-lt") || value.startsWith("intellij-linux-aws-amd-2-lt")) {
          // noinspection SpellCheckingInspection
          groupName = "Linux EC2 m5d.xlarge or 5d.xlarge or m5ad.xlarge"
        }
        else if (value.startsWith("intellij-linux-performance-aws-i-")) {
          // https://aws.amazon.com/ec2/instance-types/c6i/
          // noinspection SpellCheckingInspection
          groupName = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
        }
        else if (value.startsWith("intellij-windows-performance-aws-i-")) {
          // https://aws.amazon.com/ec2/instance-types/c6id/
          // noinspection SpellCheckingInspection
          groupName = "Windows EC2 C6id.4xlarge (16 vCPU Xeon, 32 GB)"
        }
        else if (value.startsWith("intellij-linux-hw-munit-")) {
          groupName = "Linux Munich i7-3770, 32 Gb"
        }
        else if (value.startsWith("intellij-linux-hw-EXC")) {
          // Linux, i7-9700k, 2x16GiB DDR4-3200 RAM, NVME 512GB
          groupName = "Linux JB Expo AMS i7-3770, 32 Gb"
        }
        else if (value.startsWith("intellij-linux-hw-hetzner")){
          groupName = "linux-blade-hetzner"
        }
        else if (value.startsWith("intellij-windows-hw-hetzner")){
          groupName = "windows-blade-hetzner"
        }
        else if (value.startsWith("intellij-macos-munit-741-large")){
          //https://youtrack.jetbrains.com/issue/ADM-68723/Mac-agents-in-MYO-for-IntelliJ-and-JetBrains-Runtime
          groupName = "Mac Pro Intel Xeon E5-2697v2 (4x2.7GHz), 24 RAM"
        }
        else if (value.startsWith("intellij-linux-performance-huge-aws-i")){
          groupName = "Linux EC2 C6id.metal (128 CPU Xeon, 256 GB)"
        }

        if (groupName == null) {
          groupName = "Unknown"
          console.error(`Group is unknown for machine: ${value}`)
        }
      }

      let item = this.groupNameToItem.get(groupName)
      if (item == null) {
        item = {
          value: groupName,
          children: [],
        }
        grouped.push(item)
        this.groupNameToItem.set(groupName, item)
      }

      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      item.children!.push({value})
    }
    grouped.sort((a, b) => a.value.localeCompare(b.value))
    return grouped
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

    const values: Array<string> = []
    const filter: DataQueryFilter = {f: "machine", v: values}
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
        }
        else {
          // it's group
          if(groupItem.children != null) {
            if (groupItem.children.length > 50) {
              filter.v = prefix(groupItem.children.map(it => it.value)) + "%"
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
    if (value == null || value.length === 0) {
      return false
    }

    this.configureQueryAsFilter(value, query)
    return true
  }

  private configureQueryAsFilter(selected: Array<string>, query: DataQuery) {
    const values: Array<string> = []
    for (const value of selected) {
      const groupItem = this.groupNameToItem.get(value)
      if (groupItem == null) {
        values.push(value)
      }
      else {
        // it's group
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        for (const child of groupItem.children!) {
          values.push(child.value)
        }
      }
    }

    if (values.length > 0) {
      // stable order of fields in query (caching)
      values.sort()
      if (values.length > 50) {
        query.addFilter({f:"machine", v: prefix(values) + "%", o: "like"})
      } else {
        query.addFilter({f: "machine", v: values})
      }
    }
  }
}

export interface GroupedDimensionValue {
  value: string
  children?: Array<GroupedDimensionValue>
}

function prefix(words: Array<string>):string{
  if (!words[0] || words.length ==  1) return words[0] || ""
  let i = 0
  while(words[0][i] && words.every(w => w[i] === words[0][i]))
    i++
  return words[0].slice(0, Math.max(0, i))
}