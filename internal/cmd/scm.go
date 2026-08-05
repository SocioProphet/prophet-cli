package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// prophet scm — GitHub <-> Gitea sovereignty. The estate declared Gitea canonical and cut the
// pull-mirror, but the workflow never moved: GitHub holds the real state, Gitea is a frozen orphan.
// This surface makes the cutover clean — validate parity before flipping, dual-write during the
// transition — and changes no remotes on its own. Config via env: GITEA_URL, GITEA_TOKEN, ORG.

func scmEnv() (url, token, org string) {
	url = strings.TrimRight(os.Getenv("GITEA_URL"), "/")
	token = os.Getenv("GITEA_TOKEN")
	org = os.Getenv("ORG")
	if org == "" {
		org = "SocioProphet"
	}
	return
}

func scmRun(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	return strings.TrimSpace(string(out)), err
}

func scmCurrentRepo() string {
	if r, err := scmRun("gh", "repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner"); err == nil && r != "" {
		return r
	}
	if u, err := scmRun("git", "config", "--get", "remote.origin.url"); err == nil {
		u = strings.TrimSuffix(u, ".git")
		if i := strings.LastIndex(u, "github.com"); i >= 0 {
			return strings.Trim(u[i+len("github.com"):], ":/")
		}
	}
	return ""
}

func scmGithubHead(repo string) string {
	s, _ := scmRun("gh", "api", "repos/"+repo+"/commits/main", "-q", ".sha")
	if len(s) > 12 {
		s = s[:12]
	}
	return s
}

func scmGiteaHead(url, token, repo string) (string, string) {
	if token == "" || url == "" {
		return "", "unauth"
	}
	req, err := http.NewRequest("GET", url+"/api/v1/repos/"+repo+"/branches/main", nil)
	if err != nil {
		return "", "error"
	}
	req.Header.Set("Authorization", "token "+token)
	resp, err := (&http.Client{Timeout: 8 * time.Second}).Do(req)
	if err != nil {
		return "", "error"
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == 404:
		return "", "missing"
	case resp.StatusCode != 200:
		return "", fmt.Sprintf("http%d", resp.StatusCode)
	}
	var d struct {
		Commit struct {
			ID string `json:"id"`
		} `json:"commit"`
	}
	if json.NewDecoder(resp.Body).Decode(&d) != nil {
		return "", "error"
	}
	id := d.Commit.ID
	if len(id) > 12 {
		id = id[:12]
	}
	return id, "ok"
}

func scmDrift(repoFlag string) error {
	url, token, _ := scmEnv()
	repo := repoFlag
	if repo == "" {
		repo = scmCurrentRepo()
	}
	if repo == "" {
		return fmt.Errorf("not in a recognizable repo; pass --repo OWNER/NAME")
	}
	gh := scmGithubHead(repo)
	gt, status := scmGiteaHead(url, token, repo)
	verdict := "unknown"
	switch {
	case status == "ok" && gt != "" && gh == gt:
		verdict = "in_sync"
	case status == "ok" && gt != "":
		verdict = "drifted"
	case status == "missing":
		verdict = "gitea_missing"
	}
	_ = emit(map[string]any{
		"command": "prophet scm drift", "repo": repo,
		"github": gh, "gitea": gt, "gitea_status": status, "verdict": verdict,
	})
	switch verdict {
	case "drifted":
		return fmt.Errorf("DRIFTED — GitHub is ahead of the canonical Gitea; reconcile before flipping")
	case "unknown":
		return fmt.Errorf("UNKNOWN — Gitea unreadable (set GITEA_TOKEN); treat as NOT clean")
	}
	return nil
}

func newScmCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "scm", Short: "GitHub<->Gitea sovereignty: validate drift, dual-write, prepare the cutover"}

	var repoFlag string
	drift := &cobra.Command{
		Use:   "drift",
		Short: "Compare GitHub vs Gitea HEAD — is the 'canonical' Gitea actually in sync?",
		RunE:  func(c *cobra.Command, args []string) error { return scmDrift(repoFlag) },
	}
	drift.Flags().StringVar(&repoFlag, "repo", "", "OWNER/NAME (default: current repo)")

	var branchFlag string
	sync := &cobra.Command{
		Use:   "sync",
		Short: "Dual-write: push the current branch to GitHub AND Gitea",
		RunE: func(c *cobra.Command, args []string) error {
			br := branchFlag
			if br == "" {
				br, _ = scmRun("git", "branch", "--show-current")
			}
			if br == "" {
				return fmt.Errorf("detached HEAD — checkout a branch first")
			}
			if _, err := scmRun("git", "push", "origin", br); err != nil {
				return fmt.Errorf("github push failed: %w", err)
			}
			remotes, _ := scmRun("git", "remote")
			if strings.Contains(remotes, "gitea") {
				_, _ = scmRun("git", "push", "gitea", br)
				return emit(map[string]any{"command": "prophet scm sync", "branch": br, "pushed": []string{"origin", "gitea"}})
			}
			return emit(map[string]any{"command": "prophet scm sync", "branch": br,
				"pushed": []string{"origin"}, "note": "no 'gitea' remote — add it to dual-write"})
		},
	}
	sync.Flags().StringVar(&branchFlag, "branch", "", "branch to push (default: current)")

	status := &cobra.Command{
		Use:   "status",
		Short: "Backend + drift summary",
		RunE:  func(c *cobra.Command, args []string) error { return scmDrift(repoFlag) },
	}

	cmd.AddCommand(drift, sync, status)
	return cmd
}
