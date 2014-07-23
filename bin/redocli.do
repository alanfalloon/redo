redo-ifchange redod
OUT="$PWD"/"$3"
cd ../redocli
go build -o "$OUT"
find . -name \*.go -type f -print0 | xargs -0 redo-ifchange
