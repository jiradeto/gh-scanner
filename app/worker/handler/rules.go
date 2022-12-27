package workerhandler

import "github.com/jiradeto/gh-scanner/app/domain/entities"

type ScannerRule struct {
	Type        string            `json:"type"`
	RuleId      string            `json:"rule_id"`
	Word        string            `json:"word"`
	Description string            `json:"description"`
	Severity    entities.Severity `json:"severity"`
}

var ScannerRules = []ScannerRule{
	{
		Type:        "begin_word",
		RuleId:      "secret_1",
		Word:        "public_key",
		Description: "detected secret public_key",
	},
	{
		Type:        "begin_word",
		RuleId:      "secret_2",
		Word:        "private_key",
		Description: "detected secret private_key",
	},
}
