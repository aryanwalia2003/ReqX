package services

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"reqx/internal/errs"
)

// OpenInEditor launches path in an external application — "vscode" (the
// `code` CLI), "vim" (opened in a new terminal window), or anything else
// falls back to the OS's default handler for the file (same as
// double-clicking it in Finder/Explorer/a file manager).
func (s *CollectionService) OpenInEditor(kind string, path string) error {
	if strings.TrimSpace(path) == "" {
		return errs.InvalidInput("no file path to open")
	}

	var cmd *exec.Cmd
	switch kind {
	case "vscode":
		if _, err := exec.LookPath("code"); err != nil {
			return errs.InvalidInput("VS Code CLI ('code') not found on PATH — install it via the Command Palette: Shell Command: Install 'code' command in PATH")
		}
		cmd = exec.Command("code", path)
	case "vim":
		var err error
		cmd, err = vimTerminalCommand(path)
		if err != nil {
			return err
		}
	default:
		cmd = systemOpenCommand(path)
	}

	if err := cmd.Start(); err != nil {
		return errs.Wrap(err, errs.KindExternal, "could not launch "+kind)
	}
	return nil
}

func systemOpenCommand(path string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", path)
	case "windows":
		// cmd's built-in `start` needs an empty title arg before the path
		// so a path containing spaces isn't mistaken for the window title.
		return exec.Command("cmd", "/c", "start", "", path)
	default:
		return exec.Command("xdg-open", path)
	}
}

// vimTerminalCommand opens vim in a new terminal window. Vim is a TUI app,
// so unlike VS Code/the system opener this can't just exec it detached — it
// needs an actual terminal to attach to, which is inherently OS-specific.
//
// ponytail: macOS goes through Terminal.app (well-defined single target);
// Linux tries only x-terminal-emulator (the Debian/Ubuntu default alias) —
// gnome-terminal/konsole/etc. aren't probed. Widen the Linux search (or add
// a Windows path) if that ceiling is hit.
func vimTerminalCommand(path string) (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`tell application "Terminal" to do script %s`,
			appleScriptQuote("vim "+shellQuote(path)))
		return exec.Command("osascript", "-e", script), nil
	case "linux":
		if _, err := exec.LookPath("x-terminal-emulator"); err != nil {
			return nil, errs.InvalidInput("no terminal emulator found (looked for x-terminal-emulator) — open a terminal and run: vim " + path)
		}
		return exec.Command("x-terminal-emulator", "-e", "vim", path), nil
	default:
		return nil, errs.InvalidInput("opening vim in a terminal isn't supported on this OS yet")
	}
}

// shellQuote wraps s in single quotes for safe use inside a POSIX shell
// command string, escaping embedded single quotes the standard way.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// appleScriptQuote wraps s in double quotes for safe embedding inside an
// AppleScript string literal passed to `osascript -e`.
func appleScriptQuote(s string) string {
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}
