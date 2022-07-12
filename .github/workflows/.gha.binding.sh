#!/bin/bash

# note: yarn abigen already executed in go.yml
# fixme: we only check the contract api changes here
# shellcheck disable=SC2010
diff <(ls ./binding/*.go | grep --invert-match "artifacts" | xargs git diff) <(echo -n)