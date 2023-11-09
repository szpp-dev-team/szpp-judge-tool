package storage

import (
	"encoding/json"
	"os"
)

type Storage struct {
	storagePath  string
	TaskIDByPath map[string]int
}

func New(path string) *Storage {
	return &Storage{
		storagePath:  path,
		TaskIDByPath: make(map[string]int),
	}
}

func (s *Storage) SaveAsFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

func (s *Storage) SetTaskID(taskPath string, taskID int) {
	s.TaskIDByPath[taskPath] = taskID
}

func (s *Storage) GetTaskID(taskPath string) (int, bool) {
	id, ok := s.TaskIDByPath[taskPath]
	return id, ok
}
