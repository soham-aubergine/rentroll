package main

import (
	"fmt"
	"rentroll/rlib"
	"strings"
	"time"
)

// TFMTSPACE et al control the formatting of the Journal report
const (
	TFMTSPACE   = 2  // space between cols
	TFMTINDENT  = 3  // left indent
	TFMTDESCR   = 45 // description width
	TFMTDATE    = 8  // date width
	TFMTRA      = 10 // rental agreement
	TFMTJID     = 9  // Journal id
	TFMTRN      = 15 // Rentable name
	TFMTAMOUNT  = 12 // balance width
	TFMTDECIMAL = 2  // number of decimal places
	TLINELEN    = 6*TFMTSPACE + TFMTDESCR + TFMTDATE + TFMTJID + TFMTRA + TFMTRN + TFMTAMOUNT + TFMTAMOUNT
)

var tfmt struct {
	Indent             string // left indent
	Descr              string // Description
	DescrLJ            string
	Dt                 string // date
	JID                string // transaction id
	RentalAgr          string // rental agreement
	RentableName       string // Rentable name
	Amount             string // amount
	Balance            string // balance
	Sp                 string
	Hdr                string
	AmountHdrStr       string
	JIDHdrStr          string
	DatedLedgerEntryRJ string
	DatedLedgerEntryLJ string
	LedgerHeading      string
	DescrAndBal        string
}

func initTFmt() {
	s := fmt.Sprintf("%%%ds", TFMTINDENT)
	tfmt.Indent = fmt.Sprintf(s, "")
	s = fmt.Sprintf("%%%ds", TFMTSPACE)
	tfmt.Sp = fmt.Sprintf(s, " ")
	tfmt.Descr = fmt.Sprintf("%%%ds", TFMTDESCR)                   // Description
	tfmt.DescrLJ = fmt.Sprintf("%%-%ds", TFMTDESCR)                // Description LJ
	tfmt.Dt = fmt.Sprintf("%%%ds", TFMTDATE)                       // date
	tfmt.JID = fmt.Sprintf("%%%dd", TFMTJID)                       // Journal id
	tfmt.JIDHdrStr = fmt.Sprintf("%%%ds", TFMTJID)                 // amount for header
	tfmt.RentalAgr = fmt.Sprintf("%%%ds", TFMTRA)                  // rental agreement
	tfmt.RentableName = fmt.Sprintf("%%%ds", TFMTRN)               // Rentable name
	tfmt.Amount = fmt.Sprintf("%%%d.%df", TFMTAMOUNT, TFMTDECIMAL) // digits
	tfmt.AmountHdrStr = fmt.Sprintf("%%%ds", TFMTAMOUNT)           // amount

	// Descr, Date, JID, Rental Agreement,  Rentable name,  Amount, Balance
	tfmt.DatedLedgerEntryRJ = tfmt.Indent + tfmt.Descr + tfmt.Sp + tfmt.Dt + tfmt.Sp + tfmt.JID +
		tfmt.Sp + tfmt.RentalAgr + tfmt.Sp + tfmt.RentableName + tfmt.Sp + tfmt.Amount +
		tfmt.Sp + tfmt.Amount + "\n"
	tfmt.DatedLedgerEntryLJ = tfmt.Indent + tfmt.DescrLJ + tfmt.Sp + tfmt.Dt + tfmt.Sp + tfmt.JID +
		tfmt.Sp + tfmt.RentalAgr + tfmt.Sp + tfmt.RentableName + tfmt.Sp + tfmt.Amount +
		tfmt.Sp + tfmt.Amount + "\n"
	tfmt.LedgerHeading = tfmt.Indent + tfmt.DescrLJ + "\n"
	tfmt.Hdr = tfmt.Indent + tfmt.DescrLJ + tfmt.Sp + tfmt.Dt + tfmt.Sp + tfmt.JIDHdrStr + tfmt.Sp +
		tfmt.RentalAgr + tfmt.Sp + tfmt.RentableName + tfmt.Sp + tfmt.AmountHdrStr + tfmt.Sp +
		tfmt.AmountHdrStr + "\n"
	tfmt.DescrAndBal = tfmt.Indent + tfmt.DescrLJ + tfmt.Sp + tfmt.Dt + tfmt.Sp +
		fmt.Sprintf(fmt.Sprintf("%%%ds", TFMTJID+TFMTSPACE+TFMTRA+TFMTSPACE+TFMTRN+TFMTSPACE+TFMTAMOUNT+TFMTSPACE), " ") +
		tfmt.Amount + "\n"
}

func printTLineOf(s string) {
	fmt.Println(strings.Repeat(" ", TFMTINDENT) + strings.Repeat(s, TLINELEN/len(s)))
}
func printTReportDoubleLine() {
	printTLineOf("=")
}
func printTReportLine() {
	printTLineOf("-")
}
func printTReportThinLine() {
	printTLineOf(" -")
}
func printTSubtitle(label string) {
	fmt.Printf(tfmt.LedgerHeading, label)
}
func printDatedLedgerEntryRJ(label string, d time.Time, jid int64, ra string, rn string, a, b float64) {
	fmt.Printf(tfmt.DatedLedgerEntryRJ, label, d.Format(rlib.RRDATEFMT), jid, ra, rn, a, b)
}
func printDatedLedgerEntryLJ(label string, d time.Time, jid int64, ra string, rn string, a, b float64) {
	fmt.Printf(tfmt.DatedLedgerEntryLJ, label, d.Format(rlib.RRDATEFMT), jid, ra, rn, a, b)
}
func printLedgerHeaderText(l *rlib.Ledger) {
	printTSubtitle(l.Name)
}
func printLedgerDescrAndBal(s string, d time.Time, x float64) {
	fmt.Printf(tfmt.DescrAndBal, s, d.Format(rlib.RRDATEFMT), x)
}

//
func printLedgerHeader(xbiz *rlib.XBusiness, l *rlib.Ledger, d1, d2 *time.Time) {
	printTReportDoubleLine()
	fmt.Printf("   Business:  %-13s\n", xbiz.P.Name)
	printTSubtitle(l.GLNumber + " - " + l.Name)
	fmt.Printf("   %s - %s\n", d1.Format(rlib.RRDATEFMT), d2.AddDate(0, 0, -1).Format(rlib.RRDATEFMT))
	printTReportLine()
	fmt.Printf(tfmt.Hdr, "Description", "Date", "JournalID", "RntAgr", "Rentable", "Amount", "Balance")
	printTReportLine()
}

// returns the payment/accessment reason, Rentable name
func getLedgerEntryDescription(l *rlib.LedgerEntry) (string, string, string) {
	j, _ := rlib.GetJournal(l.JID)
	sra := fmt.Sprintf("%9d", j.RAID)
	switch j.Type {
	case rlib.JNLTYPEUNAS:
		r := rlib.GetRentable(j.ID) // j.ID is set to RID when the type is unassociated
		return "Unassociated", r.Name, sra
	case rlib.JNLTYPERCPT:
		ja, _ := rlib.GetJournalAllocation(l.JAID)
		a, _ := rlib.GetAssessment(ja.ASMID)
		r := rlib.GetRentable(a.RID)
		return "Payment - " + App.AsmtTypes[a.ASMTID].Name, r.Name, sra
	case rlib.JNLTYPEASMT:
		a, _ := rlib.GetAssessment(j.ID)
		r := rlib.GetRentable(a.RID)
		return "Assessment - " + App.AsmtTypes[a.ASMTID].Name, r.Name, sra

	default:
		fmt.Printf("getLedgerEntryDescription: unrecognized type: %d\n", j.Type)
	}
	return "x", "x", "x"
}

func reportTextProcessLedgerMarker(xbiz *rlib.XBusiness, lm *rlib.LedgerMarker, d1, d2 *time.Time) {
	l, err := rlib.GetLedger(lm.LID)
	rlib.Errcheck(err)
	bal := lm.Balance
	printLedgerHeader(xbiz, &l, d1, d2)
	printLedgerDescrAndBal("Opening Balance", *d1, lm.Balance)
	rows, err := rlib.RRdb.Prepstmt.GetLedgerEntriesInRangeByGLNo.Query(l.BID, l.GLNumber, d1, d2)
	rlib.Errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var l rlib.LedgerEntry
		rlib.Errcheck(rows.Scan(&l.LEID, &l.BID, &l.JID, &l.JAID, &l.GLNumber, &l.Dt, &l.Amount, &l.Comment, &l.LastModTime, &l.LastModBy))
		bal += l.Amount
		descr, rn, sra := getLedgerEntryDescription(&l)
		printDatedLedgerEntryRJ(descr, l.Dt, l.JID, sra, rn, l.Amount, bal)
	}
	rlib.Errcheck(rows.Err())
	printTReportLine()
	printLedgerDescrAndBal("Closing Balance", d2.AddDate(0, 0, -1), bal)
	fmt.Printf("\n\n")
}

// LedgerReportText generates a textual Journal report for the supplied Business and time range
func LedgerReportText(xbiz *rlib.XBusiness, d1, d2 *time.Time) {
	t := rlib.GetLedgerList(xbiz.P.BID) // this list contains the list of all Ledger account numbers
	for i := 0; i < len(t); i++ {
		dd2 := d1.AddDate(0, 0, -1)
		dd1 := time.Date(dd2.Year(), dd2.Month(), 1, 0, 0, 0, 0, dd2.Location())
		lm, err := rlib.GetLedgerMarkerByGLNoDateRange(xbiz.P.BID, t[i].GLNumber, &dd1, &dd2)
		if lm.LMID < 1 || err != nil {
			fmt.Printf("LedgerReportText: GLNumber %s -- no Ledger Marker for: %s - %s\n",
				t[i].GLNumber, dd1.Format(rlib.RRDATEFMT), dd2.Format(rlib.RRDATEFMT))
		} else {
			reportTextProcessLedgerMarker(xbiz, &lm, d1, d2)
		}
	}
}
