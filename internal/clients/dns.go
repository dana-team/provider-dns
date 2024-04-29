/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/dana-team/provider-dns/apis/v1beta1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-dns/xpprovider"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/pkg/terraform"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal dns credentials as JSON"

	// general parameters
	keyRFC       = "rfc"
	keyServer    = "server"
	update       = "update"
	keyPort      = "port"
	keyRetries   = "retries"
	keyTimeout   = "timeout"
	keyTransport = "transport"

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
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(tfProvider *schema.Provider) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{}

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

		// set credentials in Terraform provider configuration.
		ps.Configuration = map[string]any{}
		authConfig := buildAuthConfig(creds)
		ps.Configuration[update] = []map[string]any{authConfig}

		return ps, errors.Wrap(configureNoForkDNSClient(ctx, &ps, *tfProvider), "failed to configure the Terraform DNS provider meta")
	}
}

// configureNoForkDNSClient populates the supplied *terraform.Setup with
// Terraform Plugin SDK style DNS client (Meta) and Terraform Plugin Framework
// style FrameworkProvider
func configureNoForkDNSClient(ctx context.Context, ps *terraform.Setup, p schema.Provider) error {
	diag := p.Configure(context.WithoutCancel(ctx), &tfsdk.ResourceConfig{
		Config: ps.Configuration,
	})
	if diag != nil && diag.HasError() {
		return errors.Errorf("failed to configure the provider: %v", diag)
	}
	ps.Meta = p.Meta()
	fwProvider, _ := xpprovider.GetProvider(ctx)
	ps.FrameworkProvider = fwProvider

	return nil
}

// buildAuthConfig builds the auth configuration for the provider.
func buildAuthConfig(creds map[string]string) map[string]any {
	config := map[string]any{}

	if server, ok := creds[keyServer]; ok {
		config[keyServer] = server
	}

	if rfc, ok := creds[keyRFC]; ok {
		if rfc == gsstsigRFC {
			authConfig := buildGSSTSIGAuthConfig(creds)
			config[gssapi] = []map[string]any{authConfig}
		} else if rfc == keyBasedTransactionRFC {
			secretBasedTransactionAuthConfig := buildSecretBasedTransactionAuthConfig(creds)
			mergeMaps(config, secretBasedTransactionAuthConfig)
		}
	}

	optionalConfig := buildOptionalConfig(creds)
	mergeMaps(config, optionalConfig)

	return config

}

// buildGSSTSIGAuthConfig builds the configuration for GSS-TSIG authentication (RFC 3645).
func buildGSSTSIGAuthConfig(creds map[string]string) map[string]any {
	config := make(map[string]any)

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
func buildSecretBasedTransactionAuthConfig(creds map[string]string) map[string]any {
	config := make(map[string]any)

	if keyName, ok := creds[transcationKeyName]; ok {
		config[transcationKeyName] = keyName
	}

	if keyAlgorithm, ok := creds[transactionKeyAlgorithm]; ok {
		config[transactionKeyAlgorithm] = keyAlgorithm
	}

	if keySecret, ok := creds[transactionKeySecret]; ok {
		config[transactionKeySecret] = keySecret
	}

	return config
}

// buildOptionalConfig builds the optional configuration for the provider.
func buildOptionalConfig(creds map[string]string) map[string]any {
	config := make(map[string]any)

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

// mergeMaps takes all the keys in B and inserts them into A.
func mergeMaps(a, b map[string]any) {
	for k, v := range b {
		a[k] = v
	}
}
