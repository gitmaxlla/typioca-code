package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/kirsle/configdir"
	"github.com/muesli/termenv"
)

// prepareTUIConfig points the config directories at temporary paths and returns
// the path of the JSON config file that the TUI writes.
func prepareTUIConfig(t *testing.T) string {
	t.Helper()

	cacheDir := t.TempDir()
	configDir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheDir)
	t.Setenv("XDG_CONFIG_HOME", configDir)
	configdir.Refresh()

	return filepath.Join(cacheDir, "typioca", "typioca.conf")
}

// startConfigView runs the TUI, opens the Config screen and returns the test
// model together with a quit function.
func startConfigView(t *testing.T) *teatest.TestModel {
	t.Helper()

	tm := teatest.NewTestModel(
		t,
		initialModel(termenv.Ascii, termenv.NoColor{}, 100, 40),
		teatest.WithInitialTermSize(100, 40),
	)

	teatest.WaitFor(t, tm.Output(), func(out []byte) bool {
		return bytes.Contains(out, []byte("Timer run"))
	}, teatest.WithDuration(5*time.Second))

	// The main menu lists timer run, word count run, sentence count run and
	// config; move the cursor to the last entry and open it.
	for i := 0; i < 3; i++ {
		tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	}
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	teatest.WaitFor(t, tm.Output(), func(out []byte) bool {
		return bytes.Contains(out, []byte("  wordlist"))
	}, teatest.WithDuration(5*time.Second))

	return tm
}

func stopTestModel(t *testing.T, tm *teatest.TestModel) {
	t.Helper()

	if err := tm.Quit(); err != nil {
		t.Fatalf("quit TUI: %v", err)
	}
	tm.WaitFinished(t, teatest.WithFinalTimeout(5*time.Second))
}

func readPersistedConfig(t *testing.T, path string) (Config, error) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

// waitForEmbeddedEnabled waits until the first embedded word list has the
// requested enabled flag in the config file and returns the persisted config.
func waitForEmbeddedEnabled(t *testing.T, path string, enabled bool) Config {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		config, err := readPersistedConfig(t, path)
		switch {
		case err != nil:
			lastErr = err
		case len(config.EmbededWordLists) == 0:
			lastErr = os.ErrNotExist
		case config.EmbededWordLists[0].Enabled == enabled:
			return config
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("config was not written with enabled=%v for %q: %v", enabled, "Common words", lastErr)
	return Config{}
}

func TestToggleConfigDisablesWordList(t *testing.T) {
	configPath := prepareTUIConfig(t)

	initial := ReadConfig()
	if len(initial.EmbededWordLists) == 0 || !initial.EmbededWordLists[0].Enabled {
		t.Fatalf("expected the first embedded word list to start enabled: %+v", initial.EmbededWordLists)
	}

	tm := startConfigView(t)
	defer stopTestModel(t, tm)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	persisted := waitForEmbeddedEnabled(t, configPath, false)
	if persisted.EmbededWordLists[0].Name != "Common words" {
		t.Fatalf("toggled unexpected word list: %q", persisted.EmbededWordLists[0].Name)
	}
	if !persisted.EmbededWordLists[1].Enabled {
		t.Fatalf("toggling one word list changed another one: %+v", persisted.EmbededWordLists)
	}
}

func TestToggleConfigHandlesTenConsecutiveChanges(t *testing.T) {
	configPath := prepareTUIConfig(t)

	tm := startConfigView(t)
	defer stopTestModel(t, tm)

	enabled := true
	for i := 1; i <= 10; i++ {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		enabled = !enabled

		persisted := waitForEmbeddedEnabled(t, configPath, enabled)
		if len(persisted.EmbededWordLists) != 2 {
			t.Fatalf("change %d: expected two embedded word lists, got %d", i, len(persisted.EmbededWordLists))
		}
	}

	if !enabled {
		t.Fatal("after ten toggles the word list must be back in its initial state")
	}
}
