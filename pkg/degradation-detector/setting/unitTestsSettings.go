package setting

import (
	"errors"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"time"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector/alerts"
)

const (
	unitTestsDb    = "perfUnitTests"
	unitTestsTable = "report"
	unitTestMetric = "attempt.mean.ms"
)

// The unit performance tests run on the Blade agents. Some tests, such as the Go LSP tests,
// run only on the Hetzner agents, so the detector keeps the Hetzner series for them.
const (
	bladeUnitTestsMachine   = "ij-linux-x64-perf-hw-blade-%"
	hetznerUnitTestsMachine = "intellij-linux-%-hetzner-%"
)

// unitTestsOnBladeSince is the switch of the unit performance tests to the Blade agents.
// The Blade series contain the changes that the Hetzner series already reported, so the
// detector does not report a Blade change before this time. The detector does not post a
// change older than 28 days, so this has no effect from 07.12.2026.
var unitTestsOnBladeSince = time.Date(2026, time.November, 9, 0, 0, 0, 0, time.UTC)

// unitTestsMachineSeries is the tests that the detector analyzes on one machine pattern.
type unitTestsMachineSeries struct {
	machine             string
	tests               []string
	ignoreChangesBefore time.Time
}

// splitUnitTestsByMachine gives each test one series: the Blade series, or the Hetzner series
// when the test does not run on the Blade agents.
func splitUnitTestsByMachine(bladeTests, hetznerTests []string) []unitTestsMachineSeries {
	var hetznerOnly []string
	for _, test := range hetznerTests {
		if !slices.Contains(bladeTests, test) {
			hetznerOnly = append(hetznerOnly, test)
		}
	}
	return []unitTestsMachineSeries{
		{machine: bladeUnitTestsMachine, tests: bladeTests, ignoreChangesBefore: unitTestsOnBladeSince},
		{machine: hetznerUnitTestsMachine, tests: hetznerOnly},
	}
}

// generateUnitTestSettings returns the settings of unit-test-degradation: every performance unit test is analyzed,
// and the subscribers of the notification decide who gets its messages.
func generateUnitTestSettings(notification alerts.Notification, backend Backend) []detector.Settings {
	base := detector.PerformanceSettings{
		Db: unitTestsDb, Table: unitTestsTable, Branch: "master", Metric: unitTestMetric,
	}
	bladeTests, err := backend.Tests(base.Db, base.Table, base.Branch, bladeUnitTestsMachine)
	if err != nil {
		slog.Error("error while getting tests", "error", err, "machine", bladeUnitTestsMachine)
		return nil
	}
	hetznerTests, err := backend.Tests(base.Db, base.Table, base.Branch, hetznerUnitTestsMachine)
	if err != nil {
		slog.Error("error while getting tests", "error", err, "machine", hetznerUnitTestsMachine)
		return nil
	}

	// Without the owners a test still reaches the subscribers that select it by name, or the fallback channel.
	owners, err := backend.ProjectOwners(base.Db, base.Table)
	if err != nil {
		slog.Error("error while fetching project owners, the unit tests are routed by name only", "error", err)
	}
	// Without the groups a subscriber gets the tests of an owner only when it claims the exact name of the owner.
	groups, err := backend.OwnerGroups()
	if err != nil {
		slog.Error("error while fetching code-owner groups, a claim does not cover subgroups and old names", "error", err)
	}
	router, err := newUnitTestRouter(notification, groups)
	if err != nil {
		slog.Error("invalid spec in the alerts config", "notification", unitTestDegradation, "error", err)
		return nil
	}

	var result []detector.Settings
	for _, series := range splitUnitTestsByMachine(bladeTests, hetznerTests) {
		base.Machine = series.machine
		for _, test := range series.tests {
			base.Project = test
			for _, route := range router.resolve(test, owners[test]) {
				metrics := append([]string{unitTestMetric}, additionalMetrics(route.spec.AdditionalMetrics, base, backend)...)
				for _, metric := range metrics {
					settings := base
					settings.Metric = metric
					settings.SlackSettings = route.slack
					settings.AnalysisSettings = route.spec.analysis(unitTestDegradation, metric)
					settings.IgnoreChangesBefore = series.ignoreChangesBefore
					result = append(result, settings)
				}
			}
		}
	}
	return result
}

// CheckUnitTestOwners returns an error when the owners of the unit tests or the code-owner groups are not available.
// The unit tests are analyzed anyway, but their messages can go to the fallback channel instead of the owner.
func CheckUnitTestOwners(backend Backend) error {
	_, ownersErr := backend.ProjectOwners(unitTestsDb, unitTestsTable)
	_, groupsErr := backend.OwnerGroups()
	return errors.Join(ownersErr, groupsErr)
}

// additionalMetrics returns the metrics that the spec adds for the test of the settings.
func additionalMetrics(byTestPrefix map[string][]string, test detector.PerformanceSettings, backend Backend) []string {
	var result []string
	for prefix, patterns := range byTestPrefix {
		if !strings.HasPrefix(test.Project, prefix) {
			continue
		}
		for _, pattern := range patterns {
			names, err := backend.MetricNames(test, like(pattern))
			if err != nil {
				slog.Error("error while fetching metric names by pattern", "error", err, "test", test.Project, "pattern", pattern)
				continue
			}
			result = append(result, names...)
		}
	}
	slices.Sort(result)
	return result
}

// unitTestRoute is one receiver of the messages about a test.
type unitTestRoute struct {
	slack detector.SlackSettings
	spec  spec
}

type unitTestRouter struct {
	subscribers []unitTestSubscriber
	fallback    unitTestRoute
	// canonical is the current name of a group by its current or previous name.
	canonical map[string]string
	parent    map[string]string
}

type unitTestSubscriber struct {
	alerts.Subscriber

	route  unitTestRoute
	blocks []unitTestRoute
}

func newUnitTestRouter(notification alerts.Notification, groups []detector.CodeOwnerGroup) (*unitTestRouter, error) {
	link := productLink(unitTestsTable)
	fallbackSpec, err := parseSpec(notification.Spec)
	if err != nil {
		return nil, err
	}
	router := &unitTestRouter{
		fallback:  unitTestRoute{slack: slack(alerts.Target{Channel: notification.Fallback}, false, link), spec: fallbackSpec},
		canonical: map[string]string{},
		parent:    map[string]string{},
	}
	for _, group := range groups {
		router.parent[group.Name] = group.Parent
		router.canonical[group.Name] = group.Name
		for _, alias := range group.Aliases {
			router.canonical[alias] = group.Name
		}
	}
	for _, name := range notification.SubscriberNames() {
		subscriber := unitTestSubscriber{Subscriber: notification.Subscribers[name]}
		s, err := parseSpec(subscriber.Spec)
		if err != nil {
			return nil, err
		}
		subscriber.route = unitTestRoute{slack: slack(subscriber.Target, !subscriber.Enabled, link), spec: s}
		for _, block := range subscriber.Match {
			blockSpec, err := parseSpec(block.Spec)
			if err != nil {
				return nil, err
			}
			subscriber.blocks = append(subscriber.blocks, unitTestRoute{slack: slack(block.Target, !subscriber.Enabled, link), spec: blockSpec})
		}
		router.subscribers = append(router.subscribers, subscriber)
	}
	return router, nil
}

// resolve returns who gets the messages about the test, as the README of the alerts repository tells: each
// subscriber whose block selects the test by name or who claims the owner of the test, and the fallback channel when
// there is no such subscriber. A disabled subscriber gets no message, but its test does not go to the fallback.
//
// Subscribers with the same channel and the same spec give one route: the channel gets one message.
func (r *unitTestRouter) resolve(test string, owner string) []unitTestRoute {
	var routes []unitTestRoute
	add := func(route unitTestRoute) {
		for _, existing := range routes {
			if reflect.DeepEqual(existing, route) {
				return
			}
		}
		routes = append(routes, route)
	}
	for _, subscriber := range r.subscribers {
		if index := slices.IndexFunc(subscriber.Match, func(block alerts.Block) bool {
			return alerts.Matches(block.Fields["test"], test)
		}); index >= 0 {
			add(subscriber.blocks[index])
		} else if r.claims(subscriber.Owners, owner) {
			add(subscriber.route)
		}
	}
	if len(routes) == 0 {
		return []unitTestRoute{r.fallback}
	}
	// One analysis is enough for the dashboards, so a route without a channel is needed only when it is alone.
	if withChannel := slices.DeleteFunc(slices.Clone(routes), func(route unitTestRoute) bool { return route.slack.Channel == "" }); len(withChannel) > 0 {
		return withChannel
	}
	return routes[:1]
}

// claims reports whether one of the claims covers the owner: it names the owner, or with `subgroups` a group above it.
func (r *unitTestRouter) claims(claims []alerts.Owner, owner string) bool {
	if owner == "" {
		return false
	}
	owner = r.currentName(owner)
	for _, claim := range claims {
		claimed := r.currentName(claim.Name)
		if claimed == owner {
			return true
		}
		if !claim.Subgroups {
			continue
		}
		// The limit protects from a cycle in the groups.
		for group, depth := r.parent[owner], 0; group != "" && depth < 20; group, depth = r.parent[group], depth+1 {
			if group == claimed {
				return true
			}
		}
	}
	return false
}

func (r *unitTestRouter) currentName(group string) string {
	if name, ok := r.canonical[group]; ok {
		return name
	}
	return group
}
