package intoto

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	repoURI           = "https://github.com/philips-labs/slsa-provenance-action"
	builderID         = repoURI + "/Attestations/GitHubHostedActions@v1"
	buildInvocationID = repoURI + "/actions/runs/123498765"
	buildType         = "https://github.com/Attestations/GitHubActionsWorkflow@v1"
)

func TestSLSAProvenanceStatement(t *testing.T) {
	assert := assert.New(t)

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
	assert.Equal(Completeness{Parameters: true, Environment: false, Materials: false}, stmt.Predicate.Metadata.Completeness)
	assert.False(m.Reproducible)

	provenanceActionMaterial := []Item{
		{
			URI:    "git+https://github.com/philips-labs/slsa-provenance-action",
			Digest: DigestSet{"sha1": "c4f679f131dfb7f810fd411ac9475549d1c393df"},
		},
	}

	stmt = SLSAProvenanceStatement(
		WithSubject([]Subject{{Name: "salsa.txt", Digest: DigestSet{"sha256": "f8161d035cdf328c7bb124fce192cb90b603f34ca78d73e33b736b4f6bddf993"}}}),
		WithBuilder(builderID),
		WithMetadata("https://github.com/philips-labs/slsa-provenance-action/actions/runs/1303916967"),
		WithInvocation(
			buildType,
			"ci.yaml:build",
			nil,
			nil,
			provenanceActionMaterial,
		),
	)
	assertStatement(assert, stmt, builderID, buildType, provenanceActionMaterial, nil)
}

func TestSLSAProvenanceStatementJSON(t *testing.T) {
	assert := assert.New(t)

	materialJSON := `[
			{
				"uri": "git+https://github.com/philips-labs/slsa-provenance-action",
				"digest": {
					"sha1": "a3bc1c27230caa1cc3c27961f7e9cab43cd208dc"
				}
			}
		]`
	parametersJSON := `{
				"inputs": {
					"skip_integration": true
				}
			}`
	buildFinishedOn := time.Now().UTC().Format(time.RFC3339)

	var material []Item
	err := json.Unmarshal([]byte(materialJSON), &material)
	assert.NoError(err)

	jsonStatement := fmt.Sprintf(`{
	"_type": "https://in-toto.io/Statement/v0.1",
	"subject": [
		{
			"name": "salsa.txt",
			"digest": {
				"sha256": "f8161d035cdf328c7bb124fce192cb90b603f34ca78d73e33b736b4f6bddf993"
			}
		}
	],
	"predicateType": "https://slsa.dev/provenance/v0.2",
	"predicate": {
		"builder": {
			"id": "%s"
		},
		"buildType": "%s",
		"invocation": {
			"configSource": {
				"entryPoint": "ci.yaml:build",
				"uri": "git+https://github.com/philips-labs/slsa-provenance-action",
				"digest": {
					"sha1": "a3bc1c27230caa1cc3c27961f7e9cab43cd208dc"
				}
			},
			"parameters": %s,
			"environment": null
		},
		"metadata": {
			"buildInvocationId": "https://github.com/philips-labs/slsa-provenance-action/actions/runs/1303916967",
			"buildFinishedOn": "%s",
			"completeness": {
				"parameters": true,
				"environment": false,
				"materials": false
			},
			"reproducible": false
		},
		"materials": %s
	}
}`, builderID, buildType, parametersJSON, buildFinishedOn, materialJSON)

	var stmt Statement
	err = json.Unmarshal([]byte(jsonStatement), &stmt)
	assert.NoError(err)
	assertStatement(assert, &stmt, builderID, buildType, material, []byte(parametersJSON))

	newStmt := SLSAProvenanceStatement(
		WithSubject([]Subject{{Name: "salsa.txt", Digest: DigestSet{"sha256": "f8161d035cdf328c7bb124fce192cb90b603f34ca78d73e33b736b4f6bddf993"}}}),
		WithBuilder(builderID),
		WithMetadata("https://github.com/philips-labs/slsa-provenance-action/actions/runs/1303916967"),
		WithInvocation(buildType, "ci.yaml:build", nil, []byte(parametersJSON), material),
	)

	newStmtJSON, err := json.MarshalIndent(newStmt, "", "\t")
	assert.NoError(err)

	assert.Equal(jsonStatement, string(newStmtJSON))
}

func assertStatement(assert *assert.Assertions, stmt *Statement, builderID, buildType string, material []Item, parameters json.RawMessage) {
	i := stmt.Predicate.Invocation
	assert.Equal(SlsaPredicateType, stmt.PredicateType)
	assert.Equal(StatementType, stmt.Type)
	assert.Len(stmt.Subject, 1)
	assert.Equal(Subject{Name: "salsa.txt", Digest: DigestSet{"sha256": "f8161d035cdf328c7bb124fce192cb90b603f34ca78d73e33b736b4f6bddf993"}}, stmt.Subject[0])
	assert.Equal(builderID, stmt.Predicate.Builder.ID)
	assert.Equal(buildType, stmt.Predicate.BuildType)
	assertConfigSource(assert, i.ConfigSource, stmt.Predicate.Materials)
	assert.Nil(stmt.Predicate.BuildConfig)
	assert.Equal(parameters, i.Parameters)
	assert.Equal(material, stmt.Predicate.Materials)
	assertMetadata(assert, stmt.Predicate.Metadata)
}

func assertConfigSource(assert *assert.Assertions, cs ConfigSource, materials []Item) {
	assert.Equal("ci.yaml:build", cs.EntryPoint)
	assert.Equal(materials[0].URI, cs.URI)
	assert.Equal(materials[0].Digest, cs.Digest)
}

func assertMetadata(assert *assert.Assertions, md Metadata) {
	assert.Equal("https://github.com/philips-labs/slsa-provenance-action/actions/runs/1303916967", md.BuildInvocationID)
	bft, err := time.Parse(time.RFC3339, md.BuildFinishedOn)
	assert.NoError(err)
	assert.WithinDuration(time.Now().UTC(), bft, 1200*time.Millisecond)
	assert.True(md.Completeness.Parameters)
	assert.False(md.Completeness.Materials)
	assert.False(md.Completeness.Environment)
	assert.False(md.Reproducible)
}
