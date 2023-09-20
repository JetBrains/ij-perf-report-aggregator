import { Observable } from "rxjs"
import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { useSettingsStore } from "../components/settings/settingsStore"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export class DetectChangesConfigurator implements DataQueryConfigurator, FilterConfigurator {
  private settingsStore = useSettingsStore()
  readonly value = ref(this.settingsStore.detectChanges)

  constructor() {
    watch(
      () => this.settingsStore.detectChanges,
      (newValue) => {
        this.value.value = newValue
      }
    )
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.value)
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _configuration: DataQueryExecutorConfiguration | null): boolean {
    return true
  }
}

const extractValuesFromMatrix = (matrix: (string | number)[][], index: number): (string | number)[] => {
  return matrix.map((row) => row[index])
}

export function applyFilter(seriesData: (string | number)[][]): (string | number)[][] {
  const changePointIndexes = getChangePointIndexes(seriesData[1] as number[], 1).map((value) => value + 1)
  return changePointIndexes.map((value) => extractValuesFromMatrix(seriesData, value))
}

export function getChangePointIndexes(data: number[], minDistance: number = 1): number[] {
  const n = data.length

  if (n <= 2) return []
  if (minDistance < 1 || minDistance > n) {
    throw new Error(`minDistance (${minDistance}) should be in range from 1 to data.length`)
  }

  const penalty = 3 * Math.log(n)
  const k = Math.min(n, Math.ceil(4 * Math.log(n)))

  const partialSums = getPartialSums(data, k)
  const cost = (tau1: number, tau2: number): number => getSegmentCost(partialSums, tau1, tau2, k, n)

  const bestCost: number[] = new Array(n + 1).fill(0) as number[]
  bestCost[0] = -penalty
  for (let currentTau = minDistance; currentTau < 2 * minDistance; currentTau++) {
    bestCost[currentTau] = cost(0, currentTau)
  }

  const previousChangePointIndex = new Array(n + 1).fill(0) as number[]
  let previousTaus: number[] = [0, minDistance]
  let costForPreviousTau: number[] = []

  for (let currentTau = 2 * minDistance; currentTau < n + 1; currentTau++) {
    costForPreviousTau = previousTaus.map((previousTau) => bestCost[previousTau] + cost(previousTau, currentTau) + penalty)

    const bestPreviousTauIndex = whichMin(costForPreviousTau)
    bestCost[currentTau] = costForPreviousTau[bestPreviousTauIndex]
    previousChangePointIndex[currentTau] = previousTaus[bestPreviousTauIndex]

    const currentBestCost = bestCost[currentTau]
    previousTaus = previousTaus.filter((_, i) => costForPreviousTau[i] < currentBestCost + penalty)
    previousTaus.push(currentTau - (minDistance - 1))
  }

  const changePointIndexes: number[] = []
  let currentIndex = previousChangePointIndex[n]
  while (currentIndex !== 0) {
    changePointIndexes.push(currentIndex - 1)
    currentIndex = previousChangePointIndex[currentIndex]
  }
  return changePointIndexes.reverse()
}

function getPartialSums(data: number[], k: number): number[][] {
  const n = data.length
  const partialSums: number[][] = Array.from({ length: k }, () => new Array(n + 1).fill(0) as number[])
  const sortedData = [...data].sort((a, b) => a - b)

  for (let i = 0; i < k; i++) {
    const z = -1 + (2 * i + 1) / k
    const p = 1 / (1 + Math.pow(2 * n - 1, -z))
    const t = sortedData[Math.floor((n - 1) * p)]

    for (let tau = 1; tau <= n; tau++) {
      partialSums[i][tau] = partialSums[i][tau - 1]
      if (data[tau - 1] < t) partialSums[i][tau] += 2
      if (data[tau - 1] === t) partialSums[i][tau] += 1
    }
  }
  return partialSums
}

function getSegmentCost(partialSums: number[][], tau1: number, tau2: number, k: number, n: number): number {
  let sum = 0
  for (let i = 0; i < k; i++) {
    const actualSum = partialSums[i][tau2] - partialSums[i][tau1]

    if (actualSum !== 0 && actualSum !== (tau2 - tau1) * 2) {
      const fit = (actualSum * 0.5) / (tau2 - tau1)
      const lnp = (tau2 - tau1) * (fit * Math.log(fit) + (1 - fit) * Math.log(1 - fit))
      sum += lnp
    }
  }
  const c = -Math.log(2 * n - 1)
  return ((2 * c) / k) * sum
}

function whichMin(values: number[]): number {
  return values.indexOf(Math.min(...values))
}
