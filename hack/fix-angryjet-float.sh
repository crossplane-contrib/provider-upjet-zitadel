#!/usr/bin/env bash
# Workaround for a bug in crossplane-tools/angryjet where single-resolution
# float pointer fields are generated using pkg/ptr instead of pkg/convert.

set -euo pipefail

find . -name 'zz_generated.resolvers.go' | while read -r f; do
  if grep -q 'ptr\.FromFloatPtrValue\|ptr\.ToFloatPtrValue' "$f"; then
    sed -i.bak \
      -e 's/ptr\.FromFloatPtrValue(\([^,]*\), "")/reference.FromFloatPtrValue(\1)/g' \
      -e 's/ptr\.ToFloatPtrValue/reference.ToFloatPtrValue/g' \
      "$f"
    rm -f "${f}.bak"
    echo "fixed: $f"
  fi
done
