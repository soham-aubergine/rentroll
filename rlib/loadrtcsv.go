package rlib

import (
	"fmt"
	"strconv"
	"strings"
)

// GetBusinessBID returns the BID for the Business with the supplied designation
func GetBusinessBID(des string) int64 {
	//-------------------------------------------------------------------
	// Make sure the Business exists...
	//-------------------------------------------------------------------
	b, err := GetBusinessByDesignation(des)
	if err != nil || b.BID == 0 {
		Ulog("GetBusinessBID: Business with designation %s does not exist or could not be loaded\n", des)
		return 0
	}
	return b.BID
}

// CreateRentableType reads an Rentable type string array and creates a database record for the Rentable type
//
//                                                                     Repeat as many 3-tuples as needed
//                                                                         /----------^-------------\
//  [0]        [1]      [2]   			[3]            [4]       [5]            6          7       8
// Designation,Style,	Name, 			RentCycle,  Proration,ManageToBudget,MarketRate,DtStart,DtStop
// REH,        "GM",	"Geezer Miser", 6,		       4,        1,             1100.00,   1/1/2015, 1/1/2017
// REH,        "FS",	"Flat Studio",  6,		       4,        1,             1500.00,   1/1/2015, 1/1/2017
// REH,        "SBL",	"SB Loft",     	6,		       4,        1,             1750.00,   1/1/2015, 1/1/2017
// REH,        "KDS",	"KD Suite",    	6,		       4,        1,             2000.00,   1/1/2015, 1/1/2017
// REH,        "VEH",	Vehicle,       	3,		       0,        1,             10.0,   1/1/2015, 1/1/2017
// REH,        "CPT",	Carport,       	6,		       4,        1,             35.0,   1/1/2015, 1/1/2017
func CreateRentableType(sa []string, lineno int) {
	funcname := "CreateRentableType"
	if 8 > len(sa) {
		Ulog("CreateRentableType: csv file line \"%s\" does not have 7 elements. Ignored.\n", sa)
		return
	}
	des := strings.TrimSpace(sa[0])
	if des == "Designation" {
		return // this is just the column heading
	}
	// fmt.Printf("line %d, sa = %#v\n", lineno, sa)
	required := 8
	if len(sa) < required {
		fmt.Printf("%s: line %d - found %d values, there must be at least %d\n", funcname, lineno, len(sa), required)
		return
	}

	//-------------------------------------------------------------------
	// Check to see if this Rentable type is already in the database
	//-------------------------------------------------------------------
	var a RentableType
	bid := GetBusinessBID(des)
	if bid == 0 {
		return
	}

	a.BID = bid
	a.Style = strings.TrimSpace(sa[1])
	if len(a.Style) > 0 {
		rt, err := GetRentableTypeByStyle(a.Style, bid)
		if nil != err && !IsSQLNoResultsError(err) {
			Ulog("%s: line %d - err = %v\n", funcname, lineno, err)
			return
		}
		if rt.RTID > 0 {
			Ulog("%s: line %d - RentableType named %s already exists\n", funcname, lineno, a.Style)
			return
		}
	}

	a.Name = strings.TrimSpace(sa[2])

	//-------------------------------------------------------------------
	// Load the values based on csv input
	//-------------------------------------------------------------------
	n, err := strconv.Atoi(strings.TrimSpace(sa[3])) // frequency
	if err != nil || !IsValidAccrual(int64(n)) {
		Ulog("%s: line %d - Invalid rental frequency: %s\n", funcname, lineno, sa[3])
		return
	}
	a.RentCycle = int64(n)

	n, err = strconv.Atoi(strings.TrimSpace(sa[4])) // Proration
	if err != nil || !IsValidAccrual(int64(n)) {
		Ulog("%s: line %d - Invalid rental proration frequency: %s\n", funcname, lineno, sa[4])
		return
	}
	a.Proration = int64(n)
	if a.Proration > a.RentCycle {
		Ulog("%s: line %d - Proration frequency (%d) must be greater than rental frequency (%d)\n", funcname, lineno, a.Proration, a.RentCycle)
		return
	}

	n64, err := yesnoToInt(strings.TrimSpace(sa[5])) // manage to budget
	if err != nil {
		Ulog("%s: line %d - Invalid manage to budget flag: %s\n", funcname, lineno, sa[5])
		return
	}
	a.ManageToBudget = int64(n64)

	rtid, err := InsertRentableType(&a)

	// Rentable Market Rates are provided in 3-tuples starting at index 7 - Amount,startdata,enddate
	if rtid > 0 {
		for i := 6; i < len(sa); i += 3 {
			var x float64
			var err error
			var m RentableMarketRate
			m.RTID = rtid
			if x, err = strconv.ParseFloat(strings.TrimSpace(sa[i]), 64); err != nil {
				Ulog("%s: line %d - Invalid floating point number: %s\n", funcname, lineno, sa[7])
				return
			}
			m.MarketRate = x
			DtStart, err := StringToDate(sa[i+1])
			if err != nil {
				fmt.Printf("%s: line %d - invalid start date:  %s\n", funcname, lineno, sa[i+1])
				return
			}
			m.DtStart = DtStart
			DtStop, err := StringToDate(sa[i+2])
			if err != nil {
				fmt.Printf("%s: line %d - invalid stop date:  %s\n", funcname, lineno, sa[i+2])
				return
			}
			m.DtStop = DtStop
			if m.DtStart.After(m.DtStop) {
				fmt.Printf("%s: line %d - Stop date (%s) must be after Start date (%s)\n", funcname, lineno, m.DtStop, m.DtStart)
				return
			}
			InsertRentableMarketRates(&m)
		}
	}
}

// LoadRentableTypesCSV loads a csv file with Rentable types and processes each one
func LoadRentableTypesCSV(fname string) {
	t := LoadCSV(fname)
	for i := 0; i < len(t); i++ {
		CreateRentableType(t[i], i+1)
	}
}
