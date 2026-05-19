#!/usr/bin/env python3
"""
Transform Go files to use session-scoped logger from context.

For each method that:
1. Has `ctx context.Context` as a parameter
2. Has X.logger.Xxxf(...) calls

Add `log := local_util.LoggerFromCtx(ctx, X.logger)` at the top
and replace X.logger.Xxxf( with log.Xxxf(
"""

import re
import sys
import os

# Import needed
IMPORT_TO_ADD = 'local_util "cbe-super-app-cps-action/pkgs/utils"'
IMPORT_PATH = '"cbe-super-app-cps-action/pkgs/utils"'

def has_import(content: str) -> bool:
    """Check if local_util is already imported."""
    return 'local_util "cbe-super-app-cps-action/pkgs/utils"' in content

def add_import(content: str) -> str:
    """Add the local_util import to the import block."""
    # Find the import block
    # Try multi-line import first
    multi_import_match = re.search(r'(import\s*\()', content)
    if multi_import_match:
        # Find the closing paren of the import block
        start = multi_import_match.start()
        # Find position of the opening paren
        paren_pos = content.index('(', start)
        # Insert before the closing paren
        close_paren = content.index(')', paren_pos)
        # Add a tab + import before closing paren
        new_content = content[:close_paren] + f'\t{IMPORT_TO_ADD}\n' + content[close_paren:]
        return new_content

    # Single import statement
    single_import = re.search(r'import\s+"[^"]+"\s*\n', content)
    if single_import:
        # Convert to multi-line import
        original = single_import.group(0)
        existing_pkg = re.search(r'"[^"]+"', original).group(0)
        replacement = f'import (\n\t{existing_pkg}\n\t{IMPORT_TO_ADD}\n)\n'
        return content[:single_import.start()] + replacement + content[single_import.end():]

    return content

def find_methods_needing_transform(content: str):
    """
    Parse Go file content to find methods that need transformation.
    Returns list of (method_start, method_end, receiver_var, method_body_start)
    """
    # Match method declarations with ctx context.Context
    # func (receiver *Type) MethodName(ctx context.Context, ...) ... {
    method_pattern = re.compile(
        r'func\s+\((\w+)\s+[*]?\w+\)\s+\w+\s*\([^)]*\bctx\s+context\.Context\b[^)]*\)\s*[^{]*\{'
    )

    results = []
    for m in method_pattern.finditer(content):
        receiver_var = m.group(1)
        func_open_brace = m.end() - 1  # position of '{'
        # Now scan for the matching close brace
        brace_count = 1
        pos = func_open_brace + 1
        while pos < len(content) and brace_count > 0:
            if content[pos] == '{':
                brace_count += 1
            elif content[pos] == '}':
                brace_count -= 1
            pos += 1

        method_end = pos  # after the closing brace
        method_body_start = func_open_brace + 1
        results.append((m.start(), method_end, receiver_var, method_body_start, m.group(0)))

    return results


def transform_file(filepath: str) -> tuple[bool, str]:
    """
    Transform a single Go file.
    Returns (was_modified, reason_if_skipped)
    """
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    # Find methods needing transform
    methods = find_methods_needing_transform(content)

    if not methods:
        return False, "no methods with ctx context.Context"

    # Check which receiver vars have logger calls
    # We process methods in reverse order so offsets stay valid
    modified = False

    # We'll rebuild the content by processing each method
    # Process in reverse order to preserve positions
    for method_start, method_end, receiver_var, body_start, func_sig in reversed(methods):
        method_body = content[body_start:method_end-1]

        # Check if this receiver var has logger calls in this method
        logger_pattern = re.compile(
            rf'\b{re.escape(receiver_var)}\.logger\.(Infof|Errorf|Warnf|Debugf|Fatalf)\('
        )

        if not logger_pattern.search(method_body):
            continue  # No logger calls in this method

        # Check if already transformed
        already_pattern = re.compile(
            rf'log\s*:=\s*local_util\.LoggerFromCtx\s*\(ctx\s*,\s*{re.escape(receiver_var)}\.logger\s*\)'
        )
        if already_pattern.search(method_body):
            continue  # Already transformed

        # Find the start of the method body (skip newline after {)
        # Insert after the opening brace + newline
        insert_pos = body_start
        # Skip whitespace/newline after {
        while insert_pos < len(content) and content[insert_pos] in ('\n', '\r'):
            insert_pos += 1

        # Find the indentation of the first line of the method body
        # Look for the newline before the first non-empty line
        first_line_start = insert_pos
        indent = '\t'
        # Try to detect actual indentation
        first_line_end = content.find('\n', first_line_start)
        if first_line_end > first_line_start:
            first_line = content[first_line_start:first_line_end]
            indent_match = re.match(r'^(\s+)', first_line)
            if indent_match:
                indent = indent_match.group(1)

        # Insert log line at the beginning of the method body
        log_line = f'\n{indent}log := local_util.LoggerFromCtx(ctx, {receiver_var}.logger)\n'

        # We insert right after the opening brace
        content = content[:body_start] + log_line + content[body_start:]

        # Now replace X.logger.Xxxf( with log.Xxxf( in this method
        # But we need to adjust positions since we just inserted
        new_method_end = method_end + len(log_line)

        # Replace in the modified method body
        method_section = content[body_start:new_method_end]
        new_section = logger_pattern.sub(
            lambda m: f'log.{m.group(1)}(',
            method_section
        )
        content = content[:body_start] + new_section + content[new_method_end:]

        modified = True

    if not modified:
        return False, "no methods needed transformation"

    # Add import if needed
    if not has_import(content):
        content = add_import(content)

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        return True, ""

    return False, "no changes made"


def process_files(file_list: list[str]) -> dict:
    """Process all files and return stats."""
    stats = {
        'modified': [],
        'skipped': [],
        'errors': []
    }

    for filepath in file_list:
        try:
            modified, reason = transform_file(filepath)
            if modified:
                stats['modified'].append(filepath)
            else:
                stats['skipped'].append((filepath, reason))
        except Exception as e:
            stats['errors'].append((filepath, str(e)))

    return stats


if __name__ == '__main__':
    import glob as _glob

    base = "d:/Project/Refactor/support/cbe-super-app-cps-action"

    # Collect service files
    service_files = []
    for root, dirs, files in os.walk(os.path.join(base, 'internal/service')):
        for f in files:
            if f.endswith('.go'):
                service_files.append(os.path.join(root, f))

    # Collect storage files excluding generated dirs
    storage_files = []
    skip_dirs = [
        os.path.join(base, 'internal/storage/persistance/sitota/sqlc'),
        os.path.join(base, 'internal/storage/persistance/bank/gen'),
        os.path.join(base, 'internal/storage/persistance/vault/gen'),
    ]

    for root, dirs, files in os.walk(os.path.join(base, 'internal/storage/persistance')):
        # Check if we should skip this dir
        skip = False
        for sd in skip_dirs:
            if root.startswith(sd):
                skip = True
                break
        if skip:
            continue
        for f in files:
            if f.endswith('.go'):
                storage_files.append(os.path.join(root, f))

    all_files = service_files + storage_files
    print(f"Processing {len(all_files)} files ({len(service_files)} service + {len(storage_files)} storage)")

    stats = process_files(all_files)

    print(f"\n=== RESULTS ===")
    print(f"Modified: {len(stats['modified'])} files")
    print(f"Skipped: {len(stats['skipped'])} files")
    print(f"Errors: {len(stats['errors'])} files")

    if stats['modified']:
        print("\n--- MODIFIED ---")
        for f in stats['modified']:
            print(f"  {f}")

    if stats['errors']:
        print("\n--- ERRORS ---")
        for f, err in stats['errors']:
            print(f"  {f}: {err}")

    if '--verbose' in sys.argv and stats['skipped']:
        print("\n--- SKIPPED ---")
        for f, reason in stats['skipped']:
            print(f"  {os.path.basename(f)}: {reason}")
