import os

def collect_go_code(root_dir, output_file):
    with open(output_file, 'w', encoding='utf-8') as outfile:
        for foldername, subfolders, filenames in os.walk(root_dir):
            for filename in filenames:
                if filename.endswith('.go'):
                    filepath = os.path.join(foldername, filename)
                    try:
                        with open(filepath, 'r', encoding='utf-8') as infile:
                            outfile.write(f"\n// File: {filepath}\n")
                            outfile.write(infile.read())
                            outfile.write("\n" + "-"*80 + "\n")
                    except Exception as e:
                        print(f"Failed to read {filepath}: {e}")

if __name__ == '__main__':
    root_directory = input("Enter the path to start searching for Go files: ").strip()
    output_file = 'all_go_code.txt'
    collect_go_code(root_directory, output_file)
    print(f"All Go code has been written to {output_file}")
