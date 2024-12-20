package models

import (
	commontypes "code-compiler/internal/commonTypes"
	"time"
)

type CodeSubmission struct {
	UserId          string                  `json:"userId,omitempty" bson:"userId"`
	Question        string                  `json:"question,omitempty" bson:"question"`
	FailedCase      *commontypes.TestResult `json:"failedCase,omitempty" bson:"failedCase"`
	PassedTestCases int                     `json:"passedTestCases,omitempty" bson:"passedTestCases"`
	TotalTestCases  int                     `json:"totalTestCases,omitempty" bson:"totalTestCases"`
	Err             string                  `json:"err,omitempty" bson:"err"`
	Code            string                  `json:"code,omitempty" bson:"code"`
	Language        string                  `json:"language,omitempty" bson:"language"`
	CreatedAt       time.Time               `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt,omitempty" bson:"updatedAt"`
}
