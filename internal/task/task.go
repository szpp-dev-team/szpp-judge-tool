package task

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/szpp-dev-team/szpp-judge-tool/internal/exec"
	"gopkg.in/yaml.v3"
)

type Task struct {
	Config    *Config
	Statement string
	dir       string
	logger    *slog.Logger
}

type Testcase struct {
	In  string
	Out string
}

func Load(taskPath string) (*Task, error) {
	f, err := os.Open(filepath.Join(taskPath, "task.yaml"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	// TODO: check if the statement contains "問題文", "制約", "入力", "出力"
	statement, err := os.ReadFile(filepath.Join(taskPath, "statement.md"))
	if err != nil {
		return nil, err
	}

	return &Task{
		Config:    &cfg,
		Statement: string(statement),
		dir:       taskPath,
		logger:    slog.Default().With(slog.String("task", taskPath)),
	}, nil
}

func (t *Task) ReadTestcase(slug string) (*Testcase, error) {
	in, err := os.ReadFile(filepath.Join(t.dir, "testcases", "in", slug))
	if err != nil {
		return nil, err
	}
	out, err := os.ReadFile(filepath.Join(t.dir, "testcases", "out", slug))
	if err != nil {
		return nil, err
	}
	return &Testcase{string(in), string(out)}, nil
}

func (t *Task) ReadChecker() (string, error) {
	b, err := os.ReadFile(filepath.Join(t.dir, "checker.cpp"))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (t *Task) Validate() error {
	if err := t.Config.validate(); err != nil {
		return err
	}

	if err := t.compile("correct.cpp"); err != nil {
		return fmt.Errorf("failed to compile correct.cpp: %v", err)
	}
	if err := t.compile("checker.cpp"); err != nil {
		return fmt.Errorf("failed to compile checker.cpp: %v", err)
	}
	if err := t.compile("verifier.cpp"); err != nil {
		return fmt.Errorf("failed to compile verifier.cpp: %v", err)
	}
	return t.validateTestcases()
}

func (t *Task) Cleanup() {
	t.logger.Info("cleanup")
	targets := []string{"correct", "checker", "verifier", "correct_stdout.txt"}
	for _, target := range targets {
		os.Remove(filepath.Join(t.dir, target))
	}
}

func (t *Task) compile(target string) error {
	cppCommand := getCppCommand()
	t.logger.Info("compiling", slog.String("target", target), slog.String("compiler", cppCommand))
	return exec.ExecuteCommand(cppCommand, []string{"-O2", "-std=gnu++17", "-o", basename(target), target}, exec.WithWorkdir(t.dir))
}

// testcase が正しいかを検証する
func (t *Task) validateTestcase(inPath, outPath string) error {
	t.logger.Info("verifying testcase", slog.String("testcase", inPath))

	inBytes, err := os.ReadFile(filepath.Join(t.dir, inPath))
	if err != nil {
		return err
	}

	// 入力形式の検証
	if err := exec.ExecuteCommand("./verifier", []string{}, exec.WithWorkdir(t.dir), exec.WithStdin(bytes.NewReader(inBytes))); err != nil {
		return err
	}

	outFile, err := os.OpenFile(filepath.Join(t.dir, "correct_stdout.txt"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// testcase の検証
	if err := exec.ExecuteCommand("./correct", []string{}, exec.WithWorkdir(t.dir), exec.WithStdin(bytes.NewReader(inBytes)), exec.WithStdout(outFile)); err != nil {
		return err
	}
	return exec.ExecuteCommand("./checker", []string{inPath, outPath, "correct_stdout.txt"}, exec.WithWorkdir(t.dir))
}

func (t *Task) validateTestcases() error {
	for _, testcase := range t.Config.Testcases {
		inPath := filepath.Join("testcases", "in", testcase.Slug+".txt")
		if !exists(inPath) {
			return fmt.Errorf("file %s does not exist", inPath)
		}
		outPath := filepath.Join("testcases", "out", testcase.Slug+".txt")
		if !exists(outPath) {
			return fmt.Errorf("file %s does not exist", outPath)
		}
		if err := t.validateTestcase(inPath, outPath); err != nil {
			return err
		}
	}
	return nil
}

func basename(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}

func getCppCommand() string {
	if cppCommand := os.Getenv("CXX"); cppCommand != "" {
		return cppCommand
	}
	return "g++"
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
