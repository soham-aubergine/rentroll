package rlib

import (
	"fmt"
	"sort"
	"time"
)

// RAStmtEntry describes an entry on a statement
type RAStmtEntry struct {
	T   int                // 1 = assessment, 2 = Receipt
	A   *Assessment        // for type==1, the pointer to the assessment
	R   *ReceiptAllocation // for type ==2, the pointer to the receipt
	RNT *Rentable          // the associated rentable, if known
	Amt float64
	Dt  time.Time
}

// RAStmtEntries is needed to sort the array
type RAStmtEntries []RAStmtEntry

// Len returns the size of the array
func (slice RAStmtEntries) Len() int {
	return len(slice)
}

// Less returns true if element i comes before element j
func (slice RAStmtEntries) Less(i, j int) bool {
	//return slice[i].Name < slice[j].Name;
	return slice[i].Dt.Before(slice[j].Dt)
}

// Swap swaps the two entries
func (slice RAStmtEntries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// RAAcctBal contains details about the account balance for a Rental Agreement
type RAAcctBal struct {
	DtStart    time.Time     // Period Start
	DtStop     time.Time     // Period Stop -- up to but not including
	LmStart    LedgerMarker  // this is the starting point for the calculations
	Gap        RAStmtEntries // these entries cover the gap between the LmStart and Period DtStart
	OpeningBal float64       // balance at the open of period DtStart
	Stmt       RAStmtEntries // these are the actual statement entries
	ClosingBal float64       // balance at close of period
}

// GetRAIDBalance returns the balance of the account for the supplied
// rental agreement on the date requested.
//
//
// Parameters
//     raid  = Rental Agreement ID
//       dt  = date for which balance is needed
//
// Returns
//  balance  = RAID account balance if err == nil
//      err  = any error that occurred or nil if no errors
//
//=============================================================================
func GetRAIDBalance(raid int64, dt *time.Time) (float64, error) {
	bal := float64(0)
	lm := GetRALedgerMarkerOnOrBefore(raid, dt)
	if lm.LMID == 0 {
		err := fmt.Errorf("*** ERROR ***  could not find ledger marker for RAID %d on or before %s", raid, dt.Format(RRDATEFMTSQL))
		return bal, err
	}

	var rs RAStmtEntries
	bal = lm.Balance                               // initialize
	bal += GetRAIDAcctRange(raid, &lm.Dt, dt, &rs) // update with total for this range
	return bal, nil
}

// GetRAIDStatementInfo is written in a way that will work for cash based
// systems or accrual based systems. It looks at all the transactions
// involving the RAID provided and computes a total.  The total is computed
// up-to-and-including d2.
//
// Parameters
//     raid  = Rental Agreement ID
//     d1,d2 = date range for which balance is computed
//
// Returns
//     LmStart
//         is the starting balance for calculations - LmStart.Balance is
//         the opening balance for LmStart.Dt, which is the nearest date on
//         or before d1.
//
//     Gap
//         is the list of Assessments and ReceiptAllocations that occurred
//         from LmStart.Dt up to (but not including) d1
//
//     OpeningBal
//         is the opening balance on d1.  The sum of LmStart.Balance and all
//         entries in Gap.
//
//     Stmt
//         is the list of Assessments and ReceiptAllocations that occurred
//         during the period d1 up to (but not including) d2
//
//=============================================================================
func GetRAIDStatementInfo(raid int64, d1, d2 *time.Time) (RAAcctBal, error) {
	var err error
	var m RAAcctBal

	m.DtStart = *d1
	m.DtStop = *d2

	//----------------------------------------------------------------
	//  First, find the ledger marker for this RentalAgreement...
	//----------------------------------------------------------------
	m.LmStart = GetRALedgerMarkerOnOrBefore(raid, d1)
	if m.LmStart.LMID == 0 {
		err = fmt.Errorf("*** ERROR ***  could not find ledger marker for RAID %d on or before %s", raid, d1.Format(RRDATEFMTSQL))
		return m, err
	}

	m.OpeningBal = m.LmStart.Balance                                  // initialize
	m.OpeningBal += GetRAIDAcctRange(raid, &m.LmStart.Dt, d1, &m.Gap) // update with total for this range
	sort.Sort(m.Gap)

	//----------------------------------------------------------------
	// Now get the actual Statement data and balance...
	//----------------------------------------------------------------
	m.ClosingBal = m.OpeningBal                             // initialize
	m.ClosingBal += GetRAIDAcctRange(raid, d1, d2, &m.Stmt) // update with total for the statement range
	sort.Sort(m.Stmt)

	return m, nil
}

// GetRAIDAcctRange gets the assessment and receipt allocation entries for the
// supplied time range and returns the balance of these entries.
//=============================================================================
func GetRAIDAcctRange(raid int64, d1, d2 *time.Time, p *RAStmtEntries) float64 {
	bal := float64(0)
	//----------------------------------------------------------------
	// Total all assessments in the supplied range that involve RAID.
	//----------------------------------------------------------------
	rows, err := RRdb.Prepstmt.GetAssessmentsByRAIDRange.Query(raid, d1, d2)
	Errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var a Assessment
		ReadAssessments(rows, &a)
		var rnt Rentable
		GetRentableByID(a.RID, &rnt)
		se := RAStmtEntry{
			T:   1,
			A:   &a,
			Amt: a.Amount,
			Dt:  a.Start,
			RNT: &rnt,
		}
		(*p) = append((*p), se)
		bal -= a.Amount
	}

	//----------------------------------------------------------------
	// Total all receipts in the supplied range that involve RAID.
	//----------------------------------------------------------------
	rows, err = RRdb.Prepstmt.GetASMReceiptAllocationsInRAIDDateRange.Query(raid, d1, d2)
	Errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var ra ReceiptAllocation
		ReadReceiptAllocations(rows, &ra)
		a, err := GetAssessment(ra.ASMID)
		Errcheck(err)
		var rnt Rentable
		GetRentableByID(a.RID, &rnt)
		se := RAStmtEntry{
			T:   2,
			R:   &ra,
			A:   &a,
			RNT: &rnt,
			Amt: ra.Amount,
			Dt:  ra.Dt,
		}
		bal += ra.Amount
		(*p) = append((*p), se)
	}
	return bal
}
