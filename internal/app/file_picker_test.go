package app

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFilePickerEntriesOrdersDirectoriesAndFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"zeta", ".hidden-dir", "Alpha"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0755); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"zeta.txt", ".hidden.txt", "Alpha.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	gotEntries := filePickerEntries(dir)
	got := make([]string, len(gotEntries))
	for i, entry := range gotEntries {
		got[i] = entry.label
	}
	want := []string{"..", "Alpha", "zeta", ".hidden-dir", "Alpha.txt", "zeta.txt", ".hidden.txt"}
	if !slices.Equal(got, want) {
		t.Fatalf("entries = %v, want %v", got, want)
	}
	for _, entry := range gotEntries[1:4] {
		if !entry.dir {
			t.Errorf("entry %q should be a directory", entry.label)
		}
	}
	for _, entry := range gotEntries[4:] {
		if entry.dir {
			t.Errorf("entry %q should be a file", entry.label)
		}
	}
}

func TestFilterFilePickerEntriesKeepsParent(t *testing.T) {
	entries := []filePickerEntry{
		{label: ".."},
		{label: "app.go"},
		{label: "README.md"},
	}
	got := filterFilePickerEntries(entries, "AP")
	if len(got) != 2 || got[0].label != ".." || got[1].label != "app.go" {
		t.Fatalf("filtered entries = %v, want [.. app.go]", got)
	}
}
