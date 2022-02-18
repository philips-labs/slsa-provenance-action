package intoto

import (
	"context"
	"encoding/json"
	"time"
)

const (
	// SlsaPredicateType the predicate type for SLSA intoto statements
	SlsaPredicateType = "https://slsa.dev/provenance/v0.2"
	// StatementType the type of the intoto statement
	StatementType = "https://in-toto.io/Statement/v0.1"
)

// Provenancer generates provenance statements for given artifacts
type Provenancer interface {
	GenerateProvenanceStatement(ctx context.Context, subjecter Subjecter) (*Statement, error)
	PersistProvenanceStatement(ctx context.Context, stmt *Statement, path string) error
}

// Envelope wraps an in-toto statement to be able to attach signatures to the Statement
type Envelope struct {
	PayloadType string        `json:"payloadType"`
	Payload     string        `json:"payload"`
	Signatures  []interface{} `json:"signatures"`
}

// SLSAProvenanceStatement builds a in-toto statement with predicate type https://slsa.dev/provenance/v0.1
func SLSAProvenanceStatement(opts ...StatementOption) *Statement {
	stmt := &Statement{PredicateType: SlsaPredicateType, Type: StatementType}
	for _, opt := range opts {
		opt(stmt)
	}
	return stmt
}

// StatementOption option flag to build the Statement
type StatementOption func(*Statement)

// WithSubject sets the Statement subject to the provided value
func WithSubject(s []Subject) StatementOption {
	return func(st *Statement) {
		st.Subject = s
	}
}

// WithBuilder sets the Statement builder with the given ID
func WithBuilder(id string) StatementOption {
	return func(st *Statement) {
		st.Predicate.Builder = Builder{ID: id}
	}
}

// WithMetadata sets the Predicate Metadata using the buildInvocationID and the current time
func WithMetadata(buildInvocationID string) StatementOption {
	return func(s *Statement) {
		s.Predicate.Metadata = Metadata{
			Completeness: Completeness{
				Parameters:  true,
				Environment: false,
				Materials:   false,
			},
			Reproducible:      false,
			BuildInvocationID: buildInvocationID,
			BuildFinishedOn:   time.Now().UTC().Format(time.RFC3339),
		}
	}
}

// WithInvocation sets the Predicate Invocation and Materials
func WithInvocation(buildType, entryPoint string, environment json.RawMessage, parameters json.RawMessage, materials []Item) StatementOption {
	return func(s *Statement) {
		s.Predicate.BuildType = buildType
		s.Predicate.Invocation = Invocation{
			ConfigSource: ConfigSource{
				EntryPoint: entryPoint,
				URI:        materials[0].URI,
				Digest:     materials[0].Digest,
			},
			Parameters:  parameters,
			Environment: environment,
		}
		s.Predicate.Materials = append(s.Predicate.Materials, materials...)
	}
}

// Statement The Statement is the middle layer of the attestation, binding it to a particular subject and unambiguously identifying the types of the predicate.
type Statement struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     Predicate `json:"predicate"`
}

// Subject The software artifacts that the attestation applies to.
type Subject struct {
	Name   string    `json:"name"`
	Digest DigestSet `json:"digest"`
}

// Predicate This predicate follows the in-toto attestation parsing rules.
//
// https://github.com/in-toto/attestation/blob/main/spec/README.md#parsing-rules
//
// The Predicate is the innermost layer of the attestation, containing arbitrary metadata about the Statement's subject.
//
// A predicate has a required predicateType (TypeURI) identifying what the predicate means, plus an optional predicate (object) containing additional, type-dependent parameters.
type Predicate struct {
	Builder     `json:"builder"`
	BuildType   string `json:"buildType"`
	Invocation  `json:"invocation"`
	BuildConfig *BuildConfig `json:"build_config,omitempty"`
	Metadata    `json:"metadata,omitempty"`
	Materials   []Item `json:"materials"`
}

// BuildConfig Lists the steps in the build.
// If invocation.sourceConfig is not available, buildConfig can be used to verify information about the build.
type BuildConfig struct {
}

// Builder Identifies the entity that executed the recipe, which is trusted to have correctly performed the operation and populated this provenance.
// The identity MUST reflect the trust base that consumers care about. How detailed to be is a judgement call. For example, GitHub Actions supports both GitHub-hosted runners and self-hosted runners. The GitHub-hosted runner might be a single identity because, it's all GitHub from the consumer's perspective. Meanwhile, each self-hosted runner might have its own identity because not all runners are trusted by all consumers.
//
// Consumers MUST accept only specific (signer, builder) pairs. For example, the "GitHub" can sign provenance for the "GitHub Actions" builder, and "Google" can sign provenance for the "Google Cloud Build" builder, but "GitHub" cannot sign for the "Google Cloud Build" builder.
//
// Design rationale: The builder is distinct from the signer because one signer may generate attestations for more than one builder, as in the GitHub Actions example above. The field is required, even if it is implicit from the signer, to aid readability and debugging. It is an object to allow additional fields in the future, in case one URI is not sufficient.
type Builder struct {
	ID string `json:"id"`
}

// Metadata Other properties of the build.
type Metadata struct {
	BuildInvocationID string `json:"buildInvocationId"`
	// BuildStartedOn not defined as it's not available from a GitHub Action.
	BuildFinishedOn string `json:"buildFinishedOn"`
	Completeness    `json:"completeness"`
	Reproducible    bool `json:"reproducible"`
}

// Invocation Identifies the configuration used for the build. When combined with materials, this SHOULD fully describe the build, such that re-running this recipe results in bit-for-bit identical output (if the build is reproducible).
type Invocation struct {
	ConfigSource ConfigSource    `json:"configSource"`
	Parameters   json.RawMessage `json:"parameters"`
	Environment  json.RawMessage `json:"environment"`
}

// ConfigSource Describes where the config file that kicked off the build came from.
// This is effectively a pointer to the source where buildConfig came from.
type ConfigSource struct {
	EntryPoint string    `json:"entryPoint"`
	URI        string    `json:"uri,omitempty"`
	Digest     DigestSet `json:"digest,omitempty"`
}

// Completeness Indicates that the builder claims certain fields in this message to be complete.
type Completeness struct {
	Parameters  bool `json:"parameters"`
	Environment bool `json:"environment"`
	Materials   bool `json:"materials"`
}

// DigestSet Collection of cryptographic digests for the contents of this artifact.
type DigestSet map[string]string

// Item The material used as input for producing the output artifact (subject).
type Item struct {
	URI    string    `json:"uri"`
	Digest DigestSet `json:"digest"`
}
