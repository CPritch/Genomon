# refactor.py (v4 - Final)
import re
import textwrap

def generate_handler_name(regex_variable_name):
    """Converts a regex variable name like 'healRegex' to 'parseHeal'."""
    base_name = regex_variable_name.replace("Regex", "")
    return f"parse{base_name[0].upper()}{base_name[1:]}"

def transform_handler_body(body, uses_matches):
    """Transforms the inner logic of an 'if' block into a valid Go handler body."""
    # Start by dedenting the entire block to a common baseline
    body = textwrap.dedent(body.strip("\n"))

    # PATTERN 1: Handle simple blocks with a single `strconv.Atoi` call.
    # This is the pattern that was breaking.
    # It finds the declaration line AND the `if err == nil` block together.
    simple_err_pattern = re.compile(
        r"(\w+,\s+err\s+:=\s+strconv\.Atoi\(matches\[\d+\]\))\s+if err == nil {((?:.|\n)*?)\n\s*}",
        re.MULTILINE
    )
    if simple_err_pattern.search(body):
        # Replace the entire structure with a clean version:
        # 1. The original declaration.
        # 2. A new guard clause.
        # 3. The content of the original `if` block.
        body = simple_err_pattern.sub(
            r"\1\nif err != nil {\n\treturn nil\n}\n\2",
            body
        )
    
    # PATTERN 2: Handle complex blocks with multiple `strconv.Atoi` calls (err1, err2).
    multi_err_check = re.search(r"if err1 == nil && err2 == nil", body)
    if multi_err_check:
        # First, find all strconv calls and their error variables
        strconv_calls = re.findall(r"(\w+,\s+(err\d*)\s+:=\s+strconv\.Atoi\(matches\[\d+\]\))", body)
        
        # Replace the multi-if check with just its content
        body = re.sub(r"if err1 == nil && err2 == nil\s*{((?:.|\n)*?)\n\s*}", r"\1", body, flags=re.MULTILINE)
        
        # Now, inject a guard clause after each strconv declaration
        for declaration, err_var in reversed(strconv_calls): # Reverse to avoid index issues
            guard = f"\nif {err_var} != nil {{\n\treturn nil\n}}"
            body = body.replace(declaration, declaration + guard)

    # If a block has no return statement after transformation, it means it was a faulty `if` block
    # that only had logic in the `if` and no `else`. It should return nil on failure.
    if "return []core.Effect" not in body:
        body += "\nreturn nil"

    # Add safety checks for match length at the top
    if uses_matches:
        indices = [int(i) for i in re.findall(r"matches\[(\d+)\]", body)]
        if indices:
            required_len = max(indices) + 1
            body = f"if len(matches) < {required_len} {{\n\treturn nil\n}}\n{body}"

    # Final cleanup of extra newlines and ensure proper indentation
    body = "\n".join(line for line in body.split('\n') if line.strip())
    
    return textwrap.indent(body, "\t")


def main():
    input_filename = "parser.go"
    output_filename = "parser_refactored.go"

    try:
        with open(input_filename, "r", encoding="utf-8") as f:
            content = f.read()
    except FileNotFoundError:
        print(f"Error: Could not find '{input_filename}'.")
        return

    # Extract original blocks
    regex_defs_match = re.search(r"var \((.*?)\n\)", content, re.DOTALL)
    imports_match = re.search(r"import \((.*?)\n\)", content, re.DOTALL)
    parse_func_body_match = re.search(r"func Parse\(text string\) \[\]core\.Effect {((?:.|\n)*?)\n}", content, re.MULTILINE)

    if not all([regex_defs_match, imports_match, parse_func_body_match]):
        print("Error: Failed to parse essential Go blocks.")
        return

    regex_definitions = regex_defs_match.group(1).strip()
    imports = imports_match.group(1).strip()
    parse_func_body = parse_func_body_match.group(1)
    
    processed_regex_vars = set()
    
    if_block_pattern = re.compile(
        r"\n\t// --- .*? ---\n\t+if (?:matches := )?(\w+Regex)\.(?:FindStringSubmatch|MatchString)\(text\).*? {((?:.|\n)*?)\n\t}",
        re.MULTILINE
    )

    handlers_code = []
    parser_entries = []

    # Add a fallback for the first rule which has no preceding comment
    first_rule_pattern = re.compile(
        r"^\s*if (?:matches := )?(\w+Regex)\.(?:FindStringSubmatch|MatchString)\(text\).*? {((?:.|\n)*?)\n\t}",
        re.MULTILINE
    )
    all_matches = first_rule_pattern.findall(parse_func_body) + if_block_pattern.findall(parse_func_body)


    for regex_var, body in all_matches:
        if regex_var in processed_regex_vars:
            continue
        processed_regex_vars.add(regex_var)

        handler_name = generate_handler_name(regex_var)
        # We assume `matches` is always potentially used
        uses_matches = True
        
        transformed_body = transform_handler_body(body, uses_matches)
        
        handler_func = f"""
func {handler_name}(matches []string, text string) []core.Effect {{
{transformed_body}
}}"""
        handlers_code.append(handler_func)
        parser_entries.append(f'\t{{regex: {regex_var}, handler: {handler_name}}},')

    # Assemble the final file
    refactored_content = f"""
package effects

import (
{textwrap.indent(imports, '	')}
)

// REGEX DEFINITIONS
var (
{textwrap.indent(regex_definitions, '	')}
)

// effectParser pairs a regular expression with a function that can parse its matches.
type effectParser struct {{
	regex   *regexp.Regexp
	handler func(matches []string, text string) []core.Effect
}}

// effectParsers holds our list of all known effect parsers.
var effectParsers = []effectParser{{
{"\n".join(parser_entries)}
}}

// Parse iterates through a list of registered effect parsers and returns
// a structured Effect object from the first one that matches the input text.
func Parse(text string) []core.Effect {{
	text = strings.TrimSpace(text)

	for _, p := range effectParsers {{
		if matches := p.regex.FindStringSubmatch(text); len(matches) > 0 {{
			if effects := p.handler(matches, text); effects != nil {{
				return effects
			}}
		}}
	}}

	// If no other rule matches, return an UNKNOWN effect type.
	return []core.Effect{{{{
		Type:        core.EffectUnknown,
		Description: text,
	}}}}
}}

// --- Handler Functions ---
{''.join(handlers_code)}
"""

    final_code = textwrap.dedent(refactored_content).strip()

    with open(output_filename, "w", encoding="utf-8") as f:
        f.write(final_code)

    print(f"✅ Refactoring complete! Check '{output_filename}'.")
    print("Run 'gofmt -w parser_refactored.go' to finalize formatting.")

if __name__ == "__main__":
    main()