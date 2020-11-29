package generatetest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertEqualContentLayout(t *testing.T, expectedContentLayout map[string]string, actual map[string][]byte) {
	require.Equal(t, len(expectedContentLayout), len(actual))
	for path, expectedContent := range expectedContentLayout {
		actualContent, exists := actual[path]
		if !exists {
			t.Fatalf("expectedContentLayout path %q does not exist in actual %+v", path, actual)
		}
		assert.Equal(t, expectedContent, string(actualContent))
	}
}
