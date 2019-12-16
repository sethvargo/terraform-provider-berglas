// Copyright 2019 Seth Vargo, Katie McLaughlin
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
	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBerglasIam() *schema.Resource {

	return &schema.Resource{
		Create: resourceBerglasIamCreate,
		Read:   resourceBerglasIamRead,
		Update: resourceBerglasIamUpdate,
		Delete: resourceBerglasIamDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBerglasIamImport,
		},

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the Cloud Storage bucket for the secret",
				ForceNew:    true,
				Required:    true,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the secret object in the bucket",
				ForceNew:    true,
				Required:    true,
			},

			"members": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of members",
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceBerglasIamCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	members := d.Get("members").(string)

	if err := client.Grant(ctx, &berglas.GrantRequest{
		Project: project,
		Name:    name,
		Members: []string{serviceAccount},
	}); err != nil {
		t.Fatal(err)
	}

	return nil

}
func resourceBerglasIamRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	// TODO
	return nil

}
func resourceBerglasIamUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	// TODO
	return nil
}
func resourceBerglasIamDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	if err := client.Revoke(ctx, &berglas.RevokeRequest{
		Project: project,
		Name:    name,
		Members: []string{serviceAccount},
	}); err != nil {
		t.Fatal(err)
	}
	return nil
}

func resourceBerglasIamImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// TODO
}
