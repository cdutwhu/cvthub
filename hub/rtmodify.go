package main

type IRtModify interface {
	ModifyRet(api, ret string) string
}

var (
	modifiers []IRtModify
)

func AddRtModifier(m IRtModify) {
	modifiers = append(modifiers, m)
}

// func init() {
// 	AddRtModifier(IRtModify)
// }
