// Copyright (c) 2026 OpenWALDO Project contributors
// SPDX-License-Identifier: Apache-2.0

package elementaryfoundation_test

import (
	"reflect"
	"testing"

	"github.com/openwaldo/waldo/internal/model"
)

func TestCandidateContract(t *testing.T) {
	compose, _, err := model.LoadCompose("compose.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if compose.Base != nil || compose.Interaction.Template != "" || len(compose.Stages) != 1 {
		t.Fatalf("candidate is not fresh one-stage pretraining: %+v", compose)
	}

	forecast, err := model.ForecastCompose(compose)
	if err != nil {
		t.Fatal(err)
	}
	if forecast.ApproximateParameters != 32777728 || forecast.PlannedTokens != 2454448128 {
		t.Fatalf("forecast = %d parameters/%d tokens", forecast.ApproximateParameters, forecast.PlannedTokens)
	}

	stage := compose.Stages[0]
	if stage.Name != "elementary-pretrain" || len(stage.Corpora) != 2 {
		t.Fatalf("stage changed: %+v", stage)
	}
	if stage.Corpora[0].Path != "core/synthetic/cosmopedia-v2" || stage.Corpora[0].Weight == nil || *stage.Corpora[0].Weight != 7 || stage.Corpora[1].Path != "core/synthetic/tinystories-reference" || stage.Corpora[1].Weight == nil || *stage.Corpora[1].Weight != 3 {
		t.Fatalf("mixture changed: %+v", stage.Corpora)
	}
	metadata := stage.Corpora[0].Filter.Metadata
	if !reflect.DeepEqual(metadata["audience"].Include, []string{"middle_school_students"}) || !reflect.DeepEqual(metadata["format"].Include, []string{"textbook"}) {
		t.Fatalf("Cosmopedia view changed: %+v", metadata)
	}

	architecture := compose.Architecture
	if architecture.ContextTokens != 512 || architecture.VocabularySize != 4096 || architecture.HiddenSize != 512 || architecture.IntermediateSize != 1536 || architecture.Layers != 9 || architecture.AttentionHeads != 8 || architecture.KeyValueHeads != 8 || architecture.Dropout != 0.2 || architecture.QKNormalization || architecture.Initialization != "depth-scaled" || !architecture.TieEmbeddings || architecture.ParameterDType != "float32" {
		t.Fatalf("architecture changed: %+v", architecture)
	}
	parameters := stage.Parameters
	if parameters.Profile != "causal-pretrain-weighted" || parameters.Tokens != 2454446080 || parameters.BatchSize != 12 || parameters.GradientAccumulation != 3 || parameters.SequenceLength != 512 || parameters.LearningRate != 0.0006 || parameters.Optimizer != "adamw" || parameters.Schedule != "cosine" || parameters.Seed != 42 || parameters.CheckpointEvery == nil || *parameters.CheckpointEvery != 9987 || parameters.EvaluateEvery == nil || *parameters.EvaluateEvery != 9987 || parameters.EvaluationSelection != "lowest-sha256-v1" {
		t.Fatalf("training controls changed: %+v", parameters)
	}
}
