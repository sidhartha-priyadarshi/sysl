package main

import (
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/parser"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/validate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Node can be string or node
type Node []interface{}

type CodeGenOutput struct {
	filename string
	output   Node
}

func getKeyFromValueMap(v *sysl.Value, key string) *sysl.Value {
	if m := v.GetMap(); m != nil {
		return m.Items[key]
	}
	return nil
}

func processChoice(
	g *parser.Grammar,
	obj *sysl.Value,
	choice *parser.Choice,
	logger *logrus.Logger,
) Node {
	var result Node

	for i, seq := range choice.Sequence {
		seqResult := Node{}
		fullScan := true
		for _, term := range seq.Term {
			switch x := term.Atom.Union.(type) {
			// String tokens dont have quantifiers
			case *parser.Atom_String_:
				seqResult = append(seqResult, x.String_)
			case *parser.Atom_Rulename:
				var ruleResult interface{}

				minc, maxc := parser.GetTermMinMaxCount(term)
				v := getKeyFromValueMap(obj, x.Rulename.Name)

				// raise error if required
				//  i.e.  no quantifier or +
				//        and missing from obj map
				if minc > 0 && v == nil {
					fullScan = false
					break
				}

				// skip if rule has
				//    quantifier == * or ?
				//    and does not exist in obj map
				if minc == 0 && v == nil {
					continue
				}

				if maxc > 1 {
					var valueList []*sysl.Value
					switch vv := v.Value.(type) {
					case *sysl.Value_List_:
						valueList = vv.List.Value
					case *sysl.Value_Set:
						valueList = vv.Set.Value
					default:
						logger.Warnf("Expecting a collection type, got %T for rule %s", vv, x.Rulename.Name)
						fullScan = false
					}
					ruleInstances := Node{}

					for _, valueItem := range valueList {
						// Drill down the rule
						node := processRule(g, valueItem, x.Rulename.Name, logger)
						// Check post-conditions
						if len(node) == 0 {
							fullScan = false
							break
						}
						ruleInstances = append(ruleInstances, node)
					}
					ruleResult = ruleInstances
				} else { // maxc == 1
					// Drill down the rule
					if v.GetList() != nil || v.GetSet() != nil {
						logger.Warnf("Got List or Set instead of map")
					}
					node := processRule(g, v, x.Rulename.Name, logger)
					// Check post-conditions
					if len(node) == 0 {
						logger.Warnf("could not process rule: ( %s )", x.Rulename.Name)
						fullScan = false
						break
					}
					if s, ok := node[0].(string); ok && len(node) == 1 {
						ruleResult = s
					} else {
						ruleResult = node
					}
				}
				seqResult = append(seqResult, ruleResult)
			case *parser.Atom_Choices:
				// minc, maxc := parser.GetMinMaxCount(term)
				node := processChoice(g, obj, x.Choices, logger)
				if len(node) == 0 {
					logger.Warnf("could not process Choice\n")
					fullScan = false
					break
				}
				seqResult = append(seqResult, node)
			default:
				logger.Warningf("processChoice: choice %d : %T", i, x)
				panic("Unexpected atom type")
			}
			if !fullScan {
				break
			}
		}
		if fullScan {
			result = append(result, seqResult)
		}
	}
	return result
}

func processRule(g *parser.Grammar, obj *sysl.Value, ruleName string, logger *logrus.Logger) Node {
	var str string
	if x := obj.GetMap(); x != nil {
		for key := range x.Items {
			str += key + ", "
		}
	}
	// logrus.Debugf("processRule: %s, obj keys (%s)", ruleName, str)
	rule := g.Rules[ruleName]
	if rule == nil {
		root := Node{}
		if eval.IsCollectionType(obj) {
			return nil
		}
		// Should we convert int and bools to string and return?
		return append(root, obj.GetS())
	}
	root := processChoice(g, obj, rule.Choices, logger)
	return root
}

func readGrammar(filename, grammarName, startRule string) (*parser.Grammar, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parser.ParseEBNF(string(dat), grammarName, startRule), nil
}

/*
params[0] - modelAppName
params[1] - transformAppName
params[2] - depPath
params[3] - viewName
params[4] - basePath
*/
// applyTranformToModel loads applies the transform to input model
func applyTranformToModel(model, transform *sysl.Module, params ...string) (*sysl.Value, error) {
	modelApp, has := model.Apps[params[0]]
	if !has {
		var apps []string
		for k := range model.Apps {
			apps = append(apps, k)
		}
		sort.Strings(apps)
		return nil, errors.Errorf("app %s does not exist in model, available apps: [%s]", params[0], strings.Join(apps, ", "))
	}
	view := transform.Apps[params[1]].Views[params[3]]
	if view == nil {
		return nil, errors.Errorf("Cannot execute missing view: %s, in app %s", params[3], params[1])
	}
	s := eval.Scope{}
	s.AddApp("app", modelApp)
	s.AddModule("module", model)
	s.AddString("depPath", params[2])
	s["basePath"] = eval.MakeValueString(params[4])
	var result *sysl.Value

	if perTypeTransform(view.Param) {
		result = eval.MakeValueList()
		var tNames []string
		for tName := range modelApp.Types {
			tNames = append(tNames, tName)
		}
		sort.Strings(tNames)
		for _, tName := range tNames {
			t := modelApp.Types[tName]
			s["typeName"] = eval.MakeValueString(tName)
			s["type"] = eval.TypeToValue(t)
			eval.AppendItemToValueList(result.GetList(), eval.EvaluateView(transform, params[1], params[3], s))
		}
	} else {
		result = eval.EvaluateView(transform, params[1], params[3], s)
	}

	return result, nil
}

func perTypeTransform(params []*sysl.Param) bool {
	paramMap := make(map[string]struct{})

	for _, p := range params {
		paramMap[p.Name] = struct{}{}
	}

	if _, has := paramMap["app"]; has {
		if _, has := paramMap["type"]; has {
			return true
		}
	} else {
		panic("Expecting at least an app <: sysl.App")
	}
	return false
}

// Serialize serializes node to string
func Serialize(w io.Writer, delim string, node Node) error {
	for _, n := range node {
		switch x := n.(type) {
		case string:
			if _, err := io.WriteString(w, x+delim); err != nil {
				return err
			}
		case Node:
			if err := Serialize(w, delim, x); err != nil {
				return err
			}
		}
	}
	return nil
}

// GenerateCode transform input sysl model to code in the target language described by
// grammar and a sysl transform
func GenerateCode(
	codegenParams *CmdContextParamCodegen,
	model *sysl.Module, modelAppName string,
	fs afero.Fs, logger *logrus.Logger) ([]*CodeGenOutput, error) {
	var codeOutput []*CodeGenOutput
	depPath := codegenParams.depPath
	basePath := codegenParams.basePath

	logger.Debugf("root-transform: %s\n", codegenParams.rootTransform)
	logger.Debugf("transform: %s\n", codegenParams.transform)
	logger.Debugf("dep-path: %s\n", codegenParams.depPath)
	logger.Debugf("grammar: %s\n", codegenParams.grammar)
	logger.Debugf("start: %s\n", codegenParams.start)
	logger.Debugf("basePath: %s\n", codegenParams.basePath)

	transformFs := syslutil.NewChrootFs(fs, codegenParams.rootTransform)
	tfmParser := parse.NewParser()
	tx, transformAppName, err := parse.LoadAndGetDefaultApp(codegenParams.transform, transformFs, tfmParser)
	if err != nil {
		return nil, err
	}

	g, err := readGrammar(codegenParams.grammar, "gen", codegenParams.start)
	if err != nil {
		return nil, err
	}

	if !codegenParams.disableValidator {
		grammarSysl, err := validate.LoadGrammar(codegenParams.grammar, fs)
		if err != nil {
			msg.NewMsg(msg.WarnValidationSkipped, []string{err.Error()}).LogMsg()
		} else {
			validator := validate.NewValidator(grammarSysl, tx.GetApps()[transformAppName], tfmParser)
			validator.Validate(codegenParams.start, codegenParams.depPath, codegenParams.basePath)
			validator.LogMessages()
		}
	}

	fileNames, err := applyTranformToModel(model, tx, modelAppName, transformAppName,
		depPath, "filename", basePath)
	if err != nil {
		return nil, err
	}
	result, err := applyTranformToModel(model, tx, modelAppName, transformAppName,
		depPath, g.Start, basePath)
	if err != nil {
		return nil, err
	}
	switch {
	case fileNames.GetMap() != nil:
		filename := fileNames.GetMap().Items["filename"].GetS()
		logger.Println(filename)

		if result.GetMap() != nil {
			codeOutput = appendCodeOutput(g, result, logger, codeOutput, filename)
		} else if result.GetList() != nil {
			for _, v := range result.GetList().Value {
				codeOutput = appendCodeOutput(g, v, logger, codeOutput, filename)
			}
		}
	case fileNames.GetList() != nil && result.GetList() != nil:
		fileValues := fileNames.GetList().Value
		for i, v := range result.GetList().Value {
			filename := fileValues[i].GetMap().Items["filename"].GetS()
			codeOutput = appendCodeOutput(g, v, logger, codeOutput, filename)
		}
	default:
		panic("Unexpected combination for filenames and transformation results")
	}

	return codeOutput, nil
}

func appendCodeOutput(g *parser.Grammar, v *sysl.Value,
	logger *logrus.Logger, codeOutput []*CodeGenOutput, filename string) []*CodeGenOutput {
	r := processRule(g, v, g.Start, logger)
	codeOutput = append(codeOutput, &CodeGenOutput{filename, r})
	return codeOutput
}

func outputToFiles(output []*CodeGenOutput, fs afero.Fs) error {
	for _, o := range output {
		f, err := fs.Create(o.filename)
		if err != nil {
			return errors.Wrapf(err, "unable to create %q", o.filename)
		}
		logrus.Infoln("Writing file: " + f.Name())
		if err := Serialize(f, " ", o.output); err != nil {
			return errors.Wrapf(err, "error writing to %q", o.filename)
		}
		if err := f.Close(); err != nil {
			return errors.Wrapf(err, "error closing %q", o.filename)
		}
	}
	return nil
}
