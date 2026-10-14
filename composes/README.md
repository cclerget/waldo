# Validated capability ladder

This directory contains only reproducible model recipes that passed every
precommitted numerical and behavioral promotion gate. It is a catalog of
known-good reference models, not an experiment log.

The ladder grows **capability, data, and model size together**. Each rung has a
realistic job. A model is promoted only when it performs that job; completing a
token budget or reaching a low aggregate loss is not enough.

## Directory layout

The top level contains this ladder plan and one directory for each promoted
rung. Every rung directory contains:

```text
README.md       purpose, evidence, gates, and commands
compose.yaml    immutable promoted training recipe
compose_test.go guards for the rung's model and training contract
evaluate.sh     fixed behavioral evaluation probes
```

Planned, failed, and sub-par candidates do not get rung directories. They stay
on experiment branches or in temporary workspaces until they pass promotion.

## What tokens per parameter means here

Twenty training tokens per core parameter is a useful compute-allocation
reference, not a universal definition of a trained model. It is enough for the
narrow TinyStories rung because that distribution is simple and repetitive. It
has not been enough for our broad factual models.

Very small deployment models are commonly trained far beyond the classical
compute-optimal point. For example, Hugging Face reports that SmolLM-135M was
trained on 600B tokens: about 4,444 tokens per advertised parameter. The
authors also report continued benchmark improvement beyond the Chinchilla
point, with gains beginning to slow after roughly 400B tokens. That model is
"overtrained" only in the narrow sense that a larger model could use the same
training compute more efficiently. It is not necessarily overtrained for the
goal of maximizing the capability of a fixed, inexpensive inference model.

Reference: [Hugging Face, *SmolLM - blazingly fast and remarkably
powerful*](https://huggingface.co/blog/smollm).

WALDO therefore reports and reasons about all of these separately:

- core parameters, excluding token input/output weights;
- packed training-token exposure;
- unique source tokens and effective corpus passes;
- held-out loss and bits per UTF-8 byte (BPB); and
- fixed behavioral capability gates.

For a new general-language rung, the budget checkpoints are 20, 40, 80, and
160 tokens per core parameter. The architecture, tokenizer, data, held-out
partition, and probes remain fixed across those checkpoints. The first
checkpoint that passes every gate may be promoted. Reaching 160 without
passing triggers a data/capacity review; it does not justify blindly training
to 4,444 tokens per parameter. SmolLM is evidence that tiny general models can
need much more data, not a budget prescription for this ladder.

## Promotion policy

New candidates are trained and evaluated from a compose on an experiment
branch or from a temporary path. A candidate is copied into this directory
only after all of its gates pass. A failed or sub-par candidate is removed
before the branch is merged; Git and persisted model artifacts retain its
history.

Promotion requires all of the following:

1. finite training and held-out metrics;
2. exact agreement between the selected, serialized, and reloaded artifact;
3. the precommitted normalized-loss improvement;
4. every rung-specific behavioral gate; and
5. a written result with model ID, run ID, token exposure, effective corpus
   passes, and evaluation evidence.

The behavioral probes and their scoring rules are frozen before training.
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

- **data:** quality, distribution, unique exposure, document boundaries, and
  the held-out partition;
- **tokenizer:** vocabulary size, fertility, effective byte context, and the
  fraction of parameters consumed by embeddings;
- **capacity:** core parameters, depth, width, attention, and MLP allocation;
- **optimization:** learning rate, warmup, decay, batch geometry, weight decay,
  regularization, and checkpoint selection; and
- **budget:** source passes and tokens per core parameter.

Within an experiment series, change one causal lever at a time. Compare exact
held-out BPB and behavioral scores at matched token checkpoints. Generation
quality, not a favorable training loss, decides promotion.

## Ladder definition

The first two rungs are passed references. Each promoted rung is a
self-contained directory with its compose, tests, evaluator, and evidence.
Later rungs are specifications, not promoted composes. Their directory enters
this tree only after the model passes.

| Rung | Target | Approximate core capacity | Required data | Budget policy | Minimum promotion capability | Status |
| --- | --- | ---: | --- | --- | --- | --- |
| [0001](0001-tiny-shakespeare/) | End-to-end systems control | 10M | Tiny Shakespeare | Fixed reference run | Locally plausible play continuation and correct checkpoint selection | Passed |
| [0002](0002-tinystories-bpe/) | Narrow story language model | 30M | TinyStories reference | 20 core TPP | Retain entities and causality, avoid loops, terminate with EOS | Passed |
| 0003 | Genetics and cell-biology foundation | 30-50M | PLOS Genetics, title-selected molecular-biology Wikimedia, and TinyStories retention data | 80 core TPP, then evidence review | Correct basic genetics and cell-biology continuations; retained stories; low repetition; reliable EOS | Candidate |
| 0004 | Compact general foundation | 100-150M | Rung 0003 data plus broad high-quality educational, textbook, reference, math, and technical text | 20/40/80/160 core TPP checkpoints, then review | Correct basic factual answers, coherent paragraph explanations, basic technical knowledge, low repetition | Planned |
| 0005 | Small assistant | Derived from passed 0004 | Curated supervised QA/instructions with an explicit interaction contract and assistant EOS | Short post-training stages with retention checks | Direct answers, instruction following, role isolation, EOS, and no material foundation regression | Planned |
| 0006 | Elementary reasoning and technical assistant | 250-400M, justified by 0004 scaling evidence | Passed foundation mix plus licensed math, science, and code/technical data | Scaling prediction plus matched checkpoints | Multi-step elementary reasoning, basic code explanation, and retained general/assistant behavior | Planned |

The ranges are decision boundaries, not invitations to sweep architectures.
Each rung starts from the smallest architecture supported by evidence from the
previous rung. Capacity grows only when more tokens stop producing useful
capability improvement or the rung's task clearly exceeds the model's capacity.

### Common generation gates

Every language-capability rung uses a fixed deterministic suite and a fixed
sampled suite. At minimum, promotion requires:

- no obvious repetition loop in more than 1 of 10 deterministic probes;
- at least 90% EOS termination in the long-horizon suite;
- no prompt-role or document-boundary leakage;
- no empty, malformed, or unrelated answer in more than 1 of 10 probes; and
- no regression on the prior rung's retained-capability probes.

Rung-specific correctness gates are defined with its candidate compose and
evaluation fixture before the first training run. For rung 0003 the initial
target is at least 8/10 relevant grammatical answers and 6/10 exact elementary
facts or arithmetic answers. For rung 0004 it is at least 12/15 relevant
answers and 8/10 exact basic factual answers. These targets may be revised only
after a written capacity analysis, never merely because a run failed.

## Data required to grow the ladder

The next step is not another arbitrary mixture. Each addition must fill a
declared capability gap and pass WALDO's provenance and redistribution policy.

| Priority | Corpus or view | Rung | Purpose | Action |
| --- | --- | --- | --- | --- |
| 1 | PLOS Genetics view | 0003 | Grounded volume in one bounded scientific domain | Filter the existing journal metadata; do not duplicate shards |
| 1 | Cell/genetics Wikimedia view | 0003 | Definitions and explanatory prose in the same domain | Filter existing wiki/title metadata using the frozen topic-title set |
| 1 | Full TinyStories | 0003 | Preserve proven grammar, entity, causality, and EOS behavior | Ingest all source shards when available; current reference is explicitly a subset |
| 2 | Open Textbook Library permissive subset | 0004 | Coherent instructional chapters across subjects | Ingest only books with compatible per-item licenses |
| 2 | Existing PressBooks, LibreTexts, Wikimedia, and Stack-Edu views | 0004/0006 | Textbooks, reference knowledge, and technical language | Improve selection and metadata filters before adding volume |
| 2 | DCLM-Edu high-score subset | 0004 | Web-scale educational language selected for small models | Complete source-level redistribution review before ingestion |
| 2 | FineMath-4+ | 0006 | High-quality mathematical prose and problems | Complete source-level redistribution review before ingestion |

Synthetic data may supplement a rung only when its generator, prompt, seed,
license, and derivation are recorded. It must not replace independent held-out
evaluation or dominate the mixture merely because it lowers loss.

## Metadata filtering and reingestion

Adding metadata filtering does **not automatically require reingestion**.
Canonical shards already have a `meta` JSON column. If the original fetcher
mapped the field we need into that column, WALDO can add record-level metadata
predicates and reuse the existing immutable shards. Wikimedia and Cosmopedia
already preserve useful namespace/title and audience/format metadata,
respectively.

Reingestion is required only when the necessary source metadata was never
written into the canonical shard, or when the canonical text itself must
change. Before any reingest:

1. inspect a sample of the existing shard's `meta` values;
2. prove that the desired field is absent or unusable; and
3. update the fetcher mapping and provenance contract.

WALDO's current compose filters do not expose arbitrary `meta` predicates, so
metadata-aware filtering is an implementation prerequisite for rung 0003. It
should extend selection over existing shards; it should not create duplicate
filtered corpora.

## Candidate workflow

For each next rung:

1. define the capability, frozen probes, gates, data prerequisites, and maximum
   review budget;
2. add or select only the corpora required for that capability;
3. freeze the architecture, tokenizer, mixture, held-out partition, and seed;
4. train and score matched 20/40/80/160-core-TPP checkpoints without changing
   those controls;
5. diagnose failure as data, optimization, exposure, or capacity using the
   loss curve and per-probe results; and
6. add the successful rung directory with `README.md`, `compose.yaml`,
   `compose_test.go`, and `evaluate.sh` only after every gate passes.

The immediate work is rung 0003: validate the preflight size and exposure of
the PLOS Genetics and title-selected molecular-biology Wikimedia views, then
train the fixed 30M candidate against the frozen single-domain evaluation
fixture.
