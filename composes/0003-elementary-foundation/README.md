# Candidate rung 0003: elementary foundation

This candidate asks one controlled question: can the passed 30.7M-core
TinyStories architecture acquire elementary explanatory language and basic
facts without losing its proven story behavior?

It is **not promoted**. This directory stays on its experiment branch until a
trained model passes every gate below. A failed candidate is removed rather
than merged into the validated ladder.

## Frozen experiment

- Architecture and 4K PressBooks-trained byte-BPE tokenizer are unchanged from
  rung 0002.
- Training exposure is 80 tokens per core parameter: 2,454,448,128 packed
  targets after complete optimizer steps.
- The weighted token mixture is exactly 70% filtered Cosmopedia and 30%
  TinyStories reference.
- Cosmopedia is restricted to records whose canonical metadata contains both
  `audience=middle_school_students` and `format=textbook`.
- The current TinyStories reference contains the pinned first of four source
  shards. That limitation is explicit; this candidate does not call it the
  full corpus.
- Checkpoint and evaluation observations occur every 9,987 steps, placing
  observations within one interval of 20, 40, 60, and 80 core TPP.

This deliberately uses only two source distributions. It does not add broad
Wikimedia, Stack Exchange, scientific papers, QA, or instruction data.

## Promotion gates

All gates are frozen before training:

- finite metrics and exact selected-artifact reload;
- held-out loss must improve through the run without a material late
  regression;
- at least 7 of 8 rung-0002 story probes retain their named entities and
  causal direction;
- at least 8 of 10 factual continuations are relevant and grammatical;
- at least 6 of 10 factual continuations contain the expected elementary fact
  or arithmetic result;
- at least 4 of 5 explanation probes begin a relevant explanation;
- no more than 1 of 10 deterministic probes enters an obvious repetition
  loop; and
- at least 90% of the long-horizon suite terminates with EOS.

This is a foundation model, not an assistant. The probes are raw causal
continuations and do not require chat behavior, instruction following, or
answer formatting.

## Run

From the repository root:

```console
go run ./cmd/waldo/ model forecast \
  composes/0003-elementary-foundation/compose.yaml --preflight
go run ./cmd/waldo/ model train elementary-foundation-30m-80tpp-01 \
  composes/0003-elementary-foundation/compose.yaml --hostfile ~/hostfile

./composes/0003-elementary-foundation/evaluate.sh \
  elementary-foundation-30m-80tpp-01 facts /tmp/elementary-facts.jsonl
./composes/0003-elementary-foundation/evaluate.sh \
  elementary-foundation-30m-80tpp-01 explanations /tmp/elementary-explanations.jsonl 0.7 42 160
./composes/0003-elementary-foundation/evaluate.sh \
  elementary-foundation-30m-80tpp-01 stories /tmp/elementary-stories.jsonl 0.8 42 256
./composes/0003-elementary-foundation/evaluate.sh \
  elementary-foundation-30m-80tpp-01 eos /tmp/elementary-eos.jsonl 0.8 42 1536
```

## Promotion decision

Record the model ID, run ID, actual corpus consumption, effective source
passes, held-out curve, selected checkpoint, BPB, and every behavioral score.
Only after every gate passes should this candidate be linked from the top-level
ladder README and merged as rung 0003.
