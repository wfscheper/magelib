// Copyright © 2020 The Stentor Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build mage
// +build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	// mage:import
	"github.com/wfscheper/magelib"
)

var (
	// Default mage target
	Default = All

	getGolangciLint = magelib.GetGolangciLint("v1.54.1")
	getGoreleaser   = magelib.GetGoreleaser("v1.20.0")
	getGotagger     = magelib.GetGotagger("v0.9.0")
	getGotestsum    = magelib.GetGotestsum("v1.10.0")
	getStentor      = magelib.GetStentor("v0.3.0")
)

func init() {
	magelib.ChangelogDeps = []interface{}{
		func(ctx context.Context) error { return getStentor(ctx) },
	}
	magelib.LintDeps = []interface{}{
		func(ctx context.Context) error { return getGolangciLint(ctx) },
	}
	magelib.ReleaseDeps = []interface{}{
		func(ctx context.Context) error { return getGoreleaser(ctx) },
	}
	magelib.TestDeps = []interface{}{
		func(ctx context.Context) error { return getGotestsum(ctx) },
	}
	magelib.VersionDeps = []interface{}{
		func(ctx context.Context) error { return getGotagger(ctx) },
	}

	magelib.ProjectTools = magelib.ToolMap{
		magelib.ModuleGolangciLint: getGolangciLint,
		magelib.ModuleGoreleaser:   getGoreleaser,
		magelib.ModuleGotagger:     getGotagger,
		magelib.ModuleGotestsum:    getGotestsum,
		magelib.ModuleStentor:      getStentor,
	}
}

// All runs format, lint, vet, build, and test targets
func All(ctx context.Context) {
	mg.SerialCtxDeps(ctx, magelib.Go.Lint, magelib.Go.Build, magelib.Go.Test)
}
