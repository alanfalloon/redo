redo-ifchange rredo

./rredo -h
./rredo --help
./rredo --version

rand_num() {
    shuf -n1 -i0-$1
}

for N in $(seq 100)
do
    FLAGS=$(shuf -n$(rand_num 13) <<EOF
$(shuf -n1 -e -- -d --debug)
$(shuf -n1 -e -- -v --verbose)
$(shuf -n1 -e -- -x --xtrace)
$(shuf -n1 -e -- -k --keep-going)
--overwrite
--log
--only-log
--shuffle
--debug-locks
--debug-pids
--color
--warn-stdout
--main=$(shuf -n1 -e -- redo-sources redo-targets redo-ood redo-stamp redo-always redo-ifcreate redo-ifchange redo-delegate redo-log redo-exec redo-dofile redo)
EOF
)
    TGTS=$(find *.do */ -name '*.do' -type f | sed -e 's/\.do$//' | shuf -n$(rand_num 15))
    ARGS=$(shuf -e -- $FLAGS $TGTS)
    (
        set -x
        ./rredo $ARGS
    )
done
