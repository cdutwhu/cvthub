package main

type IRtModify interface {
	ModifyRet(svr, ret string) string
}

var (
	modifiers []IRtModify
)

func AddRtModifier(m IRtModify) {
	modifiers = append(modifiers, m)
}

func initRtModifier() {
	AddRtModifier(&SifCvt{})
}
