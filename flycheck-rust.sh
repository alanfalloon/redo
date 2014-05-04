#! /bin/bash
FILENAME="$(basename "$3")"
FILENAME="${FILENAME#flycheck_}"
ERRFILE=err."$RANDOM"
{
    rustc --no-trans redo.rs &&
    rustc --test --no-trans redo.rs
} 2>&1 |
grep "$FILENAME"':[0-9]\+:' |
sort -u |
sed -e 's/ note: / warning: /' >&2
exit 0
