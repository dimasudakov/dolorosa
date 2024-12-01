package sbp

import (
	"dolorosa/internal/pipeline/contracts"
	"dolorosa/internal/pipeline/dependencies"
	"dolorosa/internal/pipeline/dependencies/nirvana"
	"dolorosa/internal/pipeline/operations/sbp/domain"
	"dolorosa/internal/pipeline/operations/sbp/rules"
)

type b2cRuleRegistry struct {
	rules []contracts.Rule[domain.Operation]
}

type deps = []contracts.Dependency[domain.Operation]

func newSbpRuleRegistry(fields CheckerFields) contracts.RuleRegistry[domain.Operation] {

	nirvanaDep := nirvana.NewNirvanaDependency[domain.Operation](
		[]nirvana.Exception{
			{
				Name: "test_rule",
			},
		},
		fields.ExceptionChecker,
		dependencies.WithOptional(),
	)

	return b2cRuleRegistry{
		rules: []contracts.Rule[domain.Operation]{
			rules.NewTestRule(deps{nirvanaDep}),
		},
	}
}

func (r b2cRuleRegistry) GetRules() []contracts.Rule[domain.Operation] {
	return r.rules
}
