import { ref } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { loadJson } from "../util/httpUtil"
import { toMutableArray, DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, ServerConfigurator } from "../dataQuery"
import { DebouncedTask, TaskHandle } from "../util/debounce"

// natural sort of alphanumerical strings
const collator = new Intl.Collator(undefined, {numeric: true, sensitivity: "base"})

export class MeasureConfigurator implements DataQueryConfigurator {
  public readonly data = ref<Array<string>>([])
  public readonly value = ref<Array<string>>([])

  private readonly debouncedLoad = new DebouncedTask(taskHandle => this.load(taskHandle))

  constructor(private readonly serverConfigurator: ServerConfigurator,
              persistentStateManager: PersistentStateManager,
              private readonly skipZeroValues: boolean = true) {
    persistentStateManager.add("metrics", this.value)
  }

  configure(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const values = toMutableArray(this.value.value)
    if (values.length === 0) {
      return false
    }

    // stable order of series (UI) and fields in query (caching)
    values.sort((a, b) => collator.compare(a, b))

    query.addField({
      name: "t",
      sql: "toUnixTimestamp(generated_time) * 1000"
    })

    configuration.dimensions.push({name: "time", type: "time"})

    for (let i = 0; i < values.length; i++) {
      const value = values[i]
      query.addField(value)
      if (this.skipZeroValues) {
        query.addFilter({field: value, operator: "!=", value: 0})
      }

      // remove _d (duration) or _i (instant) suffix
      const name = value.replace(/_[a-z]$/g, "")
      configuration.addSeries({
        name,
        type: "line",
        smooth: true,
        showSymbol: false,
        legendHoverLink: true,
        sampling: "lttb",
        encode: {
          // index if time
          x: 0,
          // +1 because time is the 0-dimension
          y: i + 1,
          tooltip: [i + 1],
        },
      }, {
        name,
        type: "int",
      })
    }
    if (query.order != null) {
      throw new Error("order must be configured only by MetricLoader")
    }
    query.order = ["t"]
    return true
  }

  scheduleLoad(immediately: boolean): void {
    this.debouncedLoad.execute(immediately)
  }

  load(taskHandle: TaskHandle): Promise<unknown> {
    const server = this.serverConfigurator.server.value
    if (server == null || server.length === 0) {
      return Promise.resolve()
    }

    return loadJson<Array<string>>(`${server}/api/v1/meta/measure?db=${this.serverConfigurator.databaseName}`, null, taskHandle, data => {
      this.data.value = data
    })
  }
}