package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kirsle/configdir"
)

func prepareResultsCache(t *testing.T) string {
	t.Helper()

	cacheDir := t.TempDir()
	configDir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheDir)
	t.Setenv("XDG_CONFIG_HOME", configDir)
	configdir.Refresh()

	return filepath.Join(cacheDir, "typioca", "results.json")
}

func testResults(wpm int) Results {
	return Results{
		identifier: ResultsIdentifier{
			testType: "WordCountBasedTest",
			numeric:  25,
			words:    "Common words",
		},
		wpm:           wpm,
		accuracy:      100,
		deltaWpm:      0,
		rawWpm:        wpm,
		cpm:           wpm * 5,
		wordList:      "Common words",
		wpmEachSecond: []float64{float64(wpm)},
	}
}

func readResultsFile(t *testing.T, path string) PersistentResults {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read results file: %v", err)
	}

	var persisted PersistentResults
	if err := json.Unmarshal(data, &persisted); err != nil {
		t.Fatalf("decode results file: %v", err)
	}
	return persisted
}

func TestPersistResultsAddsResultToExistingFile(t *testing.T) {
	path := prepareResultsCache(t)

	first := PersistResults(testResults(60))
	nodes := first.Results["WordCountBasedTest"][25]["Common words"]
	if len(nodes) != 1 || nodes[0].Wpm != 60 {
		t.Fatalf("unexpected first result: %+v", nodes)
	}

	second := PersistResults(testResults(75))
	nodes = second.Results["WordCountBasedTest"][25]["Common words"]
	if len(nodes) != 2 {
		t.Fatalf("expected two results, got %d", len(nodes))
	}
	if nodes[0].Wpm != 60 || nodes[1].Wpm != 75 {
		t.Fatalf("unexpected results order: %+v", nodes)
	}

	persisted := readResultsFile(t, path)
	nodes = persisted.Results["WordCountBasedTest"][25]["Common words"]
	if len(nodes) != 2 {
		t.Fatalf("expected two results in file, got %d", len(nodes))
	}
	if nodes[1].Wpm != 75 {
		t.Fatalf("new result is not in the file: %+v", nodes)
	}
}

func TestPersistResultsCreatesFileForFirstResult(t *testing.T) {
	path := prepareResultsCache(t)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("results file already exists before the first result: %v", err)
	}

	PersistResults(testResults(42))

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("results file was not created: %v", err)
	}

	persisted := readResultsFile(t, path)
	nodes := persisted.Results["WordCountBasedTest"][25]["Common words"]
	if len(nodes) != 1 {
		t.Fatalf("expected one result in a new file, got %d", len(nodes))
	}
	if nodes[0].Wpm != 42 {
		t.Fatalf("unexpected first result: %+v", nodes[0])
	}
}
