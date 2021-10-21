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
	repoURI := "https://github.com/philips-labs/slsa-provenance-action"
	builderID := repoURI + "/Attestations/GitHubHostedActions@v1"
	buildInvocationID := repoURI + "/actions/runs/123498765"
	recipeType := "https://github.com/Attestations/GitHubActionsWorkflow@v1"

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

	provenanceActionMaterial := []Item{
		{
			URI:    "git+https://github.com/philips-labs/slsa-provenance-action",
			Digest: DigestSet{"sha1": "c4f679f131dfb7f810fd411ac9475549d1c393df"},
		},
	}

	stmt = SLSAProvenanceStatement(
		WithSubject(make([]Subject, 1)),
		WithBuilder(builderID),
		WithRecipe(
			recipeType,
			"CI workflow",
			nil,
			nil,
			provenanceActionMaterial,
		),
	)
	r := stmt.Predicate.Recipe
	assert.Equal(predicateType, stmt.PredicateType)
	assert.Equal(statementType, stmt.Type)
	assert.Len(stmt.Subject, 1)
	assert.Equal(builderID, stmt.Predicate.Builder.ID)
	assert.Equal(recipeType, r.Type)
	assert.Equal("CI workflow", r.EntryPoint)
	assert.Nil(r.Arguments)
	assert.Equal(0, r.DefinedInMaterial)
	assert.Equal(provenanceActionMaterial, stmt.Predicate.Materials)
}
