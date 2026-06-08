import { LlmAnalysisListItem } from "../common/llmAnalysis/LlmAnalysisClient"

export interface AnalysesFilterState {
  search: string
  states: string[]
  dateFrom: Date | null
  dateTo: Date | null
  commit: string
  users: string[]
  hasTicket: boolean
  hasFeedback: boolean
  hasCommits: boolean
}

export function emptyAnalysesFilterState(): AnalysesFilterState {
  return {
    search: "",
    states: [],
    dateFrom: null,
    dateTo: null,
    commit: "",
    users: [],
    hasTicket: false,
    hasFeedback: false,
    hasCommits: false,
  }
}

export function matchesCommit(commits: string[] | undefined, hash: string): boolean {
  const needle = hash.trim().toLowerCase()
  if (needle === "") return true
  if (commits == null) return false
  return commits.some((c) => c.toLowerCase().includes(needle))
}

export function userLabel(item: LlmAnalysisListItem): string {
  if (item.userName != null && item.userName !== "") return item.userName
  const email = item.userEmail
  if (email == null || email === "") return ""
  const at = email.indexOf("@")
  return at === -1 ? email : email.slice(0, at)
}

function matchesSearch(item: LlmAnalysisListItem, search: string): boolean {
  const needle = search.trim().toLowerCase()
  if (needle === "") return true
  const haystack = [item.project, item.metric, userLabel(item)].join(" ").toLowerCase()
  return haystack.includes(needle)
}

function isAfterOrEqualFrom(createdAt: number, from: Date | null): boolean {
  if (from == null) return true
  const start = new Date(from.getFullYear(), from.getMonth(), from.getDate(), 0, 0, 0, 0).getTime()
  return createdAt >= start
}

function isBeforeOrEqualTo(createdAt: number, to: Date | null): boolean {
  if (to == null) return true
  const end = new Date(to.getFullYear(), to.getMonth(), to.getDate(), 23, 59, 59, 999).getTime()
  return createdAt <= end
}

export function matchesAnalysisFilter(item: LlmAnalysisListItem, filter: AnalysesFilterState): boolean {
  if (!matchesSearch(item, filter.search)) return false
  if (filter.states.length > 0 && !filter.states.includes(item.state)) return false

  const createdAt = new Date(item.createdAt).getTime()
  if (Number.isFinite(createdAt)) {
    if (!isAfterOrEqualFrom(createdAt, filter.dateFrom)) return false
    if (!isBeforeOrEqualTo(createdAt, filter.dateTo)) return false
  }

  if (!matchesCommit(item.llmGuiltyCommits, filter.commit)) return false
  if (filter.users.length > 0 && !filter.users.includes(userLabel(item))) return false

  if (filter.hasTicket && (item.ytIssueId == null || item.ytIssueId === "")) return false
  if (filter.hasFeedback && item.feedbackCount <= 0) return false
  return !(filter.hasCommits && (item.llmGuiltyCommits == null || item.llmGuiltyCommits.length === 0))
}

export function filterAnalyses(items: LlmAnalysisListItem[], filter: AnalysesFilterState): LlmAnalysisListItem[] {
  return items.filter((item) => matchesAnalysisFilter(item, filter))
}

export function distinctUsers(items: LlmAnalysisListItem[]): string[] {
  const set = new Set<string>()
  for (const item of items) {
    const label = userLabel(item)
    if (label !== "") set.add(label)
  }
  return [...set].toSorted((a, b) => a.localeCompare(b))
}
