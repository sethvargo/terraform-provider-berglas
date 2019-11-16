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
	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBerglasSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceBerglasSecretCreate,
		Read:   resourceBerglasSecretRead,
		Update: resourceBerglasSecretUpdate,
		Delete: resourceBerglasSecretDelete,

		Importer: &schema.ResourceImporter{
			State: resourceBerglasSecretImport,
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

			"key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Fully-qualified name of the Cloud KMS key",
				ForceNew:    true,
				Required:    true,
			},

			"plaintext": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Plaintext contents",
				Required:    true,
				Sensitive:   true,
			},

			//
			// Computed
			//
			"generation": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Generation of the object",
				Computed:    true,
			},

			"metageneration": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Metageneration of the object",
				Computed:    true,
			},
		},
	}
}

func resourceBerglasSecretCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

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
		return err
	}

	id := encodeId(bucket, secret.Name, secret.Generation)
	d.SetId(id)

	if err := setMany(d, resourceFields{
		"generation":     secret.Generation,
		"metageneration": secret.Metageneration,
	}); err != nil {
		return err
	}

	return resourceBerglasSecretRead(d, meta)
}

func resourceBerglasSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	bucket, object, generation, err := decodeId(d.Id())
	if err != nil {
		return err
	}

	secret, err := client.Read(ctx, &berglas.ReadRequest{
		Bucket:     bucket,
		Object:     object,
		Generation: generation,
	})
	if err != nil {
		return err
	}

	if err := setMany(d, resourceFields{
		"bucket":         bucket,
		"name":           secret.Name,
		"key":            secret.KMSKey,
		"plaintext":      string(secret.Plaintext),
		"generation":     secret.Generation,
		"metageneration": secret.Metageneration,
	}); err != nil {
		return err
	}

	return nil
}

func resourceBerglasSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	bucket, object, generation, err := decodeId(d.Id())
	if err != nil {
		return err
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
			return err
		}

		id := encodeId(bucket, secret.Name, secret.Generation)
		d.SetId(id)

		if err := setMany(d, resourceFields{
			"generation":     secret.Generation,
			"metageneration": secret.Metageneration,
			"plaintext":      string(secret.Plaintext),
		}); err != nil {
			return err
		}

		if err := resourceBerglasSecretRead(d, meta); err != nil {
			return err
		}
	}

	return nil
}

func resourceBerglasSecretDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)
	client, ctx := config.Client(), config.Context()

	bucket, object, _, err := decodeId(d.Id())
	if err != nil {
		return err
	}

	if err := client.Delete(ctx, &berglas.DeleteRequest{
		Bucket: bucket,
		Object: object,
	}); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceBerglasSecretImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	bucket, object, generation, err := decodeId(d.Id())
	if err != nil {
		return nil, err
	}

	if err := setMany(d, resourceFields{
		"bucket":     bucket,
		"name":       object,
		"generation": generation,
	}); err != nil {
		return nil, err
	}

	if err := resourceBerglasSecretRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
