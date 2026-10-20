package proc

import (
	"os"
	"path/filepath"
	"testing"
)

func mustMkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectGitInfo(t *testing.T) {
	root := t.TempDir()

	// Normal repo: <root>/repo/.git/ is a directory, HEAD on "main".
	repo := filepath.Join(root, "repo")
	repoGit := filepath.Join(repo, ".git")
	mustMkdir(t, repoGit)
	mustWrite(t, filepath.Join(repoGit, "HEAD"), "ref: refs/heads/main\n")

	if name, branch := detectGitInfo(repo); name != "repo" || branch != "main" {
		t.Errorf("dir repo: got (%q, %q), want (repo, main)", name, branch)
	}

	// A nested subdirectory walks up to the repo.
	sub := filepath.Join(repo, "a", "b")
	mustMkdir(t, sub)
	if name, branch := detectGitInfo(sub); name != "repo" || branch != "main" {
		t.Errorf("subdir: got (%q, %q), want (repo, main)", name, branch)
	}

	// Worktree: <root>/wt/.git is a FILE pointing at a worktree gitdir under the
	// main repo, which carries its own HEAD on "feature".
	wt := filepath.Join(root, "wt")
	mustMkdir(t, wt)
	wtGitDir := filepath.Join(repoGit, "worktrees", "wt")
	mustMkdir(t, wtGitDir)
	mustWrite(t, filepath.Join(wtGitDir, "HEAD"), "ref: refs/heads/feature\n")
	mustWrite(t, filepath.Join(wt, ".git"), "gitdir: "+wtGitDir+"\n")

	if name, branch := detectGitInfo(wt); name != "wt" || branch != "feature" {
		t.Errorf("worktree: got (%q, %q), want (wt, feature)", name, branch)
	}

	// No .git within reach: the walk caps at 10 ancestors, so a directory 12
	// levels below <root>/repo can never see <root>/repo/.git. Starting from
	// inside the repo (not beside it) also keeps the cap honest: raising the
	// cap would surface repo/.git and fail this assertion. The probe stays
	// hermetic on hosts where the temp dir sits inside a repository (dotfiles
	// checkouts, agent sandboxes, CI workspaces) — detectGitInfo cannot
	// escape far enough to find those, and must report empty.
	deep := repo
	for i := 0; i < 12; i++ {
		deep = filepath.Join(deep, "d")
		mustMkdir(t, deep)
	}
	if name, branch := detectGitInfo(deep); name != "" || branch != "" {
		t.Errorf("no repo in reach: got (%q, %q), want empty", name, branch)
	}
}

func TestGitDirFromFile(t *testing.T) {
	dir := t.TempDir()

	// Absolute pointer: returned unchanged. Derive the target from the temp dir
	// so it is absolute on every OS — a bare "/abs/..." is not absolute on
	// Windows, where git writes drive-letter paths like C:/repo/.git/worktrees/x.
	absTarget := filepath.Join(dir, "real", ".git", "worktrees", "x")
	mustWrite(t, filepath.Join(dir, "abs.git"), "gitdir: "+absTarget+"\n")
	if got := gitDirFromFile(filepath.Join(dir, "abs.git"), dir); got != absTarget {
		t.Errorf("absolute: got %q, want %q", got, absTarget)
	}

	mustWrite(t, filepath.Join(dir, "rel.git"), "gitdir: ../parent/.git/worktrees/x\n")
	want := filepath.Join(dir, "../parent/.git/worktrees/x")
	if got := gitDirFromFile(filepath.Join(dir, "rel.git"), dir); got != want {
		t.Errorf("relative: got %q, want %q", got, want)
	}

	if got := gitDirFromFile(filepath.Join(dir, "missing.git"), dir); got != "" {
		t.Errorf("missing file: got %q, want empty", got)
	}
}

func TestGitBranchFromHEAD(t *testing.T) {
	dir := t.TempDir()

	mustWrite(t, filepath.Join(dir, "HEAD"), "ref: refs/heads/dev\n")
	if got := gitBranchFromHEAD(dir); got != "dev" {
		t.Errorf("ref HEAD: got %q, want dev", got)
	}

	// Detached HEAD (raw commit) yields no branch.
	mustWrite(t, filepath.Join(dir, "HEAD"), "0123456789abcdef0123456789abcdef01234567\n")
	if got := gitBranchFromHEAD(dir); got != "" {
		t.Errorf("detached HEAD: got %q, want empty", got)
	}
}
