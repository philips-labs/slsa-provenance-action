package intoto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSLSAProvenanceStatement(t *testing.T) {
	assert := assert.New(t)

	repoURI := "https://github.com/philips-labs/slsa-provenance-action"
	builderID := repoURI + "/Attestations/GitHubHostedActions@v1"
	buildInvocationID := repoURI + "/actions/runs/123498765"
	buildType := "https://github.com/Attestations/GitHubActionsWorkflow@v1"

	stmt := SLSAProvenanceStatement()
	assert.Equal(SlsaPredicateType, stmt.PredicateType)
	assert.Equal(StatementType, stmt.Type)
	assert.Len(stmt.Subject, 0)

	stmt = SLSAProvenanceStatement(
		WithSubject(make([]Subject, 4)),
	)
	assert.Equal(SlsaPredicateType, stmt.PredicateType)
	assert.Equal(StatementType, stmt.Type)
	assert.Len(stmt.Subject, 4)

	stmt = SLSAProvenanceStatement(
		WithSubject(make([]Subject, 3)),
		WithBuilder(builderID),
	)
	assert.Equal(SlsaPredicateType, stmt.PredicateType)
	assert.Equal(StatementType, stmt.Type)
	assert.Len(stmt.Subject, 3)
	assert.Equal(builderID, stmt.Predicate.Builder.ID)

	stmt = SLSAProvenanceStatement(
		WithMetadata(buildInvocationID),
		WithBuilder(builderID),
	)
	m := stmt.Predicate.Metadata
	assert.Equal(SlsaPredicateType, stmt.PredicateType)
	assert.Equal(StatementType, stmt.Type)
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
		WithInvocation(
			buildType,
			"CI workflow",
			nil,
			nil,
			provenanceActionMaterial,
		),
	)
	i := stmt.Predicate.Invocation
	assert.Equal(SlsaPredicateType, stmt.PredicateType)
	assert.Equal(StatementType, stmt.Type)
	assert.Len(stmt.Subject, 1)
	assert.Equal(builderID, stmt.Predicate.Builder.ID)
	assert.Equal(buildType, stmt.Predicate.BuildType)
	assert.Equal("CI workflow", i.EntryPoint)
	assert.Nil(i.Arguments)
	assert.Equal(0, i.DefinedInMaterial)
	assert.Equal(provenanceActionMaterial, stmt.Predicate.Materials)
}
