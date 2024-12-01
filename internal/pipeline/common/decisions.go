package common

type RuleDecision struct {
	RuleName  string
	Decision  Decision
	AlertInfo *AlertInfo
	Reason    string
}

type FinalDecision struct {
	Decision Decision
	Reason   string
}

type Decision int8

const (
	Ok Decision = iota
	Decline
)
