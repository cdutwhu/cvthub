package main

func editRet(retStr, svr string) string {
	help := ""
	switch svr {
	case "x2j-help":
		help = "x2j"
	case "j2x-help":
		help = "j2x"
	}
	retStr, err := readStrLine(retStr, func(ln string) (bool, string) {
		if sHasPrefix(sTrimLeft(ln, " \t"), "[POST]") {
			ln = ln[:sLastIndex(ln, "]")+2]
			ln += localIP() + fSf(":%d", PORT) + mSvrPOSTPath[help]
		}
		return true, ln
	}, "")
	failOnErr("%v", err)
	return retStr
}
