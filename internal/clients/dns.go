/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/pkg/terraform"

	"github.com/dana-team/provider-dns/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal dns credentials as JSON"

	// general parameters
	keyRFC    = "rfc"
	keyServer = "server"
	update    = "update"

	// gss-tsig (RFC 3645) parameters
	gsstsigRFC  = "3645"
	gssapi      = "gssapi"
	keyTab      = "keytab"
	keyPassword = "password"
	keyRealm    = "realm"
	keyUsername = "username"

	// secret key based transaction authentication (RFC 2845) parameters
	keyBasedTransactionRFC  = "2845"
	transactionKeyAlgorithm = "key_algorithm"
	transcationKeyName      = "key_name"
	transactionKeySecret    = "key_secret"
	keyPort                 = "port"
	keyRetries              = "retries"
	keyTimeout              = "timeout"
	keyTransport            = "transport"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// Set credentials in Terraform provider configuration.
		authConfig := buildAuthConfig(creds)
		ps.Configuration = map[string]interface{}{
			update: authConfig,
		}

		return ps, nil
	}
}

// buildAuthConfig builds the auth configuration for the provider.
func buildAuthConfig(creds map[string]string) map[string]any {
	authConfig := map[string]any{}

	if server, ok := creds[keyServer]; ok {
		authConfig[keyServer] = server
	}

	if rfc, ok := creds[keyRFC]; ok {
		if rfc == gsstsigRFC {
			authConfig[gssapi] = buildGSSTSIGAuthConfig(creds)
		} else if rfc == keyBasedTransactionRFC {
			secredBasedTransactionAuthConfig := buildSecretBasedTransactionAuthConfig(creds)
			for k, v := range secredBasedTransactionAuthConfig {
				authConfig[k] = v
			}
		}
	}

	optionalConfig := buildOptionalConfig(creds)
	for k, v := range optionalConfig {
		authConfig[k] = v
	}

	return authConfig

}

// buildGSSTSIGAuthConfig builds the configuration for GSS-TSIG authentication (RFC 3645).
func buildGSSTSIGAuthConfig(creds map[string]string) map[string]string {
	config := make(map[string]string)

	if realm, ok := creds[keyRealm]; ok {
		config[keyRealm] = realm
	}

	if username, ok := creds[keyUsername]; ok {
		config[keyUsername] = username
	}

	if password, ok := creds[keyPassword]; ok {
		config[keyPassword] = password
	}

	if keytab, ok := creds[keyTab]; ok {
		config[keyTab] = keytab
	}

	return config
}

// // buildGSSTSIGAuthConfig builds the configuration for GSS-TSIG authentication (RFC 2845).
func buildSecretBasedTransactionAuthConfig(creds map[string]string) map[string]string {
	config := make(map[string]string)

	if keyName, ok := creds[transcationKeyName]; ok {
		config[transcationKeyName] = keyName
	}

	if keyAlgorithm, ok := creds[transactionKeyAlgorithm]; ok {
		config[transactionKeyAlgorithm] = keyAlgorithm
	}

	if keySecret, ok := creds[transactionKeySecret]; ok {
		config[transactionKeyAlgorithm] = keySecret
	}

	return config
}

// buildOptionalConfig builds the optional configuration for the provider.
func buildOptionalConfig(creds map[string]string) map[string]string {
	config := make(map[string]string)

	if port, ok := creds[keyPort]; ok {
		config[keyPort] = port
	}

	if retries, ok := creds[keyRetries]; ok {
		config[keyRetries] = retries
	}

	if timeout, ok := creds[keyTimeout]; ok {
		config[keyTimeout] = timeout
	}

	if transport, ok := creds[keyTransport]; ok {
		config[keyTransport] = transport
	}

	return config
}
