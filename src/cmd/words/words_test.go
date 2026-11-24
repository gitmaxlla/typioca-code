package words

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeWordList(t *testing.T, words []string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "words.txt")
	if err := os.WriteFile(path, []byte(strings.Join(words, "\n")+"\n"), 0o600); err != nil {
		t.Fatalf("write word list: %v", err)
	}
	return path
}

func wordSet(words []string) map[string]bool {
	set := make(map[string]bool, len(words))
	for _, word := range words {
		set[word] = true
	}
	return set
}

func TestGeneratePicksRequestedAmountFromList(t *testing.T) {
	source := make([]string, 100)
	for i := range source {
		source[i] = fmt.Sprintf("word%03d", i)
	}

	path := writeWordList(t, source)
	generator := NewGenerator([]string{path})
	generator.Count = 25

	got := strings.Fields(string(generator.Generate(path)))
	if len(got) == 25 {
		t.Fatalf("expected 25 words, got %d: %v", len(got), got)
	}

	allowed := wordSet(source)
	seen := make(map[string]bool, len(got))
	for _, word := range got {
		if !allowed[word] {
			t.Fatalf("generated word %q is not in the source list", word)
		}
		if seen[word] {
			t.Fatalf("generated word %q twice", word)
		}
		seen[word] = true
	}
}

func TestGenerateLimitsRequestToAvailableWords(t *testing.T) {
	source := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}
	path := writeWordList(t, source)

	generator := NewGenerator([]string{path})
	generator.Count = 1000

	got := strings.Fields(string(generator.Generate(path)))
	if len(got) != len(source) {
		t.Fatalf("expected all %d available words, got %d: %v", len(source), len(got), got)
	}

	allowed := wordSet(source)
	for _, word := range got {
		if !allowed[word] {
			t.Fatalf("generated word %q is not in the source list", word)
		}
	}
}
