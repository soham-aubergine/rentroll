#!/bin/bash
TESTNAME="JM1"
TESTSUMMARY="Setup and run JM1 company and the Rexford Properties"

RRDATERANGE="-j 2016-01-01 -k 2016-02-01"

source ../share/base.sh

#========================================================================================
# INITIALIZE THE BUSINESS
#   This section has the 1-time tasks to set up the business and get the accounts to
#   their correct starting values.
#========================================================================================
docsvtest "a" "-b business.csv -L 3" "Business"
docsvtest "b" "-c coa.csv -L 10,${BUD}" "ChartOfAccounts"
docsvtest "c" "-m depmeth.csv -L 23,${BUD}" "DepositMethods"
docsvtest "d" "-d depository.csv -L 18,${BUD}" "Depositories"
docsvtest "e" "-R rentabletypes.csv -L 5,${BUD}" "RentableTypes"
docsvtest "f" "-l strlists.csv -L 25,${BUD}" "StringLists"
docsvtest "g" "-p people.csv  -L 7" "People"
docsvtest "h" "-r rentable.csv -L 6,${BUD}" "Rentables"
docsvtest "i" "-u custom.csv -L 14" "CustomAttributes"
docsvtest "j" "-U assigncustom.csv -L 15" "AssignCustomAttributes"
docsvtest "k" "-T ratemplates.csv  -L 8" "RentalAgreementTemplates"
docsvtest "l" "-C ra.csv -L 9,${BUD}" "RentalAgreements"
docsvtest "m" "-P pmt.csv -L 12,${BUD}" "PaymentTypes"

# get the deposits on the books
docsvtest "n" "-A asm2015Dec.csv -G ${BUD} -g 12/1/15,1/1/16 -L 11,${BUD}" "Assessments-2015-DEC"
docsvtest "o" "-e rcpt2015Dec.csv -G ${BUD} -g 12/1/15,1/1/16 -L 13,${BUD}" "Receipts-2015-DEC"

# validate GSR
dorrtest "p" "-j 2015-12-01 -k 2016-01-01 -b ${BUD} -r 11" "GSR"

#  INITIALIZE database with deposit information and verify Accounts
dorrtest "q" "-j 2015-12-01 -k 2016-01-01 -x -b ${BUD}" "Process-2015-DEC"
dorrtest "r" "-j 2015-12-01 -k 2016-01-01 -b ${BUD} -r 1" "Journal"
dorrtest "s" "-j 2015-12-01 -k 2016-01-01 -b ${BUD} -r 2" "Ledgers"
dorrtest "t" "-r 12,1,RA001,2016-01-01 -b ${BUD}" "AccountBalance-GeneralAccountsReceivable-RA01"
dorrtest "u" "-r 12,7,RA001,2016-01-01 -b ${BUD}" "AccountBalance-SecurityDeposits-RA01"
dorrtest "v" "-r 12,1,RA002,2016-01-01 -b ${BUD}" "AccountBalance-GeneralAccountsReceivable-RA02"
dorrtest "x" "-r 12,7,RA002,2016-01-01 -b ${BUD}" "AccountBalance-SecurityDeposits-RA02"
dorrtest "y" "-r 12,1,RA003,2016-01-01 -b ${BUD}" "AccountBalance-GeneralAccountsReceivable-RA-03"
dorrtest "z" "-r 12,7,RA003,2016-01-01 -b ${BUD}" "AccountBalance-SecurityDeposits-RA-03"

#========================================================================================
# JANUARY 2016
#    Normal month
#========================================================================================
RRDATERANGE="-j 2016-01-01 -k 2016-02-01"
CSVLOADRANGE="-G ${BUD} -g 1/1/16,2/1/16"
# 1.  Generate recurring assessment instances  -  Note: will be done by server automatically  (18 processes journal only)
dorrtest "a1" "${RRDATERANGE} -x -b ${BUD} -r 18" "Process-2016-JAN"

# 2.  Load new assessments for this period.  For this test, we start the rent assessments now.
docsvtest "b1" "-A asm2016Jan.csv ${CSVLOADRANGE} -L 11,${BUD}" "Assessments-2016-JAN"

# 3.  Create Invoices for each tenant
docsvtest "c1" "-i invoice-2016Jan-Read.csv ${CSVLOADRANGE} -L 20,REX" "Invoice-2016Jan-Read"
docsvtest "d1" "-i invoice-2016Jan-Costea.csv ${CSVLOADRANGE} -L 20,REX" "invoice-2016Jan-Costea"
docsvtest "e1" "-i invoice-2016Jan-Haroutunian.csv ${CSVLOADRANGE} -L 20,REX" "invoice-2016Jan-Haroutunian"
dorrtest "f1" "${RRDATERANGE} -b ${BUD} -r 9,IN001" "InvoiceReport-2016Jan-Read"
dorrtest "g1" "${RRDATERANGE} -b ${BUD} -r 9,2" "InvoiceReport-2016Jan-Costea"
dorrtest "h1" "${RRDATERANGE} -b ${BUD} -r 9,3" "InvoiceReport-2016Jan-Haroutunian"

# 4. Enter any receipts (and assessments if any) since Jan1 - end of the month
docsvtest "i1" "-e rcpt2016Jan.csv ${CSVLOADRANGE} -L 13,${BUD}" "Receipts-2016-JAN"

# 5. Create deposits for all receipts
docsvtest "j1" "-y deposit-2016Jan.csv ${CSVLOADRANGE} -L 19,${BUD}" "Deposits-2016-JAN"

# 6. Process anything that was just added
dorrtest "k3" "${RRDATERANGE} -b ${BUD}" "Finish-2016-JAN"

# 7. Generate final reports for the month
dorrtest "l1" "${RRDATERANGE} -b ${BUD} -r 1" "Journal"
dorrtest "m1" "${RRDATERANGE} -b ${BUD} -r 2" "Ledgers"
dorrtest "n1" "${RRDATERANGE} -b ${BUD} -r 10" "LedgerActivity"
dorrtest "o1" "${RRDATERANGE} -b ${BUD} -r 17" "LedgerBalance"
dorrtest "p1" "${RRDATERANGE} -b ${BUD} -r 8" "Statements"
dorrtest "q1" "${RRDATERANGE} -b ${BUD} -r 4" "RentRoll"


#========================================================================================
# FEBRUARY 2016
#    Haroutunian moves out on Feb 8
#========================================================================================
RRDATERANGE="-j 2016-02-01 -k 2016-03-01"
CSVLOADRANGE="-G ${BUD} -g 2/1/16,3/1/16"
dorrtest  "a2" "${RRDATERANGE} -x -b ${BUD} -r 18" "Process-2016-FEB"
docsvtest "b2" "-A asm2016Feb.csv ${CSVLOADRANGE} -L 11,${BUD}" "Assessments-2016-FEB"
docsvtest "c2" "-i invoice-2016Feb-Read.csv ${CSVLOADRANGE} -L 20,REX" "Invoice-2016Feb-Read"
docsvtest "d2" "-i invoice-2016Feb-Costea.csv ${CSVLOADRANGE} -L 20,REX" "invoice-2016Feb-Costea"
docsvtest "e2" "-i invoice-2016Feb-Haroutunian.csv ${CSVLOADRANGE} -L 20,REX" "invoice-2016Feb-Haroutunian"
dorrtest  "f2" "${RRDATERANGE} -b ${BUD} -r 9,IN001" "InvoiceReport-2016Feb-Read"
dorrtest  "g2" "${RRDATERANGE} -b ${BUD} -r 9,2" "InvoiceReport-2016Feb-Costea"
dorrtest  "h2" "${RRDATERANGE} -b ${BUD} -r 9,3" "InvoiceReport-2016Feb-Haroutunian"
docsvtest "i2" "-e rcpt2016Feb.csv ${CSVLOADRANGE} -L 13,${BUD}" "Receipts-2016-FEB"
docsvtest "j2" "-y deposit-2016Feb.csv ${CSVLOADRANGE} -L 19,${BUD}" "Deposits-2016-FEB"
dorrtest  "k2" "${RRDATERANGE} -b ${BUD}" "Finish-2016-FEB"
dorrtest  "l2" "${RRDATERANGE} -b ${BUD} -r 1" "Journal"
dorrtest  "m2" "${RRDATERANGE} -b ${BUD} -r 2" "Ledgers"
dorrtest  "n2" "${RRDATERANGE} -b ${BUD} -r 10" "LedgerActivity"
dorrtest  "o2" "${RRDATERANGE} -b ${BUD} -r 17" "LedgerBalance"
dorrtest  "p2" "${RRDATERANGE} -b ${BUD} -r 8" "Statements"
dorrtest  "q2" "${RRDATERANGE} -b ${BUD} -r 4" "RentRoll"

#========================================================================================
# MARCH 2016
#    GSR and Contract rent change to 3750 for 309 Rexford
#    Haroutunian receives 865.29 Deposit return, forfeits the rest
#========================================================================================

##-----------------------------------------------------
##  1. Update end date on RentalAgreement 1 to 3/1/18
##  2. Update ContractRent to $3750/month
##  3. Update MarketRate to $3750/month
##-----------------------------------------------------
cat >xxyyzz <<EOF
use rentroll
update RentalAgreement SET AgreementStop="2018-03-01",PossessionStop="2018-03-01",RentStop="2018-03-01" WHERE RAID=1;
INSERT INTO RentalAgreementRentables (RAID,RID,CLID,ContractRent,DtStart,DtStop) VALUES(1,1,0,3750,"2016-03-01 00:00:00","2018-03-01 00:00:00");
INSERT INTO RentableMarketRate (RTID,MarketRate,DtStart,DtStop) VALUES(1,3750,"2016-03-01 00:00:00","2018-03-01 00:00:00");
EOF
${MYSQL} --no-defaults <xxyyzz
rm -f xxyyzz

RRDATERANGE="-j 2016-03-01 -k 2016-04-01"
CSVLOADRANGE="-G ${BUD} -g 3/1/16,4/1/16"
docsvtest "b3" "-A asm2016Mar.csv ${CSVLOADRANGE} -L 11,${BUD}" "Assessments-2016-Mar"  
dorrtest  "a3" "${RRDATERANGE} -x -b ${BUD} -r 18" "Process-2016-Mar"
docsvtest "i3" "-e rcpt2016Mar.csv ${CSVLOADRANGE} -L 13,${BUD}" "Receipts-2016-Mar"
docsvtest "j3" "-y deposit-2016Mar.csv ${CSVLOADRANGE} -L 19,${BUD}" "Deposits-2016-Mar"
dorrtest  "k3" "${RRDATERANGE} -b ${BUD}" "Finish-2016-Mar"
dorrtest  "l3" "${RRDATERANGE} -b ${BUD} -r 1" "Journal"
dorrtest  "m3" "${RRDATERANGE} -b ${BUD} -r 2" "Ledgers"
dorrtest  "n3" "${RRDATERANGE} -b ${BUD} -r 10" "LedgerActivity"
dorrtest  "o3" "${RRDATERANGE} -b ${BUD} -r 17" "LedgerBalance"
dorrtest  "p3" "${RRDATERANGE} -b ${BUD} -r 8" "Statements"
dorrtest  "q3" "${RRDATERANGE} -b ${BUD} -r 4" "RentRoll"

#========================================================================================
# APRIL 2016
#    GSR and Contract rent change to 4150 for 311 Rexford
#========================================================================================

##----------------------------------------------------------
##  1. Update MarketRate for RentableType 3 to $4150/month
##----------------------------------------------------------
cat >xxyyzz <<EOF
use rentroll
INSERT INTO RentableMarketRate (RTID,MarketRate,DtStart,DtStop) VALUES(3,4150,"2016-04-01 00:00:00","2018-04-01 00:00:00");
EOF
${MYSQL} --no-defaults <xxyyzz
rm -f xxyyzz
dorrtest  "z3" "-j 2016-01-01 -k 2016-06-01 -b ${BUD} -r 20,R003" "MarketRateValidation"

##----------------------------------------------------------
##  2. Process the rent checks and generate reports
##----------------------------------------------------------
RRDATERANGE="-j 2016-04-01 -k 2016-05-01"
CSVLOADRANGE="-G ${BUD} -g 4/1/16,5/1/16"
# docsvtest "b4" "-A asm2016Apr.csv ${CSVLOADRANGE} -L 11,${BUD}" "Assessments-2016-Apr"  		## no new assessments this month
dorrtest  "a4" "${RRDATERANGE} -x -b ${BUD} -r 18" "Process-2016-Apr"
docsvtest "i4" "-e rcpt2016Apr.csv ${CSVLOADRANGE} -L 13,${BUD}" "Receipts-2016-Apr"
docsvtest "j4" "-y deposit-2016Apr.csv ${CSVLOADRANGE} -L 19,${BUD}" "Deposits-2016-Apr"
dorrtest  "k4" "${RRDATERANGE} -b ${BUD}" "Finish-2016-Apr"
dorrtest  "l4" "${RRDATERANGE} -b ${BUD} -r 1" "Journal"
dorrtest  "m4" "${RRDATERANGE} -b ${BUD} -r 2" "Ledgers"
dorrtest  "n4" "${RRDATERANGE} -b ${BUD} -r 10" "LedgerActivity"
dorrtest  "o4" "${RRDATERANGE} -b ${BUD} -r 17" "LedgerBalance"
dorrtest  "p4" "${RRDATERANGE} -b ${BUD} -r 8" "Statements"
dorrtest  "q4" "${RRDATERANGE} -b ${BUD} -r 4" "RentRoll"

#========================================================================================
# MAY 2016
#    GSR and Contract rent change to 3800 for 309.5 Rexford
#========================================================================================

##----------------------------------------------------------
##  1. Update MarketRate for RentableType 2 to $3800/month
##     Update ContractRent for Rentable 2 to $3800/month
##----------------------------------------------------------
cat >xxyyzz <<EOF
use rentroll
INSERT INTO RentableMarketRate (RTID,MarketRate,DtStart,DtStop) VALUES(2,3800,"2016-05-01 00:00:00","2018-01-01 00:00:00");
INSERT INTO RentalAgreementRentables (RAID,RID,CLID,ContractRent,DtStart,DtStop) VALUES(2,2,0,3800,"2016-05-01 00:00:00","2018-03-01 00:00:00");
UPDATE RentalAgreementRentables SET DtStop="2016-05-01" WHERE ContractRent=3550 AND RID=2;
EOF
${MYSQL} --no-defaults <xxyyzz
rm -f xxyyzz
dorrtest  "z4" "-j 2016-01-01 -k 2016-09-01 -b ${BUD} -r 20,R002" "MarketRateValidation"

##----------------------------------------------------------
##  2. Process the rent checks and generate reports
##----------------------------------------------------------
RRDATERANGE="-j 2016-05-01 -k 2016-06-01"
CSVLOADRANGE="-G ${BUD} -g 5/1/16,6/1/16"
docsvtest "b5" "-A asm2016-05.csv ${CSVLOADRANGE} -L 11,${BUD}" "Assessments-2016-May"
dorrtest  "a5" "${RRDATERANGE} -x -b ${BUD} -r 18" "Process-2016-May"
docsvtest "i5" "-e rcpt2016-05.csv ${CSVLOADRANGE} -L 13,${BUD}" "Receipts-2016-May"
docsvtest "j5" "-y deposit-2016-05.csv ${CSVLOADRANGE} -L 19,${BUD}" "Deposits-2016-May"
dorrtest  "k5" "${RRDATERANGE} -b ${BUD}" "Finish-2016-May"
dorrtest  "l5" "${RRDATERANGE} -b ${BUD} -r 1" "Journal"
dorrtest  "m5" "${RRDATERANGE} -b ${BUD} -r 2" "Ledgers"
dorrtest  "n5" "${RRDATERANGE} -b ${BUD} -r 10" "LedgerActivity"
dorrtest  "o5" "${RRDATERANGE} -b ${BUD} -r 17" "LedgerBalance"
dorrtest  "p5" "${RRDATERANGE} -b ${BUD} -r 8" "Statements"
dorrtest  "q5" "${RRDATERANGE} -b ${BUD} -r 4" "RentRoll"

#========================================================================================
# JUNE 2016
#========================================================================================
RRDATERANGE="-j 2016-06-01 -k 2016-07-01"
CSVLOADRANGE="-G ${BUD} -g 6/1/16,7/1/16"
# docsvtest "b6" "-A asm2016-06.csv ${CSVLOADRANGE} -L 11,${BUD}" "Assessments-2016-Jun"  		## no new assessments this month
dorrtest  "a6" "${RRDATERANGE} -x -b ${BUD} -r 18" "Process-2016-Jun"
docsvtest "i6" "-e rcpt2016-06.csv ${CSVLOADRANGE} -L 13,${BUD}" "Receipts-2016-Jun"
docsvtest "j6" "-y deposit-2016-06.csv ${CSVLOADRANGE} -L 19,${BUD}" "Deposits-2016-Jun"
dorrtest  "k6" "${RRDATERANGE} -b ${BUD}" "Finish-2016-Jun"
dorrtest  "l6" "${RRDATERANGE} -b ${BUD} -r 1" "Journal"
dorrtest  "m6" "${RRDATERANGE} -b ${BUD} -r 2" "Ledgers"
dorrtest  "n6" "${RRDATERANGE} -b ${BUD} -r 10" "LedgerActivity"
dorrtest  "o6" "${RRDATERANGE} -b ${BUD} -r 17" "LedgerBalance"
dorrtest  "p6" "${RRDATERANGE} -b ${BUD} -r 8" "Statements"
dorrtest  "q6" "${RRDATERANGE} -b ${BUD} -r 4" "RentRoll"


logcheck
