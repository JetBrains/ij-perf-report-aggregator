package data_query

import (
  "fmt"
  "net/url"
  "testing"
)

func TestEscaping(t *testing.T) {
  s := url.PathEscape("[LDFXEWktellij_sourcesJRvZ'G23Y4Y5Y6Y7Y8Y9P0P1P2P3P4P5P6P7P8P9Q0Q1Q2Q3Q4Q5Q6Q7Q8Q9']BCBO.name'Hz,LDFXEWcommunityJRvZ'G23Y4Y5Y6Y7Y8Y9P0P1P2P3P4P5P6P7P8P9Q0Q1Q2Q3Q4Q5Q6Q7Q8Q9']BCBO.name'Hz]ARGB),('fSCVRqS> subtractMonths{now{}, 3}'DZ('nStRsqlStoUnixTimestamp{V} * EKc1RKc2RKc3Rmachke']~filtersZF1000'),('nSORsubNameSvalue')~tc_kstallGktellij-lkux-hw-blade-0H~vSkdexkg'BO.valueRoS':=Rv!0)]J/kdexkg'BbranchRvSmaster'BmachkeKbuild_L('dbSsharedIndexesRfieldsOmeasuresPA3QA4R'~S!'Vgenerated_timeW('fSprojectRvSXer_KidRtc_KidRYA2Z![kinz~orderSt')\u0001zkZYXWVSRQPOLKJHGFEDCBA")
  println(s)
}

func TestSwap(t *testing.T) {
  s := swap("{\"a\": \"() {} }) {(\"}")
  println(s)
  // ('a': '}} () )} (}')
  expected := "('a! '{} () )} ({')"
  if s != expected {
    t.Fatal("expected " + expected + " but got " + s)
  }
}

func TestDecompress(t *testing.T) {
  s2 := "!(Eamp(genCJDv:ztellij_sources/SHW23X4X5X6X7X8X9O0O1O2O3O4O5O6O7O8O9P0P1P2P3P4P5P6P7P8P9)BGBK.FZQEamp(genCJDv:community/SHW23X4X5X6X7X8X9O0O1O2O3O4O5O6O7O8O9P0P1P2P3P4P5P6P7P8P9)BGBK.FZ))A,VWBQ(f:CY)+*+1000'Q(n:K,subName:valueQDk2k3,machzeQfilters:!((f:project,E(db:sharedIndexes,fields:!((n:t,sql~toUnixTimestFnameRSK.value,o~!!%3D'R0)QGgenY,q~%3E+subtractMonths(now(Q+3)'HbranchRmasterBmachzeR!(VJtc_zstaller_Lid,tc_Lidk1KmeasuresLbuild_OA3PA4Q),R,v:SzdexzgBVztellij-lzux-hwW-blade-0XA2Yerated_timeZorder:tk,Lczin~:'%01~zkZYXWVSRQPOLKJHGFEDCBA"
  var err error
  s2, err = url.PathUnescape(s2)
  if err != nil {
    t.Error(err)
  }
  s := uncrush(s2)
  fmt.Print(s)
}

func TestDecompress2(t *testing.T) {
  s2 := "!((fieldsG(n:t,sqlDtoUnixTimestamp(A)*1000'4(n:C4tc_installerK,tcKBc1Bc2Bc3,E4db:ij,filtersG9product8IUHproject8'simple%20for%20IJ'HE8!(.550F551F772F773)HA,qD%3EsubtractMonths(now(41)'HC,oD!!%3D'80)4order:t).intellij-macos-hw-unit-14)Bbuild_8,v:9(f:Agenerated_timeB,6CappInit_dD:'EmachineF,.G:!(H49K_6id%01KHGFEDCBA9864.,(fieldsG(n:t,sqlDtoUnixTimestamp(A)*1000'4(n64tc_installerK,tcKC1C2C3,E4db:ij,filtersGHduct9IU4Hject9'simple%20for%20IJ'4BE9!(.550F551F772F773)4BA,qD%3EsubtractMonths(now(41)'4(f6,oD!!%3D'90)4order:t).intellij-macos-hw-unit-14),6:appInitPreparation_d8build_9,v:Agenerated_timeB(f:C,8cD:'EmachineF,.G:!(HBproK_8id%01KHGFEDCBA9864.,(fieldsH(n:t,sqlDtoUnixTimestamp(A)*1000'4(n:B4tc_installGL,tcLCc1Cc2Cc3,E4db:ij,filtGsH9product8IUKproject8'simple%20for%20IJ'KE8!(.550F551F772F773)KA,qD%3EsubtractMonths(now(41)'KB,oD!!%3D'80)4ordG:t).intellij-macos-hw-unit-14)Cbuild_8,v:9(f:AgenGated_timeBappStartG_dC,6D:'EmachineF,.GerH:!(K49L_6id%01LKHGFEDCBA9864.,(fieldsG(n:t,sqlDtoUnixTimestamp(A)*1000'4(n:B4tc_installerK,tcKCc1Cc2Cc3,E4db:ij,filtersG9product8IUHproject8'simple%20for%20IJ'HE8!(.550F551F772F773)HA,qD%3EsubtractMonths(now(41)'HB,oD!!%3D'80)4order:t).intellij-macos-hw-unit-14)Cbuild_8,v:9(f:Agenerated_timeBbootstrap_dC,6D:'EmachineF,.G:!(H49K_6id%01KHGFEDCBA9864.,(fieldsG(n:t,sqlDtoUnixTimestamp(A)*1000'4(n:B4tc_HstallerL,tcLCc1Cc2Cc3,E4db:ij,filtersG9product8IUKproject8'simple%20for%20IJ'KE8!(.550F551F772F773)KA,qD%3EsubtractMonths(now(41)'KB,oD!!%3D'80)4order:t).Htellij-macos-hw-unit-14)Cbuild_8,v:9(f:Agenerated_timeBeuaShowHg_dC,6D:'EmachHeF,.G:!(HinK49L_6id%01LKHGFEDCBA9864.,(fieldsK(n:t,sqlEtoUnixTimestamp(A)*1000'6(n46tc_HstallerO,tcOCc1Cc2Cc3,F6db:ij,filtersKNduct9IU6Nject9'simple%20for%20IJ'6BF9!(.550G551G772G773)6BA,qE%3EsubtractMonths(now(61)'6(f4,oE!!%3D'90)6order:t).Htellij-macos-hw-unit-14:plugHDescriptorLoadHg_d6)Cbuild_9,v:Agenerated_timeB(f:C,8E:'FmachHeG,.HinK:!(NBproO_8id%01ONKHGFECBA9864.)"
  var err error
  s2, err = url.PathUnescape(s2)
  if err != nil {
    t.Error(err)
  }
  s := uncrush(s2)
  fmt.Print(s)
}
