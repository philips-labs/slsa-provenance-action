package provenance

import "encoding/json"

type Envelope struct {
	PayloadType string        `json:"payloadType"`
	Payload     string        `json:"payload"`
	Signatures  []interface{} `json:"signatures"`
}

type Statement struct {
	Type          string    `json:"_type"`
	Subject       []Subject `json:"subject"`
	PredicateType string    `json:"predicateType"`
	Predicate     `json:"predicate"`
}

type Subject struct {
	Name   string    `json:"name"`
	Digest DigestSet `json:"digest"`
}

type Predicate struct {
	Builder   `json:"builder"`
	Metadata  `json:"metadata"`
	Recipe    `json:"recipe"`
	Materials []Item `json:"materials"`
}

type Builder struct {
	Id string `json:"id"`
}

type Metadata struct {
	BuildInvocationId string `json:"buildInvocationId"`
	Completeness      `json:"completeness"`
	Reproducible      bool `json:"reproducible"`
	// BuildStartedOn not defined as it's not available from a GitHub Action.
	BuildFinishedOn string `json:"buildFinishedOn"`
}

type Recipe struct {
	Type              string          `json:"type"`
	DefinedInMaterial int             `json:"definedInMaterial"`
	EntryPoint        string          `json:"entryPoint"`
	Arguments         json.RawMessage `json:"arguments"`
	Environment       *AnyContext     `json:"environment"`
}

type Completeness struct {
	Arguments   bool `json:"arguments"`
	Environment bool `json:"environment"`
	Materials   bool `json:"materials"`
}

type DigestSet map[string]string

type Item struct {
	URI    string    `json:"uri"`
	Digest DigestSet `json:"digest"`
}

type AnyContext struct {
	GitHubContext `json:"github"`
	RunnerContext `json:"runner"`
}

type GitHubContext struct {
	Action          string          `json:"action"`
	ActionPath      string          `json:"action_path"`
	Actor           string          `json:"actor"`
	BaseRef         string          `json:"base_ref"`
	Event           json.RawMessage `json:"event"`
	EventName       string          `json:"event_name"`
	EventPath       string          `json:"event_path"`
	HeadRef         string          `json:"head_ref"`
	Job             string          `json:"job"`
	Ref             string          `json:"ref"`
	Repository      string          `json:"repository"`
	RepositoryOwner string          `json:"repository_owner"`
	RunId           string          `json:"run_id"`
	RunNumber       string          `json:"run_number"`
	SHA             string          `json:"sha"`
	Token           string          `json:"token,omitempty"`
	Workflow        string          `json:"workflow"`
	Workspace       string          `json:"workspace"`
}

type RunnerContext struct {
	OS        string `json:"os"`
	Temp      string `json:"temp"`
	ToolCache string `json:"tool_cache"`
}

// See https://docs.github.com/en/actions/reference/events-that-trigger-workflows
// The only Event with dynamically-provided input is workflow_dispatch which
// exposes the user params at the key "input."
type AnyEvent struct {
	Inputs json.RawMessage `json:"inputs"`
}
