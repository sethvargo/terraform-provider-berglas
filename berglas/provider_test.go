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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatal(err)
	}
}

func testAccBucket(tb testing.TB) string {
	v := os.Getenv("TEST_ACC_BERGLAS_BUCKET")
	if v == "" {
		tb.Fatal("missing TEST_ACC_BERGLAS_BUCKET")
	}
	return v
}

func testAccKey(tb testing.TB) string {
	v := os.Getenv("TEST_ACC_BERGLAS_KEY")
	if v == "" {
		tb.Fatal("missing TEST_ACC_BERGLAS_KEY")
	}
	return v
}

func testAccPreCheck(tb testing.TB) {}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"berglas": testAccProvider,
	}
}
