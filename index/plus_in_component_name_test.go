package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// A component whose name contains '+' is present in the document and referenced identically, yet was
// reported as nonexistent: the lookup query-unescaped the JSON Pointer fragment, turning '+' into a
// space. This is the end-to-end form of that -- resolve the ref through the index, not just convert the
// id -- because the conversion being right is not the same as the reference resolving.
func TestFindComponent_PlusInComponentName(t *testing.T) {
	spec := `openapi: 3.1.0
paths:
  /thing:
    get:
      responses:
        '200':
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Widget+Inner'
components:
  schemas:
    Widget+Inner:
      type: object
      properties:
        v:
          type: string`

	var root yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(spec), &root))

	idx := NewSpecIndexWithConfig(&root, CreateOpenAPIIndexConfig())
	assert.Empty(t, idx.GetReferenceIndexErrors(),
		"a component named Widget+Inner is present and referenced identically, so it must resolve")

	ref := idx.FindComponent("#/components/schemas/Widget+Inner")
	require.NotNil(t, ref, "FindComponent must locate a '+'-named component")
	assert.Equal(t, "#/components/schemas/Widget+Inner", ref.FullDefinition)
}
