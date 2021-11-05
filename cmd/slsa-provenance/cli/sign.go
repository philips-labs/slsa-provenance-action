package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/dsse"
	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

func Sign() *cobra.Command {
	o := &options.SignOptions{}

	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Generate a signature on provenance",
		RunE: func(cmd *cobra.Command, args []string) error {
			privKey, err := o.GetKey()
			if err != nil {
				return err
			}
			outputPath, err := o.GetOutputPath()
			if err != nil {
				return err
			}
			provenancePath, err := o.GetProvenancePath()
			if err != nil {
				return err
			}

			content, err := ioutil.ReadFile(provenancePath)
			if err != nil {
				return errors.Wrapf(err, "Error reading provenance file %s", provenancePath)
			}
			var provenance intoto.Statement
			if err = json.Unmarshal(content, &provenance); err != nil {
				return errors.Wrapf(err, "Invalid JSON in provenance file %s", provenancePath)
			}

			var signer *signature.ED25519Signer
			if signer, err = signature.LoadED25519Signer(privKey); err != nil {
				return errors.Wrap(err, "Can not initiate signer")
			}

			wrappedSigner := dsse.WrapSigner(signer, intoto.InTotoPayloadType)

			// TODO canonicalize the provenance
			var toSign []byte
			if toSign, err = json.Marshal(provenance); err != nil {
				// Should be impossible, but hey!
				return errors.Wrap(err, "Could not marshal provenance for signing")
			}
			reader := bytes.NewReader(toSign)

			var signedMessage []byte
			if signedMessage, err = wrappedSigner.SignMessage(reader); err != nil {
				return errors.Wrap(err, "Could not sign data")
			}

			var envelope intoto.Envelope
			if err = json.Unmarshal(signedMessage, &envelope); err != nil {
				return errors.Wrap(err, "Signer produced date we can not use")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saving signed provenance to %s\n", outputPath)

			var payload []byte
			if payload, err = json.MarshalIndent(envelope, "", "  "); err != nil {
				return errors.Wrap(err, "Failed to marshal signed provenance")
			}

			if err = os.WriteFile(outputPath, payload, 0644); err != nil {
				return errors.Wrapf(err, "Failed to write signed provenance to %s", outputPath)
			}
			return nil
		},
	}
	o.AddFlags(cmd)
	return cmd
}
