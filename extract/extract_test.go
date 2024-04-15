package extract

import (
	"reflect"
	"testing"
)

func TestArgv(t *testing.T) {
	tests := []struct {
		name         string
		inputArgs    []string
		expectedPre  []string
		expectedPost []string
	}{
		{
			name:         "no dash",
			inputArgs:    []string{"/opt/entry", "arg1", "arg2"},
			expectedPre:  []string{"arg1", "arg2"},
			expectedPost: nil,
		},
		{
			name:         "with dash",
			inputArgs:    []string{"/opt/entry", "arg1", "--", "arg2", "arg3"},
			expectedPre:  []string{"arg1"},
			expectedPost: []string{"arg2", "arg3"},
		},
		{
			name:         "dash at start",
			inputArgs:    []string{"/opt/entry", "--", "arg1", "arg2"},
			expectedPre:  nil,
			expectedPost: []string{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre, post := Argv(tt.inputArgs)

			if !reflect.DeepEqual(pre, tt.expectedPre) {
				t.Errorf("Expected pre-dash args %v, got %v", tt.expectedPre, pre)
			}

			if !reflect.DeepEqual(post, tt.expectedPost) {
				t.Errorf("Expected post-dash args %v, got %v", tt.expectedPost, post)
			}
		})
	}
}
