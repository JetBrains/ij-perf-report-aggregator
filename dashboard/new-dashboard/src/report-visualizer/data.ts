// Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
export interface ItemV20 {
  readonly s: number
  readonly d: number
  readonly e: number
  // own duration is specified only if differs from duration
  readonly od?: number
  readonly n: string
  readonly t: string
  readonly p?: string
}

export interface CommonItem {
  readonly name: string
}

export interface TraceEvent extends CommonItem {
  readonly name: string
  readonly ph: "i" | "X"
  // timestamp in microseconds
  readonly ts: number

  readonly tid: string

  readonly cat?: string
  readonly args?: TraceEventArgs
}

export interface TraceEventArgs {
  // our extension for services
  readonly ownDur: number
  readonly plugin?: string
}

export interface InputDataV20 extends InputData {
  readonly appComponents?: ItemV20[]
  readonly projectComponents?: ItemV20[]
  readonly moduleComponents?: ItemV20[]

  readonly appServices?: ItemV20[]
  readonly projectServices?: ItemV20[]
  readonly moduleServices?: ItemV20[]

  readonly serviceWaiting?: ItemV20[]
}

export interface InputData {
  readonly traceEvents: TraceEvent[]

  readonly version: string

  readonly stats: Stats
  readonly plugins?: PluginStatItem[]

  readonly icons?: IconData[]

  // time in ms
  readonly items: ItemV20[]

  // time in ms
  readonly prepareAppInitActivities: ItemV20[]

  readonly appExtensions?: ItemV20[]
  readonly projectExtensions?: ItemV20[]
  readonly moduleExtensions?: ItemV20[]

  readonly preloadActivities?: ItemV20[]
  readonly appOptionsTopHitProviders?: ItemV20[]
  readonly projectOptionsTopHitProviders?: ItemV20[]

  readonly projectPostStartupActivities?: ItemV20[]
}

export interface Stats {
  readonly plugin: number

  readonly component: StatItem
  readonly service: StatItem

  readonly loadedClasses?: Record<string, number>
}

export interface PluginStatItem {
  readonly id: string
  readonly classCount: number
  readonly classLoadingEdtTime: number
  readonly classLoadingBackgroundTime: number
}

export interface StatItem {
  readonly app: number
  readonly project: number
  readonly module: number
}

export interface IconData {
  readonly name: string
  readonly count: number

  readonly loading: number
  readonly decoding: number
}

export class UnitConverter {
  static MICROSECONDS = new UnitConverter(1000)
  static MILLISECONDS = new UnitConverter(1)

  private constructor(readonly factor: number) {}

  convert(value: number): number {
    return value / this.factor
  }
}
