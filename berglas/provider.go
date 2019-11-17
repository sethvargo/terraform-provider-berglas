// Copyright 2019 Seth Vargo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package berglas

import (
	"context"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)

// Provider returns the actual provider instance.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_APPLICATION_CREDENTIALS",
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				Description: strings.TrimSpace(`
JSON credentials with which to authenticate against the API. This can be set to
the raw credential contents or it can be set to a file path on disk which
contains the file contents.
`),
				ConflictsWith: []string{"access_token"},
			},

			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_OAUTH_ACCESS_TOKEN",
				}, nil),
				Description: strings.TrimSpace(`
OAuth2 access token to use for communicating with Google APIs.
`),
				ConflictsWith: []string{"credentials"},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"berglas_secret": dataSourceBerglasSecret(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"berglas_secret": resourceBerglasSecret(),
		},
	}

	// Meta, but we have to pass the provider into itself
	provider.ConfigureFunc = providerConfigure(provider)

	return provider
}

// providerConfigure configures the provider
func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		ctx := p.StopContext()

		accessToken := d.Get("access_token").(string)
		credentials := d.Get("credentials").(string)

		tokenSource, err := tokenSource(ctx, accessToken, credentials)
		if err != nil {
			return nil, errors.Wrap(err, "failed to configure provider")
		}

		client, err := berglas.New(ctx, option.WithTokenSource(tokenSource))
		if err != nil {
			return nil, errors.Wrap(err, "failed to setup berglas")
		}

		config := &config{
			client: client,
			ctx:    ctx,
		}

		// Update the stored context if the provider is reset.
		p.MetaReset = func() error {
			config.lock.Lock()
			defer config.lock.Unlock()

			config.ctx = p.StopContext()
			return nil
		}

		return config, nil
	}
}

// tokenSource returns the best token source for the given environment.
func tokenSource(ctx context.Context, accessToken, credentials string) (oauth2.TokenSource, error) {
	// Try access token first
	if accessToken != "" {
		log.Printf("[INFO] authenticating via access_token")

		contents, _, err := pathorcontents.Read(accessToken)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load access token")
		}

		return oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: contents,
		}), nil
	}

	// Then credentials
	if credentials != "" {
		log.Printf("[INFO] authenticating via credentials")

		contents, _, err := pathorcontents.Read(credentials)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load credentials")
		}

		creds, err := google.CredentialsFromJSON(ctx, []byte(contents), cloudPlatformScope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse credentials")
		}

		return creds.TokenSource, nil
	}

	// Fallback to default credentials
	log.Printf("[INFO] authenticating via default credentials")
	source, err := google.DefaultTokenSource(ctx, cloudPlatformScope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default credentials")
	}
	return source, nil
}
