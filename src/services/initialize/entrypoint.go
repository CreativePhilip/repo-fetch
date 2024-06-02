package initialize

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func Init(ctx *cli.Context) error {
	cfg, err := GetUserData()

	if err != nil {
		return err
	}

	raw, err := yaml.Marshal(cfg)

	if err != nil {
		return err
	}

	err = os.WriteFile("rfetch.config.yml", raw, 0644)

	if err != nil {
		return err
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	fmt.Printf("Saved configuration to: %s\n", filepath.Join(cwd, cfg.RootPath))
	fmt.Println("To fetch the repositories run `rfetch sync`")
	return nil
}
