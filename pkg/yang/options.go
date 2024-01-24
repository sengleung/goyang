// Copyright 2017 Google Inc.
//
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

package yang

// Options defines the options that should be used when parsing YANG modules,
// including specific overrides for potentially problematic YANG constructs.
type Options struct {
	// IgnoreSubmoduleCircularDependencies specifies whether circular dependencies
	// between submodules. Setting this value to true will ensure that this
	// package will explicitly ignore the case where a submodule will include
	// itself through a circular reference.
	IgnoreSubmoduleCircularDependencies bool
	// StoreUses controls whether the Uses field of each YANG entry should be
	// populated. Setting this value to true will cause each Entry which is
	// generated within the schema to store the logical grouping from which it
	// is derived.
	StoreUses bool
	// DeviateOptions contains options for how deviations are handled.
	DeviateOptions DeviateOptions
	// StatementOptions contains options for how statements are handled.
	StatementOptions StatementOptions
	// ModuleOptions contains options for how modules are handled.
	ModuleOptions ModuleOptions
}

// DeviateOptions contains options for how deviations are handled.
type DeviateOptions struct {
	// IgnoreDeviateNotSupported indicates to the parser to retain nodes
	// that are marked with "deviate not-supported". An example use case is
	// where the user wants to interact with different targets that have
	// different support for a leaf without having to use a second instance
	// of an AST.
	IgnoreDeviateNotSupported bool
}

// IsDeviateOpt ensures that DeviateOptions satisfies the DeviateOpt interface.
func (DeviateOptions) IsDeviateOpt() {}

// DeviateOpt is an interface that can be used in function arguments.
type DeviateOpt interface {
	IsDeviateOpt()
}

func hasIgnoreDeviateNotSupported(opts []DeviateOpt) bool {
	for _, o := range opts {
		if opt, ok := o.(DeviateOptions); ok {
			return opt.IgnoreDeviateNotSupported
		}
	}
	return false
}

// StatementOptions contains options for how statements are handled.
type StatementOptions struct {
	// ExcludeStatements is a list of statement keywords that are ignored
	// when parsing Statements.
	ExcludeStatements []string
	// LatestRevisionOnly indicates whether to parse the latest revision statement
	// only and ignore the previous revision statements.
	LatestRevisionOnly bool
}

// IsStatementOpt ensures that StatementOptions satisfies the StatementOpt interface.
func (StatementOptions) IsStatementOpt() {}

// StatementOpt is an interface that can be used in function arguments.
type StatementOpt interface {
	IsStatementOpt()
}

// statementOptions contains options for how statements are handled.
type statementOptions struct {
	excludeStatements  map[string]struct{}
	latestRevisionOnly bool
}

func newStatementOptions() *statementOptions {
	return &statementOptions{
		excludeStatements: make(map[string]struct{}),
	}
}

func (opts *statementOptions) addExcludeStatements(statementOpts ...StatementOpt) {
	for _, o := range statementOpts {
		if opt, ok := o.(StatementOptions); ok {
			for _, keyword := range opt.ExcludeStatements {
				opts.excludeStatements[keyword] = struct{}{}
			}
		}
	}
}

func (opts *statementOptions) setLatestRevisionOnly(statementOpts ...StatementOpt) {
	for _, o := range statementOpts {
		if opt, ok := o.(StatementOptions); ok {
			opts.latestRevisionOnly = opt.LatestRevisionOnly
		}
	}
}

// includeStatement returns true if the statement should be included.
func (opts *statementOptions) includeStatement(keyword string) bool {
	if opts == nil || len(opts.excludeStatements) == 0 {
		return true
	}
	_, found := opts.excludeStatements[keyword]
	return !found
}

// ModuleOptions contains options for how modules are handled.
type ModuleOptions struct {
	// IncludeOnlySources is a list of statement keywords that link the Statement
	// to the Source field of the module Node. Statements not in the keywords list
	// are not linked. If the list is empty, all statements are linked.
	IncludeOnlySources []string
}

// IsModuleOpt ensures that ModuleOptions satisfies the ModuleOpt interface.
func (ModuleOptions) IsModuleOpt() {}

// ModuleOpt is an interface that can be used in function arguments.
type ModuleOpt interface {
	IsModuleOpt()
}

// moduleOptions contains options for how modules are handled.
type moduleOptions struct {
	includeOnlySources map[string]struct{}
}

func newModuleOptions() *moduleOptions {
	return &moduleOptions{
		includeOnlySources: make(map[string]struct{}),
	}
}

func (opts *moduleOptions) addIncludeOnlySources(moduleOpts ...ModuleOpt) {
	for _, o := range moduleOpts {
		if opt, ok := o.(ModuleOptions); ok {
			for _, keyword := range opt.IncludeOnlySources {
				opts.includeOnlySources[keyword] = struct{}{}
			}
		}
	}
}

// setSourceStatement returns true if the source statement should be set.
func (opts *moduleOptions) setSourceStatement(keyword string) bool {
	if opts == nil || len(opts.includeOnlySources) == 0 {
		return true
	}
	_, found := opts.includeOnlySources[keyword]
	return found
}
