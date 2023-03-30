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

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/istiometadataexchangeprocessor/exchange"
)

type istioMetadataExchangeProcessor struct {
	logger *zap.Logger
	config *Config
}

func (proc *istioMetadataExchangeProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				attributes := lr.Attributes()
				upstreamKey := proc.upstreamPeerKey()
				downstreamKey := proc.downstreamPeerKey()

				var peerInfo map[string]string
				var err error
				upstreamPeer, ok := attributes.Get(upstreamKey)
				if ok && !emptyString(upstreamPeer.AsString()) {
					peerInfo, err = istioPeerData(upstreamPeer.AsString(), true)
					if err != nil {
						proc.logger.Error("failed to parse upstream peer", zap.Error(err))
					}
				}

				downstremPeer, ok := attributes.Get(downstreamKey)
				if ok && !emptyString(downstremPeer.AsString()) {
					peerInfo, err = istioPeerData(downstremPeer.AsString(), true)
					if err != nil {
						proc.logger.Error("failed to parse downstream peer", zap.Error(err))
					}
				}

				proc.logger.Debug("istio peer info", zap.Any("peerInfo", peerInfo))
				for k, v := range peerInfo {
					attributes.PutStr(k, v)
				}

				if proc.config.RemovePeerKey {
					proc.logger.Debug("removing peer key")
					attributes.Remove(upstreamKey)
					attributes.Remove(downstreamKey)
				}
			}
		}
	}
	return ld, nil
}

func istioPeerData(peerStr string, upstream bool) (map[string]string, error) {
	nodeInfo, err := getNodeInfo(peerStr)
	if err != nil {
		return nil, err
	}

	nodesLabels := map[string]string{}
	for i := 0; i < nodeInfo.LabelsLength(); i++ {
		var kv exchange.KeyVal
		if nodeInfo.Labels(&kv, i) {
			if kv.Key() == nil || kv.Value() == nil {
				continue
			}

			nodesLabels[string(kv.Key())] = string(kv.Value())
		}
	}

	if upstream {
		return map[string]string{
			"istio.destination.name":      string(nodeInfo.Name()),
			"istio.destination.namespace": string(nodeInfo.Namespace()),
			"istio.destination.app":       nodesLabels["app"],
			"istio.destination.version":   nodesLabels["service.istio.io/canonical-revision"],
			"istio.destination.service":   nodesLabels["service.istio.io/canonical-name"],
		}, nil
	}

	return map[string]string{
		"istio.sourece.name":      string(nodeInfo.Name()),
		"istio.sourece.namespace": string(nodeInfo.Namespace()),
		"istio.source.app":        nodesLabels["app"],
		"istio.source.version":    nodesLabels["service.istio.io/canonical-revision"],
		"istio.source.service":    nodesLabels["service.istio.io/canonical-name"],
	}, nil
}

func (proc *istioMetadataExchangeProcessor) downstreamPeerKey() string {
	return "wasm.downstream_peer"
}

func (proc *istioMetadataExchangeProcessor) upstreamPeerKey() string {
	return "wasm.upstream_peer"
}

func emptyString(s string) bool {
	return s == "" || s == "-"
}

func getNodeInfo(s string) (*exchange.FlatNode, error) {
	pb := &anypb.Any{}
	if err := protojson.Unmarshal([]byte(s), pb); err != nil {
		return nil, err
	}
	return exchange.GetRootAsFlatNode(pb.GetValue(), 3), nil
}
