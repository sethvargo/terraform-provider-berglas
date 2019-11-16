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
	"sync"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
)

type config struct {
	lock sync.RWMutex

	client *berglas.Client
	ctx    context.Context
}

// Client returns the configured berglas client.
func (c *config) Client() *berglas.Client {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.client
}

// Context returns the context on the config.
func (c *config) Context() context.Context {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.ctx
}
