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

	fmt.Println("Config: ", cfg)

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
		if strings.HasPrefix(repo.Path, cfg.GroupWhitelist[0]) {
			wg.Add(1)

			go func(repository *clients.GitRepository) {
				_ = sem.Acquire(context.Background(), 1)
				defer sem.Release(1)
				defer wg.Done()
				output, err := handleRepo(cfg, repository)

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(output)
				}

			}(repo)
		}
	}

	wg.Wait()

	return nil
}

func handleRepo(cfg *config.RfetchConfig, repo *clients.GitRepository) (string, error) {
	repoPath := path.Join(cfg.GetRootPathAbs(), repo.Path)

	if _, err := os.Stat(repoPath); errors.Is(err, fs.ErrNotExist) {
		parentPath := path.Dir(repoPath)
		err := os.MkdirAll(parentPath, os.ModePerm)

		if err != nil {
			return "", err
		}

		err = os.Chdir(path.Dir(repoPath))

		if err != nil {
			return "", err
		}

		cmd := exec.Command("git", "clone", repo.SshUrl)
		err = cmd.Run()

		if err != nil {
			return "", err
		}

		return "Cloned Repository: " + repo.Path, nil
	}

	err := os.Chdir(repoPath)

	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "fetch", "--all")
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	return "Updated repository: " + repo.Path, nil
}
