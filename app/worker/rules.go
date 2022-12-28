package scannerworker

import (
	"github.com/jiradeto/gh-scanner/app/entities"
)

func GetScannerRules() []entities.ScannerRule {
	return []entities.ScannerRule{
		{
			Type:        "begin_word",
			RuleId:      "secret_1",
			Word:        "public_key",
			Description: "detected secret public_key"},
		{
			Type:        "begin_word",
			RuleId:      "secret_2",
			Word:        "private_key",
			Description: "detected secret private_key",
		},
	}
}
