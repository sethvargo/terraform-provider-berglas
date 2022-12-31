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

package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sethvargo/terraform-provider-berglas/internal/pathorcontents"
)

const (
	cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
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
		p.ConfigureContextFunc = providerConfigure(version, p)

		return p
	}
}

// providerConfigure configures the provider
func providerConfigure(version string, p *schema.Provider) schema.ConfigureContextFunc {
	return func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		accessToken := d.Get("access_token").(string)
		credentials := d.Get("credentials").(string)

		// Note that we explicitly use context.Background() instead of the provided
		// context because we want to give the client a chance to finish before
		// cleanup.
		tokenSource, err := tokenSource(context.Background(), accessToken, credentials)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to configure provider: %w", err))
		}

		client, err := berglas.New(context.Background(), option.WithTokenSource(tokenSource))
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to setup berglas: %w", err))
		}

		config := &config{
			client: client,
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
			return nil, fmt.Errorf("failed to load access token: %w", err)
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
			return nil, fmt.Errorf("failed to load credentials: %w", err)
		}

		creds, err := google.CredentialsFromJSON(ctx, []byte(contents), cloudPlatformScope)
		if err != nil {
			return nil, fmt.Errorf("failed to parse credentials: %w", err)
		}

		return creds.TokenSource, nil
	}

	// Fallback to default credentials
	log.Printf("[INFO] authenticating via default credentials")
	source, err := google.DefaultTokenSource(ctx, cloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("failed to get default credentials: %w", err)
	}
	return source, nil
}
