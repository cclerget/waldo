# Validated model ladder

This directory contains only reproducible model recipes that passed every
precommitted numerical and behavioral promotion gate. It is a catalog of
known-good reference models, not an experiment log.

## Promotion policy

New candidates are trained and evaluated from a compose on an experiment
branch or from a temporary path. A candidate is copied into this directory
only after all of its gates pass. A failed or sub-par candidate is removed
before the branch is merged; Git and the persisted model artifacts retain its
history.

Promotion requires all of the following:

1. finite training and held-out metrics;
2. exact agreement between the selected, serialized, and reloaded artifact;
3. the precommitted normalized-loss improvement;
4. every corpus-specific behavioral gate; and
5. a written result with model ID, run ID, token exposure, and evaluation
   evidence.

A lower aggregate loss alone is not a promotion.

## How loss is optimized and compared

Training directly minimizes next-token cross-entropy in nats per token. Raw
token loss is comparable only when the tokenizer and held-out partition are
identical. When tokenizers differ, this ladder compares bits per UTF-8 byte:

```text
BPB = held_out_loss / (held_out_bytes_per_token * ln(2))
```

Lower BPB means fewer bits are needed to predict the same held-out source
text. It prevents a tokenizer that emits longer or shorter token sequences
from appearing better merely because it changed the unit of measurement.

The controlled optimization levers are:

- **data:** quality filters, distribution, unique source exposure, and held-out
  partition;
- **tokenizer:** vocabulary size, fertility, effective byte context, and
  document-boundary density;
- **capacity:** core parameters and the fraction consumed by token embeddings;
- **optimization:** learning rate, warmup, decay schedule, batch geometry,
  weight decay, and checkpoint selection; and
- **budget:** source passes and tokens per core parameter.

Only one causal lever changes between a candidate and its reference whenever
possible. The current PressBooks capacity experiment, for example, holds the
tokenizer, context, data, source passes, optimizer, batch geometry, seed, and
evaluation set fixed. It changes core capacity and applies the predeclared
square-root learning-rate adjustment required by that capacity change. The
gate is BPB plus generation quality, not training loss alone.

## Rung 0001: Tiny Shakespeare reference

[`0001-tiny-shakespeare.yaml`](0001-tiny-shakespeare.yaml) is the small
end-to-end systems reference:

- model `tiny-shakespeare-control-01`, ID `b4f8477a55ad`;
- 10.7M parameters and 81.92M byte tokens;
- held-out loss improved from 5.3814 to 1.5117;
- step 1,250 was selected and reloaded after later overtraining;
- sampled output preserved play formatting and locally plausible language; and
- training completed in under 15 minutes on two H200 GPUs.

```console
go run ./cmd/waldo/ model forecast composes/0001-tiny-shakespeare.yaml
go run ./cmd/waldo/ model train tiny-shakespeare-control-01 \
  composes/0001-tiny-shakespeare.yaml
./composes/evaluate-tiny-shakespeare.sh \
  tiny-shakespeare-control-01 /tmp/tiny-shakespeare-control-01-eval.jsonl
```

## Rung 0002: TinyStories BPE reference

[`0002-tinystories-bpe.yaml`](0002-tinystories-bpe.yaml) is the first
TinyStories recipe that passed every promotion gate:

- model `tinystories-bpe-20tpp-01`, ID `5d18938fa102`, run
  `3ee94d2f3e466ce2`;
- 32.8M parameters, a 4K PressBooks-trained byte BPE, and 512-token context;
- 613,613,568 packed targets, or 20 tokens per core parameter;
- held-out loss 1.042866 and 0.498234 BPB;
- seven of eight sampled continuations retained prompt entities and causality;
- no sampled and one greedy continuation entered a clear repetition loop; and
- all eight long-horizon samples emitted EOS.

```console
go run ./cmd/waldo/ model forecast composes/0002-tinystories-bpe.yaml \
  --preflight
go run ./cmd/waldo/ model train tinystories-bpe-20tpp-01 \
  composes/0002-tinystories-bpe.yaml --hostfile ~/hostfile
./composes/evaluate-tinystories.sh \
  tinystories-bpe-20tpp-01 /tmp/tinystories-bpe-eval.jsonl
./composes/evaluate-tinystories.sh \
  tinystories-bpe-20tpp-01 /tmp/tinystories-bpe-greedy.jsonl 0 42
./composes/evaluate-tinystories.sh \
  tinystories-bpe-20tpp-01 /tmp/tinystories-bpe-eos.jsonl 0.8 42 1536
```

## Candidate workflow

For the next rung:

1. write the hypothesis and immutable gates first;
2. train using a compose outside the promoted directory;
3. retain the exact compose with the run artifacts;
4. evaluate normalized loss and fixed behavioral probes; and
5. add it here with the next contiguous number only if every gate passes.

The scripts in this directory are evaluation utilities and are not rungs.
