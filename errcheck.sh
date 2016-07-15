#!/bin/bash
ERRORS=$(ls test/*/err.txt 2>/dev/null)
ERRCOUNT=$(ls test/*/err.txt | wc -l)
if (( ERRCOUNT > 0 )); then
	echo "TESTS HAD ERRORS"
	echo "${ERRORS}"
	exit 1
else
	echo "ALL TESTS PASSED"
	echo "Good job!"
fi