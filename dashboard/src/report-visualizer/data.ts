// Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
export interface ItemV0 extends CommonItem {
  readonly name: string
  readonly description?: string

  readonly start: number
  readonly end: number

  readonly duration: number

  readonly thread: string
}

export interface ItemV20 {
  readonly s: number
  readonly d: number
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
  readonly appComponents?: Array<ItemV20>
  readonly projectComponents?: Array<ItemV20>
  readonly moduleComponents?: Array<ItemV20>

  readonly appServices?: Array<ItemV20>
  readonly projectServices?: Array<ItemV20>
  readonly moduleServices?: Array<ItemV20>

  readonly serviceWaiting?: Array<ItemV20>
}

export interface InputData {
  readonly traceEvents: Array<TraceEvent>

  readonly version: string

  readonly stats: Stats
  readonly plugins?: Array<PluginStatItem>

  readonly icons?: Array<IconData>

  // time in ms
  readonly items: Array<ItemV0>

  // time in ms
  readonly prepareAppInitActivities: Array<ItemV20>

  readonly appExtensions?: Array<ItemV0>
  readonly projectExtensions?: Array<ItemV0>
  readonly moduleExtensions?: Array<ItemV0>

  readonly preloadActivities?: Array<ItemV0>
  readonly appOptionsTopHitProviders?: Array<ItemV0>
  readonly projectOptionsTopHitProviders?: Array<ItemV0>

  readonly projectPostStartupActivities?: Array<ItemV0>
}

export interface Stats {
  readonly plugin: number

  readonly component: StatItem
  readonly service: StatItem

  readonly loadedClasses?: { [key: string]: number }
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