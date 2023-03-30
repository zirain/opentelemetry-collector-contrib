// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istiometadataexchangeprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/istiometadataexchangeprocessor"

// Config defines configuration for Resource processor.
type Config struct {
	WASMPeerSetteings `mapstructure:",squash"`
}

type WASMPeerSetteings struct {
	// The name of the upstream peer metadata exchange wasm filter.
	// If not specified, the default value is "wasm.upstream_peer".
	UpstreamPeerKey string `mapstructure:"upstream_peer_key"`
	// The name of the downstrem peer metadata exchange wasm filter.
	// If not specified, the default value is "wasm.downstream_peer".
	DownstreamPeerKey string `mapstructure:"downstream_peer_key"`
	// Whether to remove the peer metadata exchange wasm filter attributes.
	RemovePeerKey bool `mapstructure:"remove_peer_key"`
}

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
