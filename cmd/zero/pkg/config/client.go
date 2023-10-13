// Copyright 2023 Google LLC. All Rights Reserved.
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

package config

import (
	"context"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/servicecontrol/v2"
)

const serviceAccountFile = ".config/zero/control.json"
const scope = servicecontrol.CloudPlatformScope

func GetClient(ctx context.Context) (*http.Client, error) {
	var client *http.Client
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path.Join(dirname, serviceAccountFile))
	if err != nil {
		return nil, err
	}
	creds, err := google.CredentialsFromJSON(ctx, data, scope)
	if err != nil {
		return nil, err
	}
	client = oauth2.NewClient(ctx, creds.TokenSource)
	return client, nil
}
