#!/usr/bin/env bash
# Copyright (c) Microsoft. All rights reserved.
# Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

set -euo pipefail
set +x  # Never trace secrets

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SECRETS_DIR="${BUILD_SOURCESDIRECTORY:-$REPO_ROOT}/secrets"
SUBSCRIPTIONS_FILE="$SECRETS_DIR/test.subscriptions.regions.json"

if [[ ! -f "$SUBSCRIPTIONS_FILE" ]]; then
  echo "ERROR: Subscriptions JSON not found at $SUBSCRIPTIONS_FILE"
  exit 1
fi

if ! command -v jq &>/dev/null; then
  echo "ERROR: jq is required but not installed."
  exit 1
fi

export SPEECH_SUBSCRIPTION_KEY=$(jq -jr '.UnifiedSpeechSubscription.Key' "$SUBSCRIPTIONS_FILE")
export SPEECH_SUBSCRIPTION_REGION=$(jq -jr '.UnifiedSpeechSubscription.Region' "$SUBSCRIPTIONS_FILE")

if [[ -z "$SPEECH_SUBSCRIPTION_KEY" || "$SPEECH_SUBSCRIPTION_KEY" == "null" ]]; then
  echo "ERROR: Failed to extract UnifiedSpeechSubscription.Key from $SUBSCRIPTIONS_FILE"
  exit 1
fi

if [[ -z "$SPEECH_SUBSCRIPTION_REGION" || "$SPEECH_SUBSCRIPTION_REGION" == "null" ]]; then
  echo "ERROR: Failed to extract UnifiedSpeechSubscription.Region from $SUBSCRIPTIONS_FILE"
  exit 1
fi

echo "Loaded speech subscription for region: $SPEECH_SUBSCRIPTION_REGION"

GLOBAL_STRINGS_TO_REDACT=(
  "$SPEECH_SUBSCRIPTION_KEY"
)

URL_ENCODED_SPEECH_SUBSCRIPTION_KEY=$(jq -nr --arg v "$SPEECH_SUBSCRIPTION_KEY" '$v | @uri')
if [[ -n "$URL_ENCODED_SPEECH_SUBSCRIPTION_KEY" && "$URL_ENCODED_SPEECH_SUBSCRIPTION_KEY" != "$SPEECH_SUBSCRIPTION_KEY" ]]; then
  GLOBAL_STRINGS_TO_REDACT+=("$URL_ENCODED_SPEECH_SUBSCRIPTION_KEY")
fi
unset URL_ENCODED_SPEECH_SUBSCRIPTION_KEY

redact_input_with() {
  perl -MIO::Handle -lpe \
    'BEGIN {
       STDOUT->autoflush(1);
       STDERR->autoflush(1);
       if (@ARGV) {
         $re = sprintf "(?:%s)", (join "|", map { quotemeta $_ } splice @ARGV);
         $re = qr/$re/
       }
     }
     $re and s/$re/***/gi' "$@"
}

global_redact() {
  redact_input_with "${GLOBAL_STRINGS_TO_REDACT[@]}"
}

export -f redact_input_with
export -f global_redact
