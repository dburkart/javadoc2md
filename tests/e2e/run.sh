#!/bin/bash

OUTPUT_DIR=$(mktemp -d)

go run cmd/javadoc2md/main.go -input tests/e2e/input -output $OUTPUT_DIR

if [[ $SHOULD_REBASE ]]; then
	rm tests/e2e/expectations/*
	cp $OUTPUT_DIR/* tests/e2e/expectations
fi

diff -bur tests/e2e/expectations $OUTPUT_DIR

DIFF_STATUS=$?

rm -Rf $OUTPUT_DIR
exit $DIFF_STATUS
