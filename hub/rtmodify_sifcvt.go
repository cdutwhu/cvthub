package main

// --------------------- //

type SifCvt struct {
}

func (c *SifCvt) ModifyRet(svr, ret string) string {
	help := ""
	switch svr {
	case "x2j-help":
		help = "x2j"
	case "j2x-help":
		help = "j2x"
	}
	ret, err := scanStrLine(ret, func(ln string) (bool, string) {
		if sHasPrefix(sTrimLeft(ln, " \t"), "[POST]") {
			ln = ln[:sLastIndex(ln, "]")+2]
			ln += localIP() + fSf(":%d", PORT) + mSvrPOSTPath[help]
		}
		return true, ln
	}, "")
	failOnErr("%v", err)
	return ret
}

func init() {
	AddRtModifier(&SifCvt{})
}
