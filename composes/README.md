# Validated capability ladder

This directory is a catalog of reproducible models that passed every declared
gate. It is not a roadmap and it is not an experiment log. At present the
validated ladder ends at rung 0002.

Only a passed recipe may have a numbered rung directory. Candidate composes
belong on an experiment branch or in a temporary workspace. A failed candidate
is removed; its compose, telemetry, evaluation output, and decision remain in
Git and the persisted model artifact.

## Current ladder

| Rung | Proven capability | Model | Evidence | Status |
| --- | --- | --- | --- | --- |
| [0001](0001-tiny-shakespeare/) | End-to-end learning and checkpoint-selection control on one continuous literary distribution | 10.7M parameters, byte tokenizer | Plausible play continuation; best pre-overtraining checkpoint selected and reloaded exactly | Passed |
| [0002](0002-tinystories-bpe/) | Narrow simple-story language model | 30.7M core parameters, 4K byte BPE | Entity/causal retention, bounded repetition, and 8/8 long-horizon EOS | Passed |

There is deliberately no rung 0003. A cell-biology candidate previously
occupying that number failed every behavioral promotion gate and was removed.

## What the failed expansion proved

The failed `cell-biology-30m-80tpp-01` run was numerically healthy: it consumed
2.454B tokens, improved held-out loss from 8.4060 to 2.1062, selected and
exactly reloaded its best checkpoint, and finished four-GPU training normally.
It nevertheless produced approximately 1/10 correct factual continuations,
failed the explanation and retained-story gates, entered repetition loops on
nearly every factual probe, and emitted EOS on 0/10 long factual generations.

That result does **not** establish that WALDO's transformer implementation is
broken. The same architecture family learns TinyStories, and independent
numerical and memorization controls pass. It establishes that the proposed
capability experiment was unsupported:

- a TinyStories-proven 30.7M-core architecture was assumed to support factual
  domain acquisition without a matched public or internal baseline;
- elementary textbook facts were evaluated after training mostly on research
  articles, which optimize a different continuation distribution;
- only 269.5M unique targets were replayed to fill 2.454B training tokens;
- PLOS Genetics was repeated about 15.6 times and the small Wikimedia view
  about 35.2 times; and
- raw loss on the training distribution was treated as evidence for semantic
  capability even though no such relationship had been calibrated.

The failure therefore rejects that experiment. It does not authorize another
mixture, a larger token budget, or a larger model.

## Why TinyStories succeeded and factual rungs did not

TinyStories intentionally constrains vocabulary, syntax, subject matter, and
document form. Its evaluation asks the model to continue the same distribution
it learned. The TinyStories paper specifically demonstrates that very small
models can become coherent under this simplified distribution; it does not
claim that the same capacity can absorb broad or research-grade factual text.

Chinchilla's often-quoted token/parameter relationship is a compute-allocation
result, not a capability threshold. Conversely, very high token/parameter
ratios do not create capability when a small unique corpus is replayed many
times. Work on data-constrained scaling finds that repeated data remains useful
for a limited number of epochs, after which the value of more compute decays.
The failed cell-biology run exceeded that evidence for its smallest views.

References:

- [TinyStories](https://arxiv.org/abs/2305.07759)
- [Training Compute-Optimal Large Language Models](https://arxiv.org/abs/2203.15556)
- [Scaling Data-Constrained Language Models](https://arxiv.org/abs/2305.16264)
- [BabyLM](https://arxiv.org/abs/2301.11796)

## Hard stop before rung 0003

No new foundation or domain compose may be added or trained until all of the
following exist:

1. **A declared track.** The proposal is explicitly trainer validation,
   research, a from-scratch foundation model, or product post-training.
2. **A quantitative evaluator.** The report uses immutable evaluation data,
   low-variance log-likelihood or multiple-choice metrics, automatic EOS and
   repetition diagnostics, and confidence intervals where applicable. Raw
   generations are supporting diagnostics, not the primary score.
3. **Independent baselines.** The same evaluator is run on an untrained model,
   the previous WALDO rung, and at least one public model at comparable size
   and exposure. A proposed threshold must be attainable by an identified
   baseline.
4. **A data-support report.** Before training, the exact derived view reports
   unique tokens, effective source passes, duplicates, artifacts, tokenizer
   fertility, evaluation overlap, and direct support for every demanded
   capability.
5. **A scaling pilot.** Existing checkpoints or a deliberately tiny run show
   that the chosen metric improves with exposure or capacity. The pilot changes
   one causal variable and has a written early-stop rule.
6. **A bounded qualification run.** Architecture, tokenizer, data, optimizer,
   evaluator, seed, checkpoints, maximum budget, and promotion gate are frozen.
   Failure ends the candidate; it does not trigger an improvised extension.

Until these prerequisites pass, the correct next task is measurement and
baseline implementation—not training.

## Required evidence packet

Every proposed rung must have a reviewable packet before its compose is
eligible to run:

```text
track and intended use
one falsifiable hypothesis
parent or public baseline
one changed variable
immutable data and tokenizer identities
data-support and contamination report
evaluation BOM and scoring implementation
expected metric range with supporting evidence
pilot budget and early-stop rule
maximum qualification budget
promotion and rejection thresholds
```

Claims such as "30M should know elementary biology" or "80 tokens per
parameter should be enough" are not evidence and cannot populate this packet.

## Parameter and loss accounting

WALDO reports core parameters separately from token input/output weights.
Training exposure, unique source exposure, and effective source passes are
separate quantities. A tokens-per-parameter ratio is descriptive; it is never
a promotion criterion by itself.

Raw token loss is comparable only with the same tokenizer and held-out bytes.
When tokenizers differ, use exact bits per UTF-8 byte:

```text
BPB = held_out_nll_sum / (held_out_utf8_bytes * ln(2))
```

BPB measures distribution fit, not factual correctness, reasoning, stopping,
or usefulness. Every capability claim requires an external task metric.

## Architecture policy

The current evidence neither convicts nor validates the 30.7M-core architecture
for factual language modeling. It is validated only for the TinyStories rung.
Architecture changes must therefore be comparative experiments at fixed data,
tokenizer, FLOPs, optimizer, and evaluation—not guesses based on parameter
count alone.

Depth, width, attention layout, vocabulary allocation, context, batch size,
learning rate, weight decay, clipping, and packing are all controlled variables.
At least two matched configurations and multiple seeds are required before an
architecture decision is promoted.

## Next authorized work

The next work is the validation backlog in
[`docs/TRAINING-ROBUSTNESS-PLAN.md`](../docs/TRAINING-ROBUSTNESS-PLAN.md):

1. finish fixed-global-batch one/two/four-GPU and resume equivalence;
2. implement a versioned log-likelihood evaluation report;
3. evaluate matched public small-model checkpoints with that report;
4. build corpus-support, duplication, artifact, and overlap reports; and
5. use existing WALDO artifacts to test whether any metric has a defensible
   scaling trend before proposing another run.

Only after those measurements support a bounded capability may a candidate
rung 0003 be designed. Only after that candidate passes every frozen gate may
its directory enter this ladder.

## Rung directory contract

Each promoted rung contains:

```text
README.md       purpose, immutable evidence, gates, and commands
compose.yaml    exact promoted training recipe
compose_test.go guards for architecture, data, and training contract
evaluate.sh     fixed behavioral diagnostics
```

The authoritative project-wide validation sequence is documented in
[`docs/TRAINING-ROBUSTNESS-PLAN.md`](../docs/TRAINING-ROBUSTNESS-PLAN.md).
