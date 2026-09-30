# Changelog

## v1.0.0 — Strict binary SML

Breaking transition from device-specific repair to specification-based parsing.

- Remove all Tibber Pulse TL/CRC normalization, including experimental changes.
- Validate schema, required fields, choices, sequence lengths and file envelopes.
- Require original-byte message and transport CRCs, end markers and valid fill.
- Decode escaped version-1 transport payloads through `TransportParse`.
- Correct integer sign extension, unsigned widths, TL overflow and status values.
- Preserve the three standard time choices and unsigned64 values in the API.
- Support SML_ListType time values and ProcParValue list entries.
- Add specification-derived tests, fuzz seeds and full external corpus checks.
- Retain malformed real-world fixtures as negative tests.

See [STANDARD.md](STANDARD.md) for normative references, implemented scope,
unsupported layers and migration. No claim of complete XML/COSEM/transport-v2
support or formal certification is made.
