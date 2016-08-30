#!/bin/bash
TESTNAME="JM1"
TESTSUMMARY="Setup and run JM1 company and the Rexford Properties"

RRDATERANGE="-j 2016-01-01 -k 2016-02-01"

source ../share/base.sh

docsvtest "a" "-b business.csv -L 3" "Business"
docsvtest "b" "-c coa.csv -L 10,${BUD}" "ChartOfAccounts"
docsvtest "c" "-R rentabletypes.csv -L 5,${BUD}" "RentableTypes"
docsvtest "d" "-l strlists.csv -L 25,${BUD}" "StringLists"
docsvtest "e" "-p people.csv  -L 7" "People"

#docsvtest "b" "-u custom.csv -L 14" "CustomAttributes"
#docsvtest "f" "-r rentable.csv -L 6,${BUD}" "Rentables"
#docsvtest "g" "-T ratemplates.csv  -L 8" "RentalAgreementTemplates"
#docsvtest "h" "-C ra.csv -L 9,${BUD}" "RentalAgreements"
#docsvtest "i" "-P pmt.csv -L 12,${BUD}" "PaymentTypes"
#docsvtest "j" "-U assigncustom.csv -L 15" "AssignCustomAttributes"
#docsvtest "k" "-A asmt.csv -L 11,${BUD}" "Assessments"
#docsvtest "l" "-e rcpt.csv -L 13,${BUD}" "Receipts"
#
#dorrtest "m" "-j 2014-12-01 -k 2015-01-01" "ProcessDeposits"
#dorrtest "p" "${RRDATERANGE} -r 11" "GSR"
#dorrtest "m" "${RRDATERANGE} " "Process"
#dorrtest "n" "${RRDATERANGE} -r 1" "Journal"
#dorrtest "o" "${RRDATERANGE} -r 2" "Ledgers"
#dorrtest "q" "-r 12,11,RA001,2016-07-04" "AccountBalance"
#dorrtest "q1" "-r 12,9,RA001,2016-07-04" "AccountBalance"
#dorrtest "r" "${RRDATERANGE} -r 4" "RentRoll"

logcheck