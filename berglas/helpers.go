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
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// encodeId encodes the ID from the given parts.
func encodeId(bucket, object string, generation int64) string {
	bucket, object = sanitizeBucket(bucket), sanitizeObject(object)

	id := bucket + "/" + object
	if generation > 0 {
		id = id + "#" + strconv.FormatInt(generation, 10)
	}

	return id
}

// decodeId explodes the ID into the given parts.
func decodeId(id string) (string, string, int64, error) {

	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", "", 0, errors.New("id must be {bucket}/{object}#{version}")
	}

	bucket, remainder := sanitizeBucket(parts[0]), parts[1]

	parts = strings.SplitN(remainder, "#", 2)
	object := sanitizeObject(parts[0])

	var generation int64
	if len(parts) > 1 {
		i, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return "", "", 0, errors.Wrap(err, "failed to parse generation")
		}
		generation = i
	}

	return bucket, object, generation, nil
}

// resourceFields are a map of kv pairs on a resource.
type resourceFields map[string]interface{}

// setMany wraps Set and handles any errors returned
func setMany(d *schema.ResourceData, m resourceFields) error {
	for k, v := range m {
		if err := d.Set(k, v); err != nil {
			return errors.Wrapf(err, "failed to set %q", k)
		}
	}
	return nil
}

// sanitizeBucket removes any gs:// or trailing / from the bucket name.
func sanitizeBucket(s string) string {
	return sanitizeObject(strings.TrimPrefix(s, "gs://"))
}

// sanitizeObject removes any leading or trailing spaces from the object.
func sanitizeObject(s string) string {
	return strings.Trim(s, "/")
}
