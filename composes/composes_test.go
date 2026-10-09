// Copyright (c) 2026 OpenWALDO Project contributors
// Copyright (c) 2026 CtrlIQ, Inc.
// Copyright (c) 2026 Gregory M. Kurtzer
// SPDX-License-Identifier: Apache-2.0

package composes_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/openwaldo/waldo/internal/corpus"
	"github.com/openwaldo/waldo/internal/model"
	"github.com/openwaldo/waldo/internal/training"
)

var promotedFiles = []string{
	"0001-tiny-shakespeare.yaml",
	"0002-tinystories-bpe.yaml",
}

func TestModelComposeGuideNamesEverySchemaField(t *testing.T) {
	guide, err := os.ReadFile(filepath.Join("..", "docs", "MODEL-COMPOSE.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []any{
		model.Compose{}, model.ComposeBase{}, model.Architecture{}, model.Tokenizer{}, model.TokenizerTraining{}, model.Stage{}, model.CorpusSelection{}, corpus.RecordFilter{}, corpus.ValueFilter{}, corpus.DateFilter{}, training.ConversationTransform{}, training.Parameters{},
	} {
		typeOf := reflect.TypeOf(value)
		for index := 0; index < typeOf.NumField(); index++ {
			name := strings.Split(typeOf.Field(index).Tag.Get("yaml"), ",")[0]
			if name == "" || name == "-" {
				continue
			}
			if !strings.Contains(string(guide), "`"+name+"`") && !strings.Contains(string(guide), name+":") {
				t.Errorf("%s.%s YAML field %q is absent from MODEL-COMPOSE.md", typeOf.Name(), typeOf.Field(index).Name, name)
			}
		}
	}
}

func TestOnlyPromotedComposesArePersisted(t *testing.T) {
	var files []string
	err := filepath.WalkDir(".", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}
		files = append(files, strings.TrimPrefix(path, "./"))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(files, promotedFiles) {
		t.Fatalf("persisted composes = %v, want promoted ladder %v", files, promotedFiles)
	}
}

func TestPromotedComposesResolveIntoTrainingContracts(t *testing.T) {
	for _, file := range promotedFiles {
		t.Run(file, func(t *testing.T) {
			compose := loadCompose(t, file)
			if compose.Base != nil || compose.Interaction.Template != "" || len(compose.Stages) != 1 {
				t.Fatalf("promoted compose is not fresh one-stage pretraining: %+v", compose)
			}
			stage := compose.Stages[0]
			resolved, err := stage.ResolvePlanningParameters()
			if err != nil {
				t.Fatal(err)
			}
			if resolved.Profile != stage.Parameters.Profile || resolved.BatchSize != stage.Parameters.BatchSize || resolved.SequenceLength != stage.Parameters.SequenceLength || resolved.LearningRate != stage.Parameters.LearningRate || resolved.Seed != stage.Parameters.Seed {
				t.Fatalf("direct settings were not preserved: raw=%+v resolved=%+v", stage.Parameters, resolved)
			}
			if resolved.Evaluation == nil {
				t.Fatal("promoted compose has no held-out evaluation policy")
			}
		})
	}
}

func TestPromotedReferenceForecastsAndControls(t *testing.T) {
	want := []struct {
		parameters uint64
		tokens     int64
		corpus     string
	}{
		{10721280, 81920000, "core/reference/tiny-shakespeare"},
		{32777728, 613613568, "core/synthetic/tinystories-reference"},
	}
	for index, file := range promotedFiles {
		compose := loadCompose(t, file)
		forecast, err := model.ForecastCompose(compose)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		if forecast.ApproximateParameters != want[index].parameters || forecast.PlannedTokens != want[index].tokens {
			t.Fatalf("%s forecast = %d parameters/%d tokens, want %+v", file, forecast.ApproximateParameters, forecast.PlannedTokens, want[index])
		}
		if got := corpusPaths(compose.Stages[0].Corpora); !reflect.DeepEqual(got, []string{want[index].corpus}) {
			t.Fatalf("%s corpus = %v", file, got)
		}
	}

	stories := loadCompose(t, promotedFiles[1])
	architecture := stories.Architecture
	if architecture.ContextTokens != 512 || architecture.VocabularySize != 4096 || architecture.HiddenSize != 512 || architecture.IntermediateSize != 1536 || architecture.Layers != 9 || architecture.AttentionHeads != 8 || architecture.KeyValueHeads != 8 || architecture.Dropout != 0.2 || architecture.QKNormalization || architecture.Initialization != "depth-scaled" || !architecture.TieEmbeddings || architecture.ParameterDType != "float32" {
		t.Fatalf("TinyStories reference architecture changed: %+v", architecture)
	}
	tokenizer := architecture.Tokenizer.Training
	if tokenizer == nil || tokenizer.Algorithm != model.TokenizerAlgorithmByteBPEV1 || tokenizer.SampleBytes != 134217728 || tokenizer.Seed != 42 || tokenizer.DistributionPolicy != corpus.DistributionPolicyDistributable || !reflect.DeepEqual(corpusPaths(tokenizer.Corpora), []string{"core/common-pile/pressbooks"}) {
		t.Fatalf("TinyStories reference tokenizer changed: %+v", tokenizer)
	}
	parameters := stories.Stages[0].Parameters
	if parameters.Profile != "causal-pretrain-shuffled" || parameters.Tokens != 613611520 || parameters.BatchSize != 12 || parameters.GradientAccumulation != 3 || parameters.SequenceLength != 512 || parameters.LearningRate != 0.0006 || parameters.Optimizer != "adamw" || parameters.Schedule != "cosine" || parameters.Seed != 42 || parameters.EvaluationSelection != "lowest-sha256-v1" {
		t.Fatalf("TinyStories reference training controls changed: %+v", parameters)
	}
}

func TestREADMEEnforcesPromoteOnlyPolicy(t *testing.T) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"only reproducible model recipes that passed every",
		"failed or sub-par candidate is removed",
		"BPB = held_out_loss / (held_out_bytes_per_token * ln(2))",
		"A lower aggregate loss alone is not a promotion",
		"Rung 0001: Tiny Shakespeare reference",
		"Rung 0002: TinyStories BPE reference",
		"passed every promotion gate",
	} {
		if !strings.Contains(string(content), required) {
			t.Errorf("README does not contain %q", required)
		}
	}
}

func loadCompose(t *testing.T, path string) model.Compose {
	t.Helper()
	compose, _, err := model.LoadCompose(path)
	if err != nil {
		t.Fatal(err)
	}
	return compose
}

func corpusPaths(selections []model.CorpusSelection) []string {
	paths := make([]string, len(selections))
	for index, selection := range selections {
		paths[index] = selection.Path
	}
	return paths
}
