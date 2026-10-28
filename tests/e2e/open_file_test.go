package e2e

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v3"
)

func TestOpenFileFromCommandPaletteOutsideWorkspace(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()

	outside := t.TempDir()
	path := filepath.Join(outside, "outside.txt")
	if err := os.WriteFile(path, []byte("outside workspace content"), 0644); err != nil {
		t.Fatal(err)
	}
	workspaceBefore := h.app.Workspace.Paths()

	h.exec("command.palette")
	h.typeText("Open File")
	h.pressKey(tcell.KeyEnter, tcell.ModNone)
	h.assertContains("Open File")

	h.pressKey(tcell.KeyCtrlA, tcell.ModCtrl)
	h.typeText(path)
	h.pressKey(tcell.KeyEnter, tcell.ModNone)

	if got := h.app.EditorGroup.ActiveFilePath(); got != path {
		t.Fatalf("active file = %q, want %q", got, path)
	}
	if got := strings.Join(h.app.EditorGroup.Editor.Buf.Lines, "\n"); got != "outside workspace content\n" {
		t.Fatalf("opened contents = %q", got)
	}
	if got := h.app.Workspace.Paths(); !slices.Equal(got, workspaceBefore) {
		t.Fatalf("workspace paths = %v, want unchanged %v", got, workspaceBefore)
	}
	if h.app.Root.HasOverlay() {
		t.Fatal("file picker should close after opening a file")
	}
}

func TestOpenFilePickerCancelLeavesEditorAndWorkspaceUnchanged(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()

	activeBefore := h.app.EditorGroup.ActiveFilePath()
	workspaceBefore := h.app.Workspace.Paths()
	h.exec("file.open")
	if !h.app.Root.HasOverlay() {
		t.Fatal("Open File should show the file picker")
	}
	h.pressKey(tcell.KeyEscape, tcell.ModNone)

	if got := h.app.EditorGroup.ActiveFilePath(); got != activeBefore {
		t.Fatalf("active file after cancel = %q, want %q", got, activeBefore)
	}
	if got := h.app.Workspace.Paths(); !slices.Equal(got, workspaceBefore) {
		t.Fatalf("workspace paths after cancel = %v, want unchanged %v", got, workspaceBefore)
	}
	if h.app.Root.HasOverlay() {
		t.Fatal("file picker should close after cancellation")
	}
}

func TestOpenFilePickerRejectsMissingFile(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()

	missing := filepath.Join(t.TempDir(), "missing.txt")
	h.exec("file.open")
	h.pressKey(tcell.KeyCtrlA, tcell.ModCtrl)
	h.typeText(missing)
	h.pressKey(tcell.KeyEnter, tcell.ModNone)

	h.assertContains("Cannot open file")
	if !h.app.Root.HasOverlay() {
		t.Fatal("picker should remain open when the requested file does not exist")
	}
	if h.app.EditorGroup.ActiveFilePath() == missing {
		t.Fatal("missing file should not be opened as an empty new file")
	}
	h.pressKey(tcell.KeyEscape, tcell.ModNone)
}

func TestOpenFilePickerExplainsDirectorySelection(t *testing.T) {
	h := newTestHarness(t, 100, 30)
	defer h.stop()

	dir := t.TempDir()
	h.exec("file.open")
	h.pressKey(tcell.KeyCtrlA, tcell.ModCtrl)
	h.typeText(dir)
	h.pressKey(tcell.KeyEnter, tcell.ModNone)

	h.assertContains("path is a directory")
	if !h.app.Root.HasOverlay() {
		t.Fatal("picker should remain open when the selected path is a directory")
	}
	h.pressKey(tcell.KeyEscape, tcell.ModNone)
}
