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
	"testing"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBerglasSecret_basic(t *testing.T) {
	t.Parallel()

	bucket := testAccBucket(t)
	name := "terraform-" + acctest.RandString(24)
	key := testAccKey(t)
	ctx := context.Background()

	// Create a secret for reading
	secret, err := berglas.Create(ctx, &berglas.CreateRequest{
		Bucket:    bucket,
		Object:    name,
		Plaintext: []byte("testing123"),
		Key:       key,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup the secret
	defer func() {
		if err := berglas.Delete(ctx, &berglas.DeleteRequest{
			Bucket: bucket,
			Object: name,
		}); err != nil {
			t.Error(err)
		}
	}()

	rn := "data.berglas_secret.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataBerglasSecret_basic(t, bucket, name, secret.Generation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "id",
						fmt.Sprintf("%s/%s#%d", bucket, name, secret.Generation)),
					resource.TestCheckResourceAttrSet(rn, "bucket"),
					resource.TestCheckResourceAttrSet(rn, "name"),
					resource.TestCheckResourceAttrSet(rn, "plaintext"),
				),
			},
		},
	})
}
func testDataBerglasSecret_basic(t testing.TB, bucket, name string, generation int64) string {
	return fmt.Sprintf(`
data "berglas_secret" "test" {
	bucket     = "%s"
	name       = "%s"
	generation = "%d"
}`, bucket, name, generation)
}
