package cli_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

const (
	provenanceData = `
	{
		"_type": "https://in-toto.io/Statement/v0.1",
		"subject": [
		  {
			"name": "slsa-provenance",
			"digest": {
			  "sha256": "d42c69de30b2e9ad94d524cc4c8658a8c4d56440837ef64adcde88deb72b5ff0"
			}
		  }
		],
		"predicateType": "https://slsa.dev/provenance/v0.1",
		"predicate": {
		  "builder": {
			"id": "https://github.com/philips-labs/slsa-provenance-action/Attestations/SelfHostedActions@v1"
		  },
		  "metadata": {
			"buildInvocationId": "https://github.com/philips-labs/slsa-provenance-action/actions/runs/1332651620",
			"completeness": {
			  "arguments": true,
			  "environment": false,
			  "materials": false
			},
			"reproducible": false,
			"buildFinishedOn": "2021-11-04T15:55:21Z"
		  },
		  "recipe": {
			"type": "https://github.com/Attestations/GitHubActionsWorkflow@v1",
			"definedInMaterial": 0,
			"entryPoint": "Integration test file provenance",
			"arguments": null,
			"environment": null
		  },
		  "materials": [
			{
			  "uri": "git+https://github.com/philips-labs/slsa-provenance-action",
			  "digest": {
				"sha1": "c4f679f131dfb7f810fd411ac9475549d1c393df"
			  }
			},
			{
			  "uri": "pkg:deb/debian/stunnel4@5.50-3?arch=amd64",
			  "digest": {
				"sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa"
			  }
			}
		  ]
		}
	  }
	`
)

func TestSignCliOptions(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	provenanceFile := path.Join(rootDir, "bin/unittest.provenance")
	keyFile := path.Join(rootDir, "bin/unittest.key.file")
	key := "0305334e381af78f141cb666f6199f57bc3495335a256a95bd2a55bf546663f6"
	outFile := path.Join(rootDir, "bin/provenance.signed.json")

	ioutil.WriteFile(provenanceFile, []byte(provenanceData), 0644)
	defer (func() {
		_ = os.Remove(provenanceFile)
	})()

	testCases := []struct {
		name      string
		err       error
		arguments []string
	}{
		{
			name:      "no flags",
			err:       errors.New("Neither key nor key-path specified"),
			arguments: make([]string, 0),
		}, {
			name: "both --key-path and --key set",
			err:  errors.New("Both key and key-path specified"),
			arguments: []string{
				"--key",
				key,
				"--key-path",
				keyFile,
			},
		}, {
			name: "no --provenance-path set",
			err:  cli.RequiredFlagError("--provenance-path"),
			arguments: []string{
				"--provenance-path",
				"",
				"--key",
				key,
			},
		}, {
			name: "no --output-path set",
			err:  cli.RequiredFlagError("--output-path"),
			arguments: []string{
				"--key",
				key,
				"--output-path",
				"",
			},
		}, {
			name: "no provenance file",
			err:  errors.New("Error reading provenance file provenance.json: open provenance.json: no such file or directory"),
			arguments: []string{
				"--key",
				key,
			},
		}, {
			name: "bad key (wrong size)",
			err:  errors.New("Unable to decode key: encoding/hex: invalid byte: U+0058 'X'"),
			arguments: []string{
				"--provenance-path",
				provenanceFile,
				"--key",
				key + "XX",
			},
		}, {
			name: "bad key (not hex)",
			err:  errors.New("Decoded key has wrong size, expected 32 bytes, got 33"),
			arguments: []string{
				"--provenance-path",
				provenanceFile,
				"--key",
				key + "AB",
			},
		}, {
			name: "no keyfile",
			err:  errors.New("Error reading key file nope: open nope: no such file or directory"),
			arguments: []string{
				"--provenance-path",
				provenanceFile,
				"--key-path",
				"nope",
			},
		}, {
			name: "successful sign with key",
			err:  nil,
			arguments: []string{
				"--provenance-path",
				provenanceFile,
				"--key",
				key,
				"--output-path",
				outFile,
			},
		}, {
			name: "successful sign with key file",
			err:  nil,
			arguments: []string{
				"--provenance-path",
				provenanceFile,
				"--key-path",
				keyFile,
				"--output-path",
				outFile,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

			assert.NoError(ioutil.WriteFile(keyFile, []byte(key), 0644))

			defer func() {
				_ = os.Remove(keyFile)
				_ = os.Remove(outFile)
			}()

			cmd := cli.Sign()
			_, err := executeCommand(cmd, tc.arguments...)
			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
				if assert.FileExists(outFile) {
					content, err := os.ReadFile(outFile)
					assert.NoError(err)
					assert.Greater(len(content), 1)
					var envelope intoto.Envelope
					assert.NoError(json.Unmarshal(content, &envelope))
				}
			}
		})
	}
}

func TestSignSignature(t *testing.T) {
	// TODO check if we indeed generate a good signature
}

func BenchmarkSign(b *testing.B) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	provenanceFile := path.Join(rootDir, "bin/unittest.provenance")
	keyFile := path.Join(rootDir, "bin/unittest.key.file")
	key := "0305334e381af78f141cb666f6199f57bc3495335a256a95bd2a55bf546663f6"
	outFile := path.Join(rootDir, "bin/provenance.signed")

	ioutil.WriteFile(provenanceFile, []byte(provenanceData), 0644)
	defer (func() {
		_ = os.Remove(provenanceFile)
	})()

	ioutil.WriteFile(keyFile, []byte(key), 0644)

	defer func() {
		_ = os.Remove(keyFile)
		_ = os.Remove(outFile)
	}()

	b.Run("using commandline key", func(b *testing.B) {
		cmd := cli.Sign()
		for i := 0; i < b.N; i++ {
			executeCommand(cmd,
				"-provenance_path",
				provenanceFile,
				"-key",
				key,
				"-output_path",
				outFile,
			)
		}
	})

	b.Run("using key file", func(b *testing.B) {
		cmd := cli.Sign()
		for i := 0; i < b.N; i++ {
			executeCommand(cmd,
				"-provenance_path",
				provenanceFile,
				"-key-file",
				keyFile,
				"-output_path",
				outFile,
			)
		}
	})
}