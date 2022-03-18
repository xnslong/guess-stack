package utils

type FixNeeded interface {
	NeedFix() bool
	SetNeedFix(need bool)
}
