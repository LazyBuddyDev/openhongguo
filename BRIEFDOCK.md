# Vendored App protocol SDK

Source: https://github.com/beilunyang/openhongguo

Commit: `015a388112ee148a7d2dfb7ecfafe171bbab957f` (2026-11-12).

Fork review: upstream changes through `cac5dbc0e1150c0f93e309d8de68bb483f50f499` were reviewed on 2026-11-24. Selected fixes and exclusions are recorded in [docs/upstream-review.md](docs/upstream-review.md). The commit above remains the source baseline; this is not a complete upstream version upgrade. Updating this fork does not update hongguo-service or an installed module.

License: Apache-2.0, retained in LICENSE. Only required SDK sources, the memory profile store and focused BriefDock regression tests are included. Local changes are listed below.

BriefDock supplies its own narrowly scoped transport. Account cookies pass through private pipes only to explicitly allowed video, shelf and watch-progress endpoints. Account-bound device profiles use application memory and separate Windows Credential Manager storage; QA only reads that storage. Retries and automatic device resets are disabled. The upstream SDK is an unofficial client, not an official authorization API. A media model is not proof of successful playback or permission for every episode.

BriefDock adds sdk/briefdock_account.go: one explicit authenticated GET video-shelf endpoint, using the existing fixed SDK device and signing code. The custom transport permits this read and the fixed POST video-shelf update.

BriefDock also independently adds fixed GET read_progress/list/v1/ and POST
read_progress/upload/v1/ account checkpoint operations. The transport permits only
these exact method/path pairs, no history/favorite deletion or arbitrary endpoint.
Protocol field reference (not copied implementation):
https://github.com/tsukiyo-zero/hongguoTV/blob/77a99f509b8a5c586bf6ef453147cf1003203970/kotlin-tv/core/src/main/kotlin/com/hongguotv/core/CloudAccount.kt
Playback records are timestamp-merged and uploads are read back before confirmation;
background test instances require a separate explicit progress-write flag.

The fixed POST /reading/bookapi/bookshelf/video/update/v1/ independently implements
single-series add/cancel wire fields from the same pinned protocol reference above.
No collection replacement, bulk cancellation, arbitrary URL or automatic write retry
is accepted. Favorite state is read before the operation and read back afterward;
isolated writes additionally require BRIEFDOCK_FAVORITES_REAL=1.

Playback adds `sdk/briefdock_playback.go` and `internal/protocol/sign/briefdock_playback.go`:
one fixed GET player-model endpoint, full Argus/Ladon using existing Apache SDK
primitives, no account-style legacy signature headers, no endpoint/device rotation.
Wire-field reference only (no copied implementation):
https://github.com/tsukiyo-zero/hongguoTV/blob/77a99f509b8a5c586bf6ef453147cf1003203970/kotlin-tv/core/src/main/kotlin/com/hongguotv/core/ContentRepository.kt
and `Signer.kt` / `Argus.kt` at the same commit.

`internal/protocol/parse.go` and `sdk/model/media.go` retain the upstream generic
parser and add a quality-aware path: select supported AVC/HEVC renditions, prefer
the closest resolution at or above the request, keep that rendition's own wrapped
material, and report actual/available resolutions. ByteVC2 is rejected instead of
being treated as HEVC. Extra compatibility-model requests were tested without
benefit and removed; signed fallback addresses are not used. Focused tests cover
wrong-episode rejection, mixed codecs/materials, higher-resolution selection and
no silent downgrade.
