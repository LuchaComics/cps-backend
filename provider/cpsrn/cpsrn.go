package cpsrn

import "fmt"

// Provider provides an interface for abstracting `CPS Registry Number`.
type Provider interface {
	GenerateNumber(currentTotalSubmissionsCount int64) string
}

type cpsrnProvider struct {
	sectionA int64
	sectionB int64
	sectionC int64
}

// NewProvider function is a contructor that returns the default `CPS Registry Number` provider.
func NewProvider() Provider {
	return cpsrnProvider{
		sectionA: 788346,
		sectionB: 26649,
		sectionC: 1001,
	}
}

// Generates the unique `CPS Registry Number` required for tracking submissions.
func (p cpsrnProvider) GenerateNumber(currentTotalSubmissionsCount int64) string {
	newSectionC := p.sectionC + currentTotalSubmissionsCount
	return fmt.Sprintf("%d-%d-%d", p.sectionA, p.sectionB, newSectionC)
}
