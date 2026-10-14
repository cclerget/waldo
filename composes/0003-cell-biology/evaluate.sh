#!/usr/bin/env bash
# Evaluate candidate rung 0003 with genetics and cell-biology probes.
set -euo pipefail

if (($# < 3 || $# > 6)); then
  echo "usage: $0 MODEL {facts|explanations|stories|eos} OUTPUT.jsonl [TEMPERATURE] [SEED] [MAX_TOKENS]" >&2
  exit 2
fi

model="$1"
suite="$2"
output="$3"
temperature="${4:-0}"
seed="${5:-42}"
max_tokens="${6:-80}"
run_id="$(go run ./cmd/waldo/ --json model summary "$model" | jq -r '.bom.current_run_id')"

facts=(
  'DNA is the molecule that stores'
  'A gene is a segment of'
  'The four bases in DNA are'
  'In DNA, adenine pairs with'
  'RNA contains uracil instead of'
  'Ribosomes build'
  'The cell membrane controls'
  'Mitosis produces'
  'Meiosis produces'
  'Proteins are made from'
)

explanations=(
  'DNA replication begins when'
  'During transcription, a cell'
  'During translation, a ribosome'
  'Mitosis and meiosis differ because'
  'A genetic mutation can affect a protein when'
)

stories=(
  $'Alice was bored and wanted to find some adventures. She walked up to her friend Ben, who looked very busy playing with his toys. Alice said, "Why don\x27t we'
  $'Alice walked up to her friend Ben\x27s house. She was planning to ask him to go to the park with her. When Ben opened the door, she asked him if he had any plans. He said, "I\x27m sorry, Alice, but'
  $'The day before Ben\x27s birthday, Alice walked into the kitchen and saw Ben sitting there, looking gloomy. She said, "Ben, why are you'
  $'Once upon a time there was a curious boy who lived in a house with a big garden. Every day he explored the garden and found new surprises. But one day, it was raining so hard that his mother told him'
  $'Once upon a time, there was a tiger who liked to play the guitar. One day, a bunny heard the guitar from a distance and'
  $'Alice wanted to play with her doll, but she could not remember where she had put it. She looked all around the house but could not find it, so she decided'
  $'"Ben, what do you have in your pocket?" Alice asked. "Oh, nothing," Ben replied. But Alice saw that there was definitely something in Ben\x27s pocket, and she was very curious what it was, so she'
  $'One day, a bird was flying high over the sea. The bird noticed a small boat with a boy sitting inside. The boy looked lost, so'
)

case "$suite" in
  facts) prompts=("${facts[@]}") ;;
  explanations) prompts=("${explanations[@]}") ;;
  stories) prompts=("${stories[@]}") ;;
  eos) prompts=("${facts[@]}") ;;
  *) echo "unknown suite: $suite" >&2; exit 2 ;;
esac

: >"$output"
for prompt in "${prompts[@]}"; do
  go run ./cmd/waldo/ --json model chat "$model" \
    --run-id "$run_id" --raw --temperature "$temperature" --seed "$seed" \
    --max-tokens "$max_tokens" "$prompt" >>"$output"
done

jq -s --arg suite "$suite" --argjson max_tokens_requested "$max_tokens" '{
  suite: $suite,
  responses: length,
  max_tokens_requested: $max_tokens_requested,
  eos: map(select(.result.finish_reason == "eos")) | length,
  max_tokens: map(select(.result.finish_reason == "max_tokens")) | length,
  results: map({prompt, tokens: .result.tokens, finish_reason: .result.finish_reason, text: .result.text})
}' "$output"
