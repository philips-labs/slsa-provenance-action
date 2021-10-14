package intoto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSLSAProvenanceStatement(t *testing.T) {
	assert := assert.New(t)

	predicateType := "https://slsa.dev/provenance/v0.1"
	statementType := "https://in-toto.io/Statement/v0.1"
	repoURI := "https://github.com/slsa-provenance-action"
	builderID := repoURI + "/Attestations/GitHubHostedActions@v1"
	buildInvocationID := repoURI + "/actions/runs/123498765"

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

	stmt = SLSAProvenanceStatement(
		WithMetadata(buildInvocationID),
		WithBuilder(builderID),
	)
	m := stmt.Predicate.Metadata
	assert.Equal(predicateType, stmt.PredicateType)
	assert.Equal(statementType, stmt.Type)
	assert.Len(stmt.Subject, 0)
	assert.Equal(builderID, stmt.Predicate.Builder.ID)
	assert.Equal(buildInvocationID, m.BuildInvocationID)
	bft, err := time.Parse(time.RFC3339, m.BuildFinishedOn)
	assert.NoError(err)
	assert.WithinDuration(time.Now().UTC(), bft, 1200*time.Millisecond)
	assert.Equal(Completeness{Arguments: true, Environment: false, Materials: false}, stmt.Predicate.Metadata.Completeness)
	assert.False(m.Reproducible)
}
