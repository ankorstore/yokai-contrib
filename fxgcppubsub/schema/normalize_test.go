package schema_test

import (
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub/schema"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeSchemaID(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		in          string
		expectedOut string
	}{
		"empty": {
			in:          "",
			expectedOut: "",
		},
		"without project info": {
			in:          "bar",
			expectedOut: "bar",
		},
		"with project info": {
			in:          "projects/foo/schemas/bar",
			expectedOut: "bar",
		},
	}

	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expectedOut, schema.NormalizeSchemaID(tc.in))
		})
	}
}
