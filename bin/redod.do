OUT="$PWD"/"$3"
cd ../redod
go build -o "$OUT"
find . -name \*.go -type f -print0 | xargs -0 redo-ifchange
