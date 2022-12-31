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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBerglasSecret_basic(t *testing.T) {
	t.Parallel()

	bucket := testAccBucket(t)
	name := "terraform-" + acctest.RandString(24)
	key := testAccKey(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testProviderFactories,
		CheckDestroy:      testAccBerglasSecretDestroy(t, bucket, name),
		Steps: []resource.TestStep{
			{
				Config: testBerglasSecret_basic(t, bucket, name, key),
				Check: resource.ComposeTestCheckFunc(
					testAccBerglasSecret(t, bucket, name),
					resource.TestCheckResourceAttr("berglas_secret.test", "bucket", bucket),
					resource.TestCheckResourceAttr("berglas_secret.test", "name", name),
					resource.TestCheckResourceAttr("berglas_secret.test", "plaintext", "super-secret"),
				),
			},
			{
				ResourceName:      "berglas_secret.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBerglasSecret(t testing.TB, bucket, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := New("test")().Meta().(*config)
		client := config.Client()

		ctx := context.Background()
		if _, err := client.Read(ctx, &berglas.ReadRequest{
			Bucket: bucket,
			Object: name,
		}); err != nil {
			return fmt.Errorf("failed to get secret: %w", err)
		}

		return nil
	}
}

func testAccBerglasSecretDestroy(t testing.TB, bucket, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := New("test")().Meta().(*config)
		client := config.Client()

		ctx := context.Background()
		if _, err := client.Read(ctx, &berglas.ReadRequest{
			Bucket: bucket,
			Object: name,
		}); err == nil {
			return fmt.Errorf("expected resource to be deleted")
		}

		return nil
	}
}

func testBerglasSecret_basic(t testing.TB, bucket, name, key string) string {
	return fmt.Sprintf(`
resource "berglas_secret" "test" {
	bucket    = "%s"
	name      = "%s"
	key       = "%s"
	plaintext = "super-secret"
}`, bucket, name, key)
}
