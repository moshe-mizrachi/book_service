import os

def print_directory_structure(start_path, output_file, include_dirs=None, ignore_files=None, ignore_dirs=None):
    if include_dirs is None:
        include_dirs = []
    if ignore_files is None:
        ignore_files = []
    if ignore_dirs is None:
        ignore_dirs = []

    with open(output_file, 'w', encoding='utf-8') as out_file:
        for root, dirs, files in os.walk(start_path):
            # Filter directories to include only specified ones
            dirs[:] = [d for d in dirs if d in include_dirs or d not in ignore_dirs and not d.startswith('.idea')]

            if any(dir in root for dir in include_dirs):
                level = root.replace(start_path, '').count(os.sep)
                indent = ' ' * 4 * level
                out_file.write(f"{indent}{os.path.basename(root)}/\n")

                sub_indent = ' ' * 4 * (level + 1)
                for file in files:
                    if file in ignore_files:
                        continue

                    out_file.write(f"{sub_indent}{file}\n")
                    try:
                        with open(os.path.join(root, file), 'r', encoding='utf-8', errors='ignore') as f:
                            content = f.read()
                            out_file.write(f"{sub_indent}-- CONTENT START --\n")
                            out_file.write(content[:500])  # Write the first 500 characters for brevity
                            out_file.write(f"\n{sub_indent}-- CONTENT END --\n")
                    except Exception as e:
                        out_file.write(f"{sub_indent}-- ERROR READING FILE: {e} --\n")

# Specify the path to your project directory
project_path = '.'  # Change this to the root directory of your project
output_file = 'project_generated_structure.txt'

# Specify directories to include and files/directories to ignore
include_dirs = ['cmd', 'pkg', 'test']
ignore_files = ['go.mod', 'go.sum']
ignore_dirs = ['tmp', 'node_modules', '.idea']

print_directory_structure(project_path, output_file, include_dirs=include_dirs, ignore_files=ignore_files, ignore_dirs=ignore_dirs)
