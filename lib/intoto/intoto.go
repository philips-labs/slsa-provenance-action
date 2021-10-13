package intoto

import (
	"encoding/json"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
)

// Envelope wraps an in-toto statement to be able to attach signatures to the Statement
type Envelope struct {
	PayloadType string        `json:"payloadType"`
	Payload     string        `json:"payload"`
	Signatures  []interface{} `json:"signatures"`
}

// Statement The Statement is the middle layer of the attestation, binding it to a particular subject and unambiguously identifying the types of the predicate.
type Statement struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     `json:"predicate"`
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
	Builder   `json:"builder"`
	Metadata  `json:"metadata"`
	Recipe    `json:"recipe"`
	Materials []Item `json:"materials"`
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
	Completeness      `json:"completeness"`
	Reproducible      bool `json:"reproducible"`
	// BuildStartedOn not defined as it's not available from a GitHub Action.
	BuildFinishedOn string `json:"buildFinishedOn"`
}

// Recipe Identifies the configuration used for the build. When combined with materials, this SHOULD fully describe the build, such that re-running this recipe results in bit-for-bit identical output (if the build is reproducible).
type Recipe struct {
	Type              string             `json:"type"`
	DefinedInMaterial int                `json:"definedInMaterial"`
	EntryPoint        string             `json:"entryPoint"`
	Arguments         json.RawMessage    `json:"arguments"`
	Environment       *github.AnyContext `json:"environment"`
}

// Completeness Indicates that the builder claims certain fields in this message to be complete.
type Completeness struct {
	Arguments   bool `json:"arguments"`
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
