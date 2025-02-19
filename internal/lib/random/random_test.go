package random_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"url-service/internal/lib/random"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "length 3", length: 3},
		{name: "length 5", length: 5},
		{name: "length 10", length: 10},
		{name: "length 30", length: 30},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s1 := random.NewRandomString(test.length)
			s2 := random.NewRandomString(test.length)

			assert.Len(t, s1, test.length)
			assert.Len(t, s2, test.length)

			assert.NotEqual(t, s1, s2)
		})
	}
}
