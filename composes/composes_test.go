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

var promotedRungs = []string{
	"0001-tiny-shakespeare",
	"0002-tinystories-bpe",
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

func TestOnlyPromotedRungBundlesArePersisted(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	want := append(append([]string{}, promotedRungs...), "README.md", "composes_test.go")
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("compose directory entries = %v, want %v", names, want)
	}

	for _, rung := range promotedRungs {
		for _, required := range []string{"README.md", "compose.yaml", "compose_test.go", "evaluate.sh"} {
			path := filepath.Join(rung, required)
			info, err := os.Stat(path)
			if err != nil {
				t.Errorf("promoted rung %s is missing %s: %v", rung, required, err)
				continue
			}
			if info.IsDir() {
				t.Errorf("promoted rung component %s is a directory", path)
			}
			if required == "evaluate.sh" && info.Mode()&0o111 == 0 {
				t.Errorf("promoted rung evaluator %s is not executable", path)
			}
		}
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
		"0001-tiny-shakespeare",
		"0002-tinystories-bpe",
	} {
		if !strings.Contains(string(content), required) {
			t.Errorf("README does not contain %q", required)
		}
	}
}
