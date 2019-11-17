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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBerglasSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBerglasSecretRead,

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

			"generation": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Generation of the object",
				Optional:    true,
			},

			//
			// Computed
			//
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Fully-qualified name of the Cloud KMS key",
				ForceNew:    true,
				Computed:    true,
			},

			"plaintext": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Plaintext contents",
				Computed:    true,
				Sensitive:   true,
			},

			"metageneration": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Metageneration of the object",
				Computed:    true,
			},
		},
	}
}

func dataSourceBerglasSecretRead(d *schema.ResourceData, meta interface{}) error {
	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	generation := d.Get("generation").(int)

	id := encodeId(bucket, name, int64(generation))
	d.SetId(id)
	return resourceBerglasSecretRead(d, meta)
}
