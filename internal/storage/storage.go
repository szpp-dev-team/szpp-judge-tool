package storage

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Storage struct {
	filename string
	data     *Data
}

type Data struct {
	TaskIDByPath map[string]int `json:"tasks"`
}

func LoadOrInit(filename string) (*Storage, error) {
	var data Data
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		log.Println(err)
		if err == io.EOF {
			return &Storage{
				filename: filename,
				data: &Data{
					TaskIDByPath: make(map[string]int),
				},
			}, nil
		}
		return nil, err
	}
	log.Println(data)
	return &Storage{
		filename: filename,
		data:     &data,
	}, nil
}

func (s *Storage) Save() error {
	f, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s.data)
}

func (s *Storage) SetTaskID(taskPath string, taskID int) {
	s.data.TaskIDByPath[taskPath] = taskID
}

func (s *Storage) GetTaskID(taskPath string) (int, bool) {
	id, ok := s.data.TaskIDByPath[taskPath]
	return id, ok
}
