package sync

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/semaphore"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"repo-fetch/src/clients"
	"repo-fetch/src/config"
	"strings"
	"sync"
)

func Sync(ctx *cli.Context) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	fmt.Println(cwd)
	cfg, err := config.WalkFsAndLoadConfig(cwd)

	if err != nil {
		return err
	}

	gitClient, err := clients.CreateProviderFromConfig(cfg)

	if err != nil {
		return err
	}

	repos, err := gitClient.ListRepositories()

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(10)

	for _, repo := range repos {
		if !strings.HasPrefix(repo.Path, cfg.GroupWhitelist[0]) {
			continue
		}

		wg.Add(1)
		go handleRepoAsync(cfg, repo, &wg, sem)
	}

	wg.Wait()

	return nil
}

func handleRepoAsync(cfg *config.RfetchConfig, repo *clients.GitRepository, wg *sync.WaitGroup, sem *semaphore.Weighted) {
	_ = sem.Acquire(context.Background(), 1)
	defer sem.Release(1)
	defer wg.Done()
	output, err := handleRepo(cfg, repo)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}
}

func handleRepo(cfg *config.RfetchConfig, repo *clients.GitRepository) (string, error) {
	repoPath := path.Join(cfg.GetRootPathAbs(), repo.Path)

	// Path does not exist, therefore repository does not exist, needs cloning
	if _, err := os.Stat(repoPath); errors.Is(err, fs.ErrNotExist) {
		err = cloneRepo(repo, path.Dir(repoPath))

		if err != nil {
			return "", err
		}

		return "Cloned Repository: " + repo.Path, nil
	}

	err := updateRepo(repo, repoPath)

	if err != nil {
		return "", err
	}

	return "Updated repository: " + repo.Path, nil
}

func cloneRepo(repo *clients.GitRepository, repoPath string) error {
	err := os.MkdirAll(repoPath, os.ModePerm)

	if err != nil {
		return err
	}

	err = os.Chdir(repoPath)

	if err != nil {
		return err
	}

	return exec.Command("git", "clone", repo.SshUrl).Run()
}

func updateRepo(repo *clients.GitRepository, repoPath string) error {
	err := os.Chdir(repoPath)

	if err != nil {
		return err
	}

	return exec.Command("git", "fetch", "--all").Run()
}
