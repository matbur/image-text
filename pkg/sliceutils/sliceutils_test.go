package sliceutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matbur/image-text/pkg/sliceutils"
)

func TestCoalesce(t *testing.T) {
	type TestCase[T comparable] struct {
		name  string
		items []T
		want  T
	}

	tests := []TestCase[string]{
		{
			name:  "empty",
			items: []string{},
			want:  "",
		},
		{
			name:  "all empty",
			items: []string{"", ""},
			want:  "",
		},
		{
			name:  "one of one",
			items: []string{"a"},
			want:  "a",
		},
		{
			name:  "1st empty, get 2nd",
			items: []string{"", "b"},
			want:  "b",
		},
		{
			name:  "both not empty, get 1st",
			items: []string{"a", "b"},
			want:  "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceutils.Coalesce(tt.items...)
			assert.Equal(t, tt.want, got)
		})
	}

	testsPtrStr := []TestCase[*string]{
		{
			name:  "empty",
			items: []*string{},
			want:  nil,
		},
		{
			name:  "all nils",
			items: []*string{nil, nil},
			want:  nil,
		},
		{
			name:  "one of one",
			items: []*string{toPointer("a")},
			want:  toPointer("a"),
		},
		{
			name:  "1st nil, get 2nd",
			items: []*string{nil, toPointer("b")},
			want:  toPointer("b"),
		},
		{
			name:  "both not nil, get 1st",
			items: []*string{toPointer("a"), toPointer("b")},
			want:  toPointer("a"),
		},
	}

	for _, tt := range testsPtrStr {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceutils.Coalesce(tt.items...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func toPointer[T any](t T) *T {
	return &t
}
