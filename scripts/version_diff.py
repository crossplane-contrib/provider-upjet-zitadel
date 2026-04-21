#!/usr/bin/env python3
"""Compare schema versions between two provider schema files."""

import json
import sys


def main():
    if len(sys.argv) != 4:
        print(f"Usage: {sys.argv[0]} <generated.lst> <old-schema.json> <new-schema.json>")
        sys.exit(1)

    generated_lst = sys.argv[1]
    old_schema_path = sys.argv[2]
    new_schema_path = sys.argv[3]

    with open(generated_lst) as f:
        resources = json.load(f)

    with open(old_schema_path) as f:
        old_schema = json.load(f)

    with open(new_schema_path) as f:
        new_schema = json.load(f)

    old_key = next(iter(old_schema.get("provider_schemas", {})))
    new_key = next(iter(new_schema.get("provider_schemas", {})))

    old_resources = old_schema["provider_schemas"][old_key].get("resource_schemas", {})
    new_resources = new_schema["provider_schemas"][new_key].get("resource_schemas", {})

    changes = []
    for r in resources:
        old_ver = old_resources.get(r, {}).get("version", 0)
        new_ver = new_resources.get(r, {}).get("version", 0)
        if old_ver != new_ver:
            changes.append((r, old_ver, new_ver))

    if changes:
        print("Schema version changes detected:")
        for name, old_v, new_v in changes:
            print(f"  {name}: {old_v} -> {new_v}")
        sys.exit(1)
    else:
        print("No schema version changes detected.")


if __name__ == "__main__":
    main()
