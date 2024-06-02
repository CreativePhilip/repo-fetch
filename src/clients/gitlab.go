package clients

import (
	"github.com/xanzy/go-gitlab"
	"strconv"
)

type GitlabClient struct {
	client *gitlab.Client
}

func InitializeGitlabClient(token string) (*GitlabClient, error) {
	client, err := gitlab.NewClient(token)

	if err != nil {
		return nil, err
	}

	return &GitlabClient{
		client: client,
	}, nil
}

func (c *GitlabClient) ListRepositories() ([]*GitRepository, error) {
	repositories := make([]*GitRepository, 0)

	requestOptions := &gitlab.ListGroupProjectsOptions{
		Archived:         gitlab.Ptr(false),
		IncludeSubGroups: gitlab.Ptr(true),
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}
	for {
		projects, response, err := c.client.Groups.ListGroupProjects("4311008", requestOptions)

		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			repositories = append(repositories, &GitRepository{
				Id:     strconv.Itoa(project.ID),
				Name:   project.Name,
				Path:   project.PathWithNamespace,
				SshUrl: project.SSHURLToRepo,
			})
		}

		if response.CurrentPage >= response.TotalPages {
			break
		}

		requestOptions.ListOptions = gitlab.ListOptions{
			Page:    response.CurrentPage + 1,
			PerPage: 100,
		}
	}

	return repositories, nil
}
