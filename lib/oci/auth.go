package oci

import (
	"context"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

// WithDefaultClientOptions sets some sane default options for crane to authenticate
// private registries
func WithDefaultClientOptions(ctx context.Context, k8sKeychain bool) []crane.Option {
	opts := []crane.Option{
		crane.WithContext(ctx),
	}

	if k8sKeychain {
		kc := authn.NewMultiKeychain(
			authn.DefaultKeychain,
			google.Keychain,
			authn.NewKeychainFromHelper(ecr.ECRHelper{ClientFactory: api.DefaultClientFactory{}}),
			authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper()),
		)
		opts = append(opts, crane.WithAuthFromKeychain(kc))
	} else {
		opts = append(opts, crane.WithAuthFromKeychain(authn.DefaultKeychain))
	}

	return opts
}
