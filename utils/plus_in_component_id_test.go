package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// A $ref fragment is a JSON Pointer, in which '+' is a literal character: only '~0' and '~1' are
// escapes there. QueryUnescape applies x-www-form-urlencoded semantics and decodes '+' as a space, so a
// component whose name contains '+' became unresolvable even though it was present in the document.
//
// .NET OpenAPI generators emit '+' for nested types -- it is the CLR nested-type separator -- so this
// affected real published documents, e.g. Aruba's `ApiKeyQuotasDto+QuotaValue`.
func TestConvertComponentIdIntoFriendlyPathSearch_PlusIsLiteral(t *testing.T) {
	name, search := ConvertComponentIdIntoFriendlyPathSearch("#/components/schemas/Widget+Inner")
	assert.Equal(t, "Widget+Inner", name, "'+' is literal in a JSON Pointer and must survive unescaping")
	assert.Equal(t, "$.components.schemas['Widget+Inner']", search)
}

// Percent-encoding must still decode: PathUnescape percent-decodes, it simply leaves '+' alone.
func TestConvertComponentIdIntoFriendlyPathSearch_PercentStillDecodes(t *testing.T) {
	name, _ := ConvertComponentIdIntoFriendlyPathSearch("#/components/schemas/Widget%20Inner")
	assert.Equal(t, "Widget Inner", name, "%20 still decodes to a space")

	name, _ = ConvertComponentIdIntoFriendlyPathSearch("#/components/schemas/Widget%2BInner")
	assert.Equal(t, "Widget+Inner", name, "%2B decodes to '+', which then stays '+'")
}

// A literal space keeps working, and '~1' still means '/'. These passed before the change and must
// keep passing: the fix narrows what is decoded, it does not stop decoding.
func TestConvertComponentIdIntoFriendlyPathSearch_UnaffectedCases(t *testing.T) {
	name, _ := ConvertComponentIdIntoFriendlyPathSearch("#/components/schemas/Widget Inner")
	assert.Equal(t, "Widget Inner", name)

	name, _ = ConvertComponentIdIntoFriendlyPathSearch("#/paths/~1pets~1{id}")
	assert.Equal(t, "/pets/{id}", name, "~1 is a JSON Pointer escape and must still decode")

	name, _ = ConvertComponentIdIntoFriendlyPathSearch("#/components/schemas/Widget")
	assert.Equal(t, "Widget", name)
}

// ConvertComponentIdIntoPath takes a DOT-separated JSONPath-style id, not a '#/' ref -- it splits on
// '.', so a slash-shaped input comes back whole. Same unescaping, so the same '+' guarantee applies.
func TestConvertComponentIdIntoPath_PlusIsLiteral(t *testing.T) {
	name, _ := ConvertComponentIdIntoPath("$.components.schemas.Widget+Inner")
	assert.Equal(t, "Widget+Inner", name, "'+' must survive unescaping here too")
}
