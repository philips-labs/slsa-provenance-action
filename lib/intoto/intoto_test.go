package intoto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSLSAProvenanceStatement(t *testing.T) {
	assert := assert.New(t)

	predicateType := "https://slsa.dev/provenance/v0.1"
	statementType := "https://in-toto.io/Statement/v0.1"
	builderID := "https://github.com/slsa-provenance-action/Attestations/GitHubHostedActions@v1"

	stmt := SLSAProvenanceStatement()
	assert.Equal(predicateType, stmt.PredicateType)
	assert.Equal(statementType, stmt.Type)
	assert.Len(stmt.Subject, 0)

	stmt = SLSAProvenanceStatement(
		WithSubject(make([]Subject, 4)),
	)
	assert.Equal(predicateType, stmt.PredicateType)
	assert.Equal(statementType, stmt.Type)
	assert.Len(stmt.Subject, 4)

	stmt = SLSAProvenanceStatement(
		WithSubject(make([]Subject, 3)),
		WithBuilder(builderID),
	)
	assert.Equal(predicateType, stmt.PredicateType)
	assert.Equal(statementType, stmt.Type)
	assert.Len(stmt.Subject, 3)
	assert.Equal(builderID, stmt.Predicate.Builder.ID)
}
