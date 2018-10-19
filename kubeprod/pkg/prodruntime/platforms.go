/*
 * Bitnami Kubernetes Production Runtime - A collection of services that makes it
 * easy to run production workloads in Kubernetes.
 *
 * Copyright 2018 Bitnami
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prodruntime

import (
	"fmt"
	"net/url"
)

type Platform struct {
	Name        string
	Description string
}

// Platforms should append themselves to this in init()
var Platforms = []Platform{}

func FindPlatform(name string) *Platform {
	for i := range Platforms {
		p := &Platforms[i]
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (p *Platform) ManifestURL(base *url.URL) (*url.URL, error) {
	return base.Parse(fmt.Sprintf("platforms/%s.jsonnet", p.Name))
}
