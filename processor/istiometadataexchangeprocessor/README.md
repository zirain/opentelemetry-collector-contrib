# Istio Metadata Exchange Processor

| Status                   |               |
|--------------------------|---------------|
| Stability                | [development] |
| Supported pipeline types | logs          |
| Distributions            | none          |

NOTE - This processor is experimental, with the intention that its functionality will be reimplemented in the [istio metadata exchange processor](../istiometadataexchangeprocessor/README.md) in the future.

The logs transform processor can be used to apply [log operators](../../pkg/stanza/docs/operators) to logs coming from any receiver.
Please refer to [config.go](./config.go) for the config spec.

Examples:

```yaml
processors:
  istiometadataexchange:
```

Refer to [config.yaml](./testdata/config.yaml) for detailed
examples on using the processor.

[development]: https://github.com/open-telemetry/opentelemetry-collector#development
