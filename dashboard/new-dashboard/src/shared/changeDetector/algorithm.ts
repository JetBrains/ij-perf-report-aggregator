import { hodgesLehmannEstimator, median, pooledShamos } from "./statistic"

export enum ChangePointClassification {
  DEGRADATION = "Degradation",
  OPTIMIZATION = "Optimization",
  NO_CHANGE = "No Change",
}

export function detectChanges(seriesData: (string | number)[][]): Map<string, ChangePointClassification> {
  const dataset = seriesData[1] as number[] | undefined
  const changePointIndexes = getChangePointIndexes(dataset, 5)
  const classifications = classifyChangePoint(changePointIndexes, dataset)
  const resultMap = new Map<string, ChangePointClassification>()

  for (const [index, value] of changePointIndexes.entries()) {
    const extractedValues = extractValuesFromMatrix(seriesData, value)
    resultMap.set(JSON.stringify(extractedValues), classifications[index])
  }
  return resultMap
}

const whichMin = (values: number[]): number => {
  return values.indexOf(Math.min(...values))
}

const getSegmentCost = (partialSums: number[][], tau1: number, tau2: number, k: number, n: number): number => {
  let sum = 0
  for (let i = 0; i < k; i++) {
    const actualSum = partialSums[i][tau2] - partialSums[i][tau1]

    if (actualSum !== 0 && actualSum !== (tau2 - tau1) * 2) {
      const fit = (actualSum * 0.5) / (tau2 - tau1)
      const lnp = (tau2 - tau1) * (fit * Math.log(fit) + (1 - fit) * Math.log1p(-fit))
      sum += lnp
    }
  }
  return ((2 * -Math.log(2 * n - 1)) / k) * sum
}

export const classifyChangePoint = (changePointIndexes: number[], dataset: number[] | undefined) => {
  if (dataset == undefined) return []
  const classifications: ChangePointClassification[] = []

  for (let i = 0; i < changePointIndexes.length; i++) {
    // If it's the first change point, take data from the beginning, otherwise from the previous change point.
    const startBefore = i === 0 ? 0 : changePointIndexes[i - 1]
    const endBefore = changePointIndexes[i]

    const startAfter = changePointIndexes[i]
    // If it's the last change point, take data till the end, otherwise till the next change point.
    const endAfter = i === changePointIndexes.length - 1 ? dataset.length : changePointIndexes[i + 1]

    const segmentBefore = dataset.slice(startBefore, endBefore)
    const segmentAfter = dataset.slice(startAfter, endAfter)

    const medianBefore = median(segmentBefore)
    const medianAfter = median(segmentAfter)

    const percentageDifference = Math.abs(((medianAfter - medianBefore) / medianBefore) * 100)

    const hlValue = hodgesLehmannEstimator(segmentBefore, segmentAfter)
    const shamos = pooledShamos(segmentBefore, segmentAfter)
    const effectSize = hlValue / shamos
    let classification

    if (
      (medianBefore < 2000 && percentageDifference < 5) ||
      (medianBefore >= 2000 && medianBefore < 10000 && percentageDifference < 2) ||
      (medianBefore >= 10000 && percentageDifference < 1)
    ) {
      classification = ChangePointClassification.NO_CHANGE
    } else if (Math.abs(medianBefore - medianAfter) < 10) {
      classification = ChangePointClassification.NO_CHANGE
    } else if (Math.abs(effectSize) < 2) {
      classification = ChangePointClassification.NO_CHANGE
    } else {
      classification = hlValue > 0 ? ChangePointClassification.DEGRADATION : ChangePointClassification.OPTIMIZATION
    }

    classifications.push(classification)
  }
  return classifications
}

const getPartialSums = (data: number[], k: number): number[][] => {
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
const extractValuesFromMatrix = (matrix: (string | number)[][], index: number): (string | number)[] => {
  return matrix.map((row) => row[index])
}

export function getChangePointIndexes(data: number[] | undefined, minDistance: number = 1): number[] {
  if (data == undefined) return []
  const n = data.length

  if (n <= 2) return []
  minDistance = Math.min(minDistance, n / 3)

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
    changePointIndexes.push(currentIndex)
    currentIndex = previousChangePointIndex[currentIndex]
  }
  return changePointIndexes.reverse()
}
