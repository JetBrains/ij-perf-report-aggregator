package tc_properties

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	// /Volumes/data/Downloads/build.finish.properties
	t.Setenv("CLICKHOUSE_DATA_PATH", "test")

	data, err := os.ReadFile("../../testData/build.finish.properties")
	require.NoError(t, err)
	properties, err := ReadProperties(data)
	require.NoError(t, err)

	err = os.WriteFile("/tmp/foo.txt", properties, 0o777)
	require.NoError(t, err)

	// err = backup()
	// if err == nil || !strings.HasPrefix(err.Error(), "can't connect to ") {
	//   t.Error(err)
	// }
}

func TestPropertiesToJson(t *testing.T) {
	t.Parallel()
	p := &Properties{m: map[string]string{
		"bool.true":       "true",
		"bool.false":      "false",
		"int":             "42",
		"int.negative":    "-7",
		"int.zero":        "0",
		"int.leadingZero": "007",
		"int.plus":        "+5",
		"int.huge":        "99999999999999999999999999",
		"string":          "hello",
		"string.escape":   `say "hi" \ <bye>`,
		`key"quote`:       "v",
		"empty":           "",
	}}

	data, err := propertiesToJson(p)
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(data, &parsed))

	assert.Equal(t, true, parsed["bool.true"])
	assert.Equal(t, false, parsed["bool.false"])
	assert.InDelta(t, float64(42), parsed["int"], 0)
	assert.InDelta(t, float64(-7), parsed["int.negative"], 0)
	assert.InDelta(t, float64(0), parsed["int.zero"], 0)
	// leading zeros and '+' are not valid JSON numbers, must stay strings
	assert.Equal(t, "007", parsed["int.leadingZero"])
	assert.Equal(t, "+5", parsed["int.plus"])
	assert.Contains(t, parsed, "int.huge")
	assert.Equal(t, "hello", parsed["string"])
	assert.Equal(t, `say "hi" \ <bye>`, parsed["string.escape"])
	assert.Equal(t, "v", parsed[`key"quote`])
	assert.NotContains(t, parsed, "empty")
}

func TestIsJsonInt(t *testing.T) {
	t.Parallel()
	valid := []string{"0", "5", "42", "-7", "-0", "99999999999999999999999999"}
	for _, s := range valid {
		assert.True(t, isJsonInt(s), s)
	}
	invalid := []string{"", "-", "007", "-01", "+5", "1a", "1.5", " 1"}
	for _, s := range invalid {
		assert.False(t, isJsonInt(s), s)
	}
}

func TestReadProperties(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../testData/build.finish.properties")
	require.NoError(t, err)
	properties, _ := LoadBytes(data, nil)
	v, ok := properties.Get("vcsroot.authMethod")
	assert.True(t, ok)
	assert.Equal(t, "TEAMCITY_SSH_KEY", v)
}
