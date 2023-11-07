package task

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Title         string                  `yaml:"title" validate:"required"`                // 問題名
	TimeLimitMs   uint32                  `yaml:"time_limit" validate:"gte=100,lte=10000"`  // 実行時間制限
	MemoryLimitMb uint32                  `yaml:"memory_limit" validate:"gte=128,lte=1024"` // 実行メモリ制限
	Difficulty    string                  `yaml:"difficulty" validate:"required"`           // 難易度
	TestcaseSets  map[string]*TestcaseSet `yaml:"testcase_sets"`
	Testcases     []*Testcase             `yaml:"testcases"`
}

type TestcaseSet struct {
	ScoreRatio    uint32   `yaml:"score_ratio"` // 得点比率(総和が100になるように)
	TestcaseSlugs []string `yaml:"list"`        // その TestcaseSet に属する Testcase の Slug 一覧
	IsSample      bool     `yaml:"is_sample"`   // サンプルかどうか
}

type Testcase struct {
	Slug        string `yaml:"name"`        // TestcaseSlug
	Description string `yaml:"description"` // 説明
}

func (c *Config) validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}

	scoreRatioSum := uint32(0)
	for _, ts := range c.TestcaseSets {
		scoreRatioSum += ts.ScoreRatio
	}
	if scoreRatioSum != 100 {
		return errors.New("the sum of scoreRatios must be 100")
	}

	return nil
}
