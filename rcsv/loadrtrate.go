package rcsv

import (
	"fmt"
	"rentroll/rlib"
	"strings"
)

//  CSV file format:
//                             RT name or style  string with or without %
// 0    1             2        3                 4
// BUD, RPName,       RPRID    RentableType,     Amount
// REX, FAA-P, RPR0001, GM,               85%
// REX, FAA-P,   1,     Flat Studio,      1400
// REX, FAA-P,   1,     SBL,    			 1500
// REX, FAA-P,   1,     KDS,    			 75%
// REX, FAA-T,      1,     GM,               90%
// REX, FAA-T,      1,     Flat Studio,      90%
// REX, FAA-T,      1,     SBL,    			 1500
// REX, FAA-T,      1,     KDS,    			 80%

// CreateRatePlanRefRTRate reads a rental specialty type string array and creates a database record for the rental specialty type.
func CreateRatePlanRefRTRate(sa []string, lineno int) (string, int) {
	funcname := "CreateRatePlanRefRTRate"
	var b rlib.Business

	const (
		BUD          = 0
		RPName       = iota
		RPRID        = iota
		RentableType = iota
		Amount       = iota
	)

	// csvCols is an array that defines all the columns that should be in this csv file
	var csvCols = []CSVColumn{
		{"BUD", BUD},
		{"RPName", RPName},
		{"RPRID", RPRID},
		{"RentableType", RentableType},
		{"Amount", Amount},
	}

	rs, x := ValidateCSVColumns(csvCols, sa, funcname, lineno)
	if x > 0 {
		return rs, 1
	}
	if lineno == 1 {
		return rs, 0
	}

	des := strings.ToLower(strings.TrimSpace(sa[BUD]))

	//-------------------------------------------------------------------
	// BUD
	//-------------------------------------------------------------------
	if len(des) > 0 {
		b = rlib.GetBusinessByDesignation(des)
		if len(b.Designation) == 0 {
			rs += fmt.Sprintf("%s: line %d, rlib.Business with designation %s does not exist\n", funcname, lineno, sa[BUD])
			return rs, CsvErrorSensitivity
		}
	}

	//-------------------------------------------------------------------
	// RatePlan Name
	//-------------------------------------------------------------------
	var rp rlib.RatePlan
	rpname := strings.ToLower(strings.TrimSpace(sa[RPName]))
	if len(rpname) > 0 {
		rlib.GetRatePlanByName(b.BID, rpname, &rp)
		if rp.RPID < 1 {
			rs += fmt.Sprintf("%s: line %d - RatePlan named %s not found\n", funcname, lineno, rpname)
			return rs, CsvErrorSensitivity
		}
	}

	var a rlib.RatePlanRefRTRate
	var err error
	var errmsg string

	//-------------------------------------------------------------------
	// RPRef
	//-------------------------------------------------------------------
	a.RPRID = CSVLoaderGetRPRID(strings.TrimSpace(sa[RPRID]))
	if 0 == a.RPRID {
		rs += fmt.Sprintf("%s: line %d - Bad value for RatePlanRef ID: %s\n", funcname, lineno, sa[RPRID])
		return rs, CsvErrorSensitivity
	}

	//-------------------------------------------------------------------
	// RT Style
	// identifies the RentableType
	//-------------------------------------------------------------------
	name := strings.TrimSpace(sa[RentableType])
	rt, err := rlib.GetRentableTypeByStyle(name, b.BID)
	if err != nil {
		rs += fmt.Sprintf("%s: line %d - could not load RentableType with Style = %s,  err:  %s\n", funcname, lineno, sa[RentableType], err.Error())
		return rs, CsvErrorSensitivity
	}
	a.RTID = rt.RTID

	//-------------------------------------------------------------------
	// Amount
	// Entered as a string. If the string contains a % then the amount
	// is a percentage and we set the % flag. Otherwise, it is an absolute amount
	//-------------------------------------------------------------------
	amt := strings.TrimSpace(sa[Amount])
	a.Val, errmsg = rlib.FloatFromString(amt, "bad amount")
	if len(errmsg) > 0 {
		return rs, CsvErrorSensitivity
	}
	if strings.Contains(amt, "%") {
		a.FLAGS |= rlib.FlRTRpct
	}

	//-------------------------------------------------------------------
	// Insert the record
	//-------------------------------------------------------------------
	err = rlib.InsertRatePlanRefRTRate(&a)
	if nil != err {
		rs += fmt.Sprintf("%s: lineno %d  - error inserting RatePlanRefRTRate = %v\n", funcname, lineno, err)
		return rs, CsvErrorSensitivity
	}
	return rs, 0
}

// LoadRatePlanRefRTRatesCSV loads a csv file with RatePlan rates for specific rentable types
func LoadRatePlanRefRTRatesCSV(fname string) string {
	rs := ""
	t := rlib.LoadCSV(fname)
	for i := 0; i < len(t); i++ {
		s, err := CreateRatePlanRefRTRate(t[i], i+1)
		rs += s
		if err > 0 {
			break
		}
	}
	return rs
}
