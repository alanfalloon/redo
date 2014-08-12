#! /bin/bash

sqlite3 .redo.db <<EOF
.dump
.exit
EOF
