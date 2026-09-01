package services

import (
	"path/filepath"
	"runtime"
	"testing"
)

// Note: these tests exercise command *construction* only (never .Start()),
// since actually launching VS Code/a terminal would be a real, visible side
// effect on whatever machine runs `go test`.

func TestCollectionService_OpenInEditor_EmptyPath(t *testing.T) {
	s := NewCollectionService(nil)
	if err := s.OpenInEditor("system", "  "); err == nil {
		t.Fatal("expected an error for an empty path")
	}
}

func TestCollectionService_OpenInEditor_VSCodeNotFound(t *testing.T) {
	s := NewCollectionService(nil)
	t.Setenv("PATH", t.TempDir()) // no `code` binary reachable from here
	if err := s.OpenInEditor("vscode", "/tmp/x.json"); err == nil {
		t.Fatal("expected an error when the 'code' CLI isn't on PATH")
	}
}

func TestSystemOpenCommand_MatchesGOOS(t *testing.T) {
	cmd := systemOpenCommand("/tmp/x.json")
	var want string
	switch runtime.GOOS {
	case "darwin":
		want = "open"
	case "windows":
		want = "cmd"
	default:
		want = "xdg-open"
	}
	got := filepath.Base(cmd.Path)
	if runtime.GOOS == "windows" {
		got = got[:len(want)] // strip a possible ".exe" suffix
	}
	if got != want {
		t.Errorf("systemOpenCommand() path = %q, want %q", cmd.Path, want)
	}
}

func TestShellQuote(t *testing.T) {
	got := shellQuote(`it's a "test" path.json`)
	want := `'it'\''s a "test" path.json'`
	if got != want {
		t.Errorf("shellQuote() = %q, want %q", got, want)
	}
}

func TestAppleScriptQuote(t *testing.T) {
	got := appleScriptQuote(`say "hi" \ done`)
	want := `"say \"hi\" \\ done"`
	if got != want {
		t.Errorf("appleScriptQuote() = %q, want %q", got, want)
	}
}

func TestVimTerminalCommand_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux-specific")
	}
	cmd, err := vimTerminalCommand("/tmp/x.json")
	if err != nil {
		t.Fatalf("vimTerminalCommand() error = %v", err)
	}
	if filepath.Base(cmd.Path) != "x-terminal-emulator" {
		t.Errorf("vimTerminalCommand() path = %q, want x-terminal-emulator", cmd.Path)
	}
}

func TestVimTerminalCommand_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-specific")
	}
	cmd, err := vimTerminalCommand("/tmp/x.json")
	if err != nil {
		t.Fatalf("vimTerminalCommand() error = %v", err)
	}
	if filepath.Base(cmd.Path) != "osascript" {
		t.Errorf("vimTerminalCommand() path = %q, want osascript", cmd.Path)
	}
}
