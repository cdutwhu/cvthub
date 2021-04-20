package main

import (
	"github.com/digisan/gotk/slice/ts"
)

// --------------------- //

type SifCvt struct {
}

func init() {
	AddRtModifier(&SifCvt{})
}

func (c *SifCvt) ModifyRet(api, ret string) string {

	if ts.NotIn(api, "/sif-xml2json/help", "/sif-json2xml/help") {
		return ret
	}

	ret, err := scanStrLine(ret, func(ln string) (bool, string) {
		if sHasPrefix(sTrimLeft(ln, " \t"), "[POST]") {
			ln = ln[:sLastIndex(ln, "]")+2]
			ln += localIP() + fSf(":%d", PORT) + mApiReDirPOST[api]
		}
		return true, ln
	}, "")
	failOnErr("%v", err)

	return ret
}
