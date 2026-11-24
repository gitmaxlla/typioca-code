#!/bin/sh
set -eu

dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
repo=$(dirname "$dir")
cd "$repo/src"

go test ./... -coverprofile="$dir/coverage.out"
go tool cover -func="$dir/coverage.out" > "$dir/coverage-func.txt"

total=$(tail -n 1 "$dir/coverage-func.txt" | awk '{print $3}')

block=$(mktemp)
{
  echo "Обновлено: $(date '+%Y-%m-%d %H:%M %Z')"
  echo
  echo "Итоговое покрытие: **$total**."
  echo "Покрытие по функциям — в \`coverage-func.txt\`, сырой профиль — в"
  echo "\`coverage.out\`."
} > "$block"

readme="$dir/README.md"
out=$(mktemp)
awk -v block="$block" '
  /<!-- coverage:start -->/ {
    print
    while ((getline line < block) > 0) print line
    close(block)
    skip = 1
    next
  }
  /<!-- coverage:end -->/ { skip = 0 }
  !skip { print }
' "$readme" > "$out"
mv "$out" "$readme"
chmod 0644 "$readme"
rm -f "$block"