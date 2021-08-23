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
	"fmt"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBerglasSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBerglasSecretCreate,
		ReadContext:   resourceBerglasSecretRead,
		UpdateContext: resourceBerglasSecretUpdate,
		DeleteContext: resourceBerglasSecretDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceBerglasSecretImport,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Description: "Name of the Cloud Storage bucket for the secret",
				ForceNew:    true,
				Required:    true,
			},

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the secret object in the bucket",
				ForceNew:    true,
				Required:    true,
			},

			"key": {
				Type:        schema.TypeString,
				Description: "Fully-qualified name of the Cloud KMS key",
				ForceNew:    true,
				Required:    true,
			},

			"plaintext": {
				Type:        schema.TypeString,
				Description: "Plaintext contents",
				Required:    true,
				Sensitive:   true,
			},

			//
			// Computed
			//
			"generation": {
				Type:        schema.TypeInt,
				Description: "Generation of the object",
				Computed:    true,
			},

			"metageneration": {
				Type:        schema.TypeInt,
				Description: "Metageneration of the object",
				Computed:    true,
			},
		},
	}
}

func resourceBerglasSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	client := config.Client()

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)
	key := d.Get("key").(string)
	plaintext := d.Get("plaintext").(string)

	secret, err := client.Create(ctx, &berglas.CreateRequest{
		Bucket:    bucket,
		Object:    name,
		Key:       key,
		Plaintext: []byte(plaintext),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create secret: %w", err))
	}

	id := encodeId(bucket, secret.Name, secret.Generation)
	d.SetId(id)

	if err := setMany(d, resourceFields{
		"generation":     secret.Generation,
		"metageneration": secret.Metageneration,
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update resource fields: %w", err))
	}

	return resourceBerglasSecretRead(ctx, d, meta)
}

func resourceBerglasSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	client := config.Client()

	bucket, object, generation, err := decodeId(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode id: %w", err))
	}

	secret, err := client.Read(ctx, &berglas.ReadRequest{
		Bucket:     bucket,
		Object:     object,
		Generation: generation,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read secret: %w", err))
	}

	if err := setMany(d, resourceFields{
		"bucket":         bucket,
		"name":           secret.Name,
		"key":            secret.KMSKey,
		"plaintext":      string(secret.Plaintext),
		"generation":     secret.Generation,
		"metageneration": secret.Metageneration,
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update resource fields: %w", err))
	}

	return nil
}

func resourceBerglasSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	client := config.Client()

	bucket, object, generation, err := decodeId(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode id: %w", err))
	}

	if d.HasChange("plaintext") {
		secret, err := client.Update(ctx, &berglas.UpdateRequest{
			Bucket:         bucket,
			Object:         object,
			Generation:     generation,
			Metageneration: int64(d.Get("metageneration").(int)),
			Key:            d.Get("key").(string),
			Plaintext:      []byte(d.Get("plaintext").(string)),
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update secret: %w", err))
		}

		id := encodeId(bucket, secret.Name, secret.Generation)
		d.SetId(id)

		if err := setMany(d, resourceFields{
			"generation":     secret.Generation,
			"metageneration": secret.Metageneration,
			"plaintext":      string(secret.Plaintext),
		}); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update resource fields: %w", err))
		}

		return resourceBerglasSecretRead(ctx, d, meta)
	}

	return nil
}

func resourceBerglasSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	client := config.Client()

	bucket, object, _, err := decodeId(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode id: %w", err))
	}

	if err := client.Delete(ctx, &berglas.DeleteRequest{
		Bucket: bucket,
		Object: object,
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete secret: %w", err))
	}

	d.SetId("")

	return nil
}

func resourceBerglasSecretImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	bucket, object, generation, err := decodeId(d.Id())
	if err != nil {
		return nil, fmt.Errorf("failed to decode id: %w", err)
	}

	if err := setMany(d, resourceFields{
		"bucket":     bucket,
		"name":       object,
		"generation": generation,
	}); err != nil {
		return nil, fmt.Errorf("failed to update resource fields: %w", err)
	}

	if diag := resourceBerglasSecretRead(ctx, d, meta); diag.HasError() {
		return nil, fmt.Errorf("failed to read secret")
	}

	return []*schema.ResourceData{d}, nil
}
