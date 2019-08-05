package main

import (
	"io/ioutil"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

const plantumlHeader = `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

@startuml
hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
}`

func TestGenerateIntegrations(t *testing.T) {
	m, _ := parse.Parse("demo/simple/sysl-ints.sysl", "../../")
	require.NotNil(t, m)

	stmt := &sysl.Statement{}
	args := &Args{"", "Project", false, false}
	apps := []string{"System1", "IntegratedSystem", "System2"}
	highlights := MakeStrSet("IntegratedSystem", "System1", "System2")
	s1 := AppElement{"IntegratedSystem", "integrated_endpoint_1"}
	t1 := AppElement{"System1", "endpoint"}
	dep1 := AppDependency{
		Self:      s1,
		Target:    t1,
		Statement: stmt,
	}
	deps := []AppDependency{
		dep1,
	}
	endpt := &sysl.Endpoint{
		Name: "_",
		Stmt: []*sysl.Statement{
			{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "IntegratedSystem"}}},
			{Stmt: &sysl.Statement_Action{Action: &sysl.Action{Action: "System1"}}},
		},
	}
	intsParam := &IntsParam{apps, highlights, deps, m.GetApps()["Project"], endpt}
	r := GenerateView(args, intsParam, m)

	require.NotNil(t, r)

	expected := plantumlHeader + `
[IntegratedSystem] as _0 <<highlight>>
[System1] as _1 <<highlight>>
_0 --> _1
@enduml`

	assert.Equal(t, expected, r)
}

type intsArg struct {
	rootModel string
	title     string
	output    string
	project   string
	filter    string
	modules   string
	exclude   []string
	clustered bool
	epa       bool
}

func comparePUML(t *testing.T, expected, actual map[string]string) {
	for name, goldenFile := range expected {
		golden, err := ioutil.ReadFile(goldenFile)
		assert.Nil(t, err)
		if string(golden) != actual[name] {
			err := ioutil.WriteFile("tests/"+name+".puml", []byte(actual[name]), 0777)
			assert.Nil(t, err)
		}
		assert.Equal(t, string(golden), actual[name])
	}

	// Then
	assert.Equal(t, len(expected), len(actual))
}

func TestGenerateIntegrationsWithTestFile(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "indirect_1.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	expected := map[string]string{
		"all.png":            "tests/indirect_1-all-golden.puml",
		"indirect_arrow.png": "tests/indirect_1-indirect_arrow-golden.puml",
		"my_callers.png":     "tests/indirect_1-my_callers-golden.puml",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithTestFileAndFilters(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_test.sysl",
		output:    "%(epname).png",
		project:   "Project",
		filter:    "test",
	}
	expected := map[string]string{}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	// Then
	assert.Equal(t, expected, result)
}

func TestGenerateIntegrationsWithImmediatePredecessors(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_immediate_predecessors_test.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}
	expected := map[string]string{
		"immediate_predecessors.png": "tests/immediate_predecessors-golden.puml",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithExclude(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_excludes_test.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	expected := map[string]string{
		"excludes.png": "tests/excludes-golden.puml",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithPassthrough(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_passthrough_test.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	expected := map[string]string{
		"passthrough.png": "tests/passthrough-golden.puml",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	// Then
	comparePUML(t, expected, result)
}

func TestDoGenerateIntegrations(t *testing.T) {
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "indirect_1.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}
	argsData := []string{"sysl", "ints", "--root", args.rootModel, "-o", args.output, "-j", args.project, args.modules}
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	configureCmdlineForIntgen(sysl)
	selectedCommand, err := sysl.Parse(argsData[1:])
	assert.Nil(t, err, "Cmd line parse failed for sysl ints")
	assert.Equal(t, selectedCommand, "ints")
}

func TestGenerateIntegrationsWithCluster(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_with_cluster.sysl",
		output:    "%(epname).png",
		project:   "Project",
		clustered: true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"cluster.png": "tests/cluster-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithEpa(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_with_epa.sysl",
		output:    "%(epname).png",
		project:   "Project",
		epa:       true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"epa.png": "tests/epa-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithIndirectArrow(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "indirect_2.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"all_indirect_2.png":  "tests/all_indirect_2-golden.puml",
		"no_passthrough.png":  "tests/no_passthrough-golden.puml",
		"passthrough_b.png":   "tests/passthrough_b-golden.puml",
		"passthrough_c.png":   "tests/passthrough_c-golden.puml",
		"passthrough_d.png":   "tests/passthrough_d-golden.puml",
		"passthrough_c_e.png": "tests/passthrough_c_e-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithRestrictBy(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_with_restrict_by.sysl",
		output:    "%(epname).png",
		project:   "Project",
		epa:       true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"with_restrict_by.png":    "tests/with_restrict_by-golden.puml",
		"without_restrict_by.png": "tests/without_restrict_by-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithFilter(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_with_filter.sysl",
		output:    "%(epname).png",
		project:   "Project",
		filter:    "matched",
	}

	expected := map[string]string{
		"matched.png": "tests/matched-golden.puml",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationWithOrWithoutPassThrough(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_with_or_without_passthrough.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"with_passthrough.png":    "tests/with_passthrough-golden.puml",
		"without_passthrough.png": "tests/without_passthrough-golden.puml",
		"with_systema.png":        "tests/with_systema-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)

}

func TestPassthrough2(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "passthrough_1.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"pass_1_all.png":   "tests/pass_1_all-golden.puml",
		"pass_1_sys_a.png": "tests/pass_1_sys_a-golden.puml",
		"pass_b.png":       "tests/pass_b-golden.puml",
		"pass_b_c.png":     "tests/pass_b_c-golden.puml",
		"pass_f.png":       "tests/pass_f-golden.puml",
		"pass_D.png":       "tests/pass_D-golden.puml",
		"pass_e.png":       "tests/pass_e-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}

func TestGenerateIntegrationsWithPubSub(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "integration_with_pubsub.sysl",
		output:    "%(epname).png",
		project:   "Project",
		epa:       true,
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"pubsub.png": "tests/pubsub-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}

func TestAllStmts(t *testing.T) {
	// Given
	args := &intsArg{
		rootModel: "./tests/",
		modules:   "ints_stmts.sysl",
		output:    "%(epname).png",
		project:   "Project",
	}

	// When
	result, err := GenerateIntegrationsWithParams(args.rootModel, args.title, args.output,
		args.project, args.filter, args.modules, args.exclude, args.clustered,
		args.epa, "warn", false)
	require.NoError(t, err)

	expected := map[string]string{
		"all_stmts.png": "tests/all_stmts-golden.puml",
	}

	// Then
	comparePUML(t, expected, result)
}
