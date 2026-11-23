# LazyBuddy fork

`main` preserves the upstream branch at the time of the fork. `lazybuddy` starts
from the SDK revision used by hongguo-service, `015a388112ee148a7d2dfb7ecfafe171bbab957f`.
It contains the playback and account adaptations described in [BRIEFDOCK.md](BRIEFDOCK.md).

The [upstream workflow](.github/workflows/upstream.yml) checks upstream `main`
daily and opens one review issue per new tip. Review the comparison, select
relevant fixes and run regression tests before updating the reviewed baseline
in the workflow. The workflow does not merge upstream commits.
GitHub can disable the schedule after 60 days of repository inactivity. Re-enable
it from the Actions tab when needed.

hongguo-service builds from a tracked SDK snapshot in `vendor/openhongguo`.
When adopting a change from this fork, update that snapshot in the service
repository and record the exact fork commit alongside the update. Review
playback on an affected episode and a known-good episode before shipping changes
to playback resolution.
