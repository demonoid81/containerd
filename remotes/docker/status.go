/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package docker

import (
	"sync"

	"github.com/demonoid81/containerd/content"
	"github.com/demonoid81/containerd/errdefs"
	"github.com/pkg/errors"
)

// Status of a content operation
type Status struct {
	content.Status

	// UploadUUID is used by the Docker registry to reference blob uploads
	UploadUUID string
}

// StatusTracker to track status of operations
type StatusTracker interface {
	GetStatus(string) (Status, error)
	SetStatus(string, Status)
}

type memoryStatusTracker struct {
	statuses map[string]Status
	m        sync.Mutex
}

// NewInMemoryTracker returns a StatusTracker that tracks content status in-memory
func NewInMemoryTracker() StatusTracker {
	return &memoryStatusTracker{
		statuses: map[string]Status{},
	}
}

func (t *memoryStatusTracker) GetStatus(ref string) (Status, error) {
	t.m.Lock()
	defer t.m.Unlock()
	status, ok := t.statuses[ref]
	if !ok {
		return Status{}, errors.Wrapf(errdefs.ErrNotFound, "status for ref %v", ref)
	}
	return status, nil
}

func (t *memoryStatusTracker) SetStatus(ref string, status Status) {
	t.m.Lock()
	t.statuses[ref] = status
	t.m.Unlock()
}
