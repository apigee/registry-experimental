// Copyright 2021 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorm

// Key represents a key in a storage provider
type Key struct {
	Kind string
	Name string
}

// NewKey creates a new storage key.
func (c *Client) NewKey(kind, name string) *Key {
	return &Key{Kind: kind, Name: name}
}

func (k *Key) String() string {
	return k.Kind + ":" + k.Name
}
