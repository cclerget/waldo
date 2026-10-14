# Candidate rung 0003: genetics and cell biology

This candidate tests whether the passed 30.7M-core TinyStories architecture can
learn one bounded knowledge domain without losing its proven language and EOS
behavior. The target is genetics and cell biology—not general elementary
knowledge.

The previous broad Cosmopedia candidate is not a validated rung. It reached a
held-out loss of 1.6674 and retained stories and EOS, but scored only 1/10 on
the frozen factual probes. It learned synthetic textbook style rather than
dependable facts. Git and model `9733035d296b` retain that evidence.

This replacement remains **unpromoted** until every gate below passes.

## Frozen experiment

- Architecture and 4K PressBooks-trained byte-BPE tokenizer remain unchanged
  from rung 0002.
- Training exposure remains 80 tokens per core parameter: 2,454,448,128 packed
  targets after complete optimizer steps.
- The weighted token mixture is:
  - 70% peer-reviewed articles from the `PLOS Genetics` journal;
  - 10% Wikipedia article families and Wikibook chapters whose titles are
    explicitly about molecular genetics, cell biology, or biochemistry; and
  - 20% TinyStories reference for language, causality, and EOS retention.
- Wikimedia selection requires `wiki` to be `wikipedia.com` or
  `wikibooks.com` plus one of the explicit title patterns in `compose.yaml`.
  The Wikibooks hierarchy depths are enumerated because shell `*` does not
  cross `/`. Namespace-zero main-content selection was already applied during
  ingestion.
- No Cosmopedia, broad Wikimedia, other PLOS journals, arithmetic, geography,
  astronomy, history, literature, QA, or instruction data is included.

The PLOS and Wikimedia records are independently authored, grounded sources.
PLOS provides volume and factual consistency within one research field;
the title-selected reference material supplies definitions and explanatory
structure closer to the raw continuation probes.

## Promotion gates

All gates are frozen before training:

- finite metrics and exact selected-artifact reload;
- held-out loss improves without a material late regression;
- at least 8 of 10 deterministic factual continuations begin with the expected
  genetics or cell-biology fact;
- at least 4 of 5 explanation probes remain scientifically relevant through
  their first two sentences;
- at least 7 of 8 rung-0002 story probes retain named entities and causal
  direction;
- no more than 1 of 10 deterministic factual probes enters an obvious
  repetition loop; and
- at least 90% of the long-horizon factual suite terminates with EOS.

This remains a raw causal foundation model. Chat behavior, instruction
following, concise answers, and assistant formatting are out of scope.

## Preflight evidence

The frozen selection was measured on the indexed corpora before training:

- 269.5M unique tokenizer targets and 2.5B planned targets, for 9.11 aggregate
  effective passes;
- 5.6K `PLOS Genetics` records at 15.61 effective passes;
- 3.6K title-selected Wikimedia records at 35.20 effective passes;
- 508.9K TinyStories records at 3.22 effective passes; and
- 90.3% of selected records fit within one 512-token sequence.

The Wikimedia replay is intentionally concentrated, but it accounts for only
10% of planned tokens and supplies the bounded reference facts measured by the
frozen probes. This candidate must not be generalized beyond the declared
molecular-genetics and cell-biology capability.

## Run

From the repository root:

```console
go run ./cmd/waldo/ model forecast \
  composes/0003-cell-biology/compose.yaml --preflight
go run ./cmd/waldo/ model train cell-biology-30m-80tpp-01 \
  composes/0003-cell-biology/compose.yaml --hostfile ~/hostfile

./composes/0003-cell-biology/evaluate.sh \
  cell-biology-30m-80tpp-01 facts /tmp/cell-biology-facts.jsonl
./composes/0003-cell-biology/evaluate.sh \
  cell-biology-30m-80tpp-01 explanations /tmp/cell-biology-explanations.jsonl 0.7 42 160
./composes/0003-cell-biology/evaluate.sh \
  cell-biology-30m-80tpp-01 stories /tmp/cell-biology-stories.jsonl 0.8 42 256
./composes/0003-cell-biology/evaluate.sh \
  cell-biology-30m-80tpp-01 eos /tmp/cell-biology-eos.jsonl 0.8 42 1536
```

Before training, the preflight must show only `PLOS Genetics`, the declared
molecular-biology Wikimedia view, and TinyStories. Review both total and
per-corpus effective passes; do not infer suitability from the aggregate token
budget alone.

## Promotion decision

Record the model ID, run ID, actual weighted consumption, effective source
passes, held-out curve, selected checkpoint, BPB, and every behavioral score.
Only after every gate passes should this candidate be considered the validated
rung 0003.
