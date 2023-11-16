package storage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Save(t *testing.T) {
	jsonPath := filepath.Join(t.TempDir(), "data.json")

	s, err := LoadOrInit(jsonPath)
	require.NoError(t, err)
	s.SetTaskID("/hoge", 1)
	s.SetTaskID("/fuga", 2)
	require.NoError(t, s.Save())

	s, err = LoadOrInit(jsonPath)
	require.NoError(t, err)
	hogeTaskID, ok := s.GetTaskID("/hoge")
	require.True(t, ok)
	assert.Equal(t, 1, hogeTaskID)
	fugaTaskID, ok := s.GetTaskID("/fuga")
	require.True(t, ok)
	assert.Equal(t, 2, fugaTaskID)
	_, ok = s.GetTaskID("/piyo")
	require.False(t, ok)
}
