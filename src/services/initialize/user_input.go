package initialize

import "repo-fetch/src/utils"

func GetUserData() (*InitConfiguration, error) {
	provider, err := utils.StdinSelectFrom(
		"Select your git provider",
		[]string{"GitHub", "GitLab", "BitBucket"},
	)

	if err != nil {
		return nil, err
	}

	tokenVar, err := utils.StdinInput("Provide the env var for the api token")

	if err != nil {
		return nil, err
	}

	rootPath, err := utils.StdinInput("Provide the path to the root to clone")

	if err != nil {
		return nil, err
	}

	return &InitConfiguration{
		Provider: provider,
		TokenVar: tokenVar,
		RootPath: rootPath,
	}, nil
}
