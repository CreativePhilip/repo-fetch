package clients

import (
	"repo-fetch/src/config"
	"strings"
)

type GitProviderClient interface {
	ListRepositories() ([]*GitRepository, error)
}

type UnsupportedGitProviderError struct {
}

func (e UnsupportedGitProviderError) Error() string {
	return "The selected git provider is not supported"
}

type GitRepository struct {
	Id     string
	Name   string
	Path   string
	SshUrl string
}

func CreateProvider(provider string, token string) (GitProviderClient, error) {
	switch strings.ToLower(provider) {
	case "gitlab":
		return InitializeGitlabClient(token)
	default:
		return nil, UnsupportedGitProviderError{}
	}
}

func CreateProviderFromConfig(cfg *config.RfetchConfig) (GitProviderClient, error) {
	return CreateProvider(cfg.Provider, cfg.GetToken())
}
