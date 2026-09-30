# Strict binary SML

## Normative reference

This implementation follows the published **BSI TR-03109-1, Anlage IVb,
Version 1.0, 18 March 2013**, containing SML **1.04**, draft 13 July 2010.

[Official BSI publication](https://www.bsi.bund.de/SharedDocs/Downloads/DE/BSI/Publikationen/TechnischeRichtlinien/TR03109/TR-03109-1_Anlage_Feinspezifikation_Drahtgebundene_LMN-Schnittstelle_Teilb.html)

The library version is independent of the SML protocol version. This is a
parser library, not a smart-meter server or a certification claim.

## Implemented scope and traceability

| Specification | Implementation |
| --- | --- |
| Chapter 6.1, tables 4–5 | TL types, continuation bits, multi-byte lengths, bounds and overflow checks (`shared.go`) |
| Chapter 6.2.1–6.2.4 | Octet strings, signed/unsigned integers, Boolean `42 XX`; no type aliases (`number.go`, `value.go`, `boolean.go`) |
| Chapter 6.2.5, 6.3.2–6.3.4 | Sequence arity, CHOICE tags, required/optional fields (`schema.go`) |
| Definitions F–H | All three time choices, including local/seasonal offsets in minutes (`time.go`) |
| Definitions V, X | SML_Value primitives and SML_ListType/smlTime (`value.go`) |
| Definitions FF–II | Parameter trees, all five ProcParValue choices (`tree.go`) |
| Chapters 5.1.1–5.1.15 | Open/Close, GetList, profile and parameter messages, Attention; schema validation before decoding |
| Paragraph 29 | Mandatory message CRC over **original bytes**, excluding crc16 TL/value and EndOfSmlMsg (`message.go`) |
| Chapter 8.1, paragraphs 153–158 | Transport version 1, escaping, alignment, zero fill, mandatory CRC over original wire bytes (`transport.go`) |
| Chapter 4.1 and paragraph 27 | File Open/Close envelope and contiguous message groups (`file.go`) |

Unsupported: XML encoding, COSEM services and transport-version-2 block
negotiation/state machines. Unsupported message identifiers and detected
version-2 transport are errors; the library does not guess their meaning.
The referenced draft itself contains unfinished COSEM definitions. A complete
implementation of every layer must not be inferred from this release.

Application responsibilities include authentication, request/response matching,
OBIS semantics, signature verification and execution of abortOnError policies.
These are distinct from structural decoding and CRC integrity checks.

## Integrity policy

No Pulse-specific repair or CRC variant is implemented. A Boolean is not an
unsigned integer; `40`, `43` or `45` cannot substitute for an unsigned TL field.
CRC validation never substitutes bytes or searches for possible corrections.
Failure returns an error and file/transport parsing returns no partial readings.

`TransportRead` returns validated **wire bytes**. To decode a complete frame,
use `TransportParse`, or `TransportPayload` followed by `FileParse`. Do not
strip eight bytes from each end: that misses escaping and declared fill.
`FileParse` expects an unframed, unescaped SML file, without transport padding.
For explicitly application-defined standalone messages, use `MessageParse`.
Primitive helpers allow the `01` optional marker; enclosing schema validators
decide where it is permitted. Empty octet strings have the same encoding.

## Breaking migration from v0.1.x

- `Time` is a struct with `Present`, `Tag`, `Timestamp`, `LocalOffset` and
  `SeasonTimeOffset`, rather than a uint32 that loses the time choice.
- Unsigned `Value` numbers are in `DataUnsigned` (uint64); signed numbers stay
  in `DataInt`. `DataTime` represents the SML_ListType time alternative.
- `ListEntry.Status` and `StatusParse` are uint64 and preserve the actual bits.
- `PulseU8Parse` and all implicit Pulse compatibility rules are removed.
- `MessageParse` always validates CRC, including when its obsolete variadic
  boolean argument is false. That argument is retained only for source migration.
- Reserved types, malformed TL fields, missing required fields/terminators,
  invalid file envelopes, bad transport escapes/fill and CRC failures are errors.

## Tests

Portable tests: `go test ./...`, `go vet ./...`.

External corpus (unchanged upstream revision
`a3c7869fca8ad73c34248f6b4eda115acf7c79dc`): set `LIBSML_TESTING_DIR` to a clone
of [devZer0/libsml-testing](https://github.com/devZer0/libsml-testing), then run
`go test -run TestExternalLibSMLFixtures -v`. The test checks every bin/hex pair
and audits every reachable frame against the local fork, asserting the recorded
strict results. The corpus is not a promise that all captured bytes are valid.
Results: 37 pairs; 220 accepted frames; 36 rejected frames or incomplete tails.
Rejections include bad CRCs, damaged dumps and Holley ZDBA's missing time CHOICE.
The upstream fixtures and expectations are not modified.

The old damaged Pulse fixtures remain byte-for-byte regression inputs, now
expected to be rejected. The valid live capture still yields three messages
and the expected consumption value. To audit the 60-capture stream, set
`PULSE_CAPTURE_DIR` and run `go test -run TestStrictPulseLongStream -v`.
An independent bitwise CRC reference confirms 27 valid frames and 33 frames
with invalid transport CRCs. Input bytes are never modified.

Fuzzing: `go test -run '^$' -fuzz FuzzStrictMessageParser -fuzztime 10s`.
