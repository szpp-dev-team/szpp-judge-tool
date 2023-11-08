package contest

import (
	"log/slog"
	"os"
	"path/filepath"

	pkgtask "github.com/szpp-dev-team/szpp-judge-tool/internal/task"
	"gopkg.in/yaml.v3"
)

type Contest struct {
	dir    string
	config *Config
	logger *slog.Logger
}

func Load(contestPath string) (*Contest, error) {
	f, err := os.Open(filepath.Join(contestPath, "contest.yaml"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &Contest{
		dir:    contestPath,
		config: &cfg,
		logger: slog.Default(),
	}, nil
}

func (c *Contest) Validate() error {
	for _, task := range c.config.Tasks {
		c.logger.Info("validating task", slog.String("task", task.Slug))
		controller, err := pkgtask.Load(filepath.Join(c.dir, task.Slug))
		if err != nil {
			return err
		}
		defer controller.Cleanup()
		if err := controller.Validate(); err != nil {
			return err
		}
	}

	return nil
}
