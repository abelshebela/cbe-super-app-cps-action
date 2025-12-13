"""
Convert BPS Action Service LLD Markdown to DOCX
"""
import re
from docx import Document
from docx.shared import Pt, Inches, RGBColor
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.style import WD_STYLE_TYPE

def parse_markdown_to_docx(md_file, docx_file):
    """Convert Markdown file to DOCX with formatting"""
    
    # Create document
    doc = Document()
    
    # Set document margins
    sections = doc.sections
    for section in sections:
        section.top_margin = Inches(1)
        section.bottom_margin = Inches(1)
        section.left_margin = Inches(1)
        section.right_margin = Inches(1)
    
    # Read markdown file
    with open(md_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into lines
    lines = content.split('\n')
    
    in_code_block = False
    in_table = False
    code_language = ''
    table_headers = []
    table_rows = []
    
    i = 0
    while i < len(lines):
        line = lines[i]
        
        # Skip PlantUML blocks
        if line.strip().startswith('```plantuml'):
            in_code_block = True
            i += 1
            while i < len(lines) and not lines[i].strip().startswith('```'):
                i += 1
            in_code_block = False
            i += 1
            continue
        
        # Code blocks
        if line.strip().startswith('```'):
            if not in_code_block:
                in_code_block = True
                code_language = line.strip()[3:]
                i += 1
                continue
            else:
                in_code_block = False
                i += 1
                continue
        
        if in_code_block:
            # Add code with monospace font
            p = doc.add_paragraph(line)
            p.style = 'Normal'
            for run in p.runs:
                run.font.name = 'Courier New'
                run.font.size = Pt(9)
                run.font.color.rgb = RGBColor(0, 0, 0)
            i += 1
            continue
        
        # Headings
        if line.startswith('# '):
            p = doc.add_heading(line[2:], level=1)
        elif line.startswith('## '):
            p = doc.add_heading(line[3:], level=2)
        elif line.startswith('### '):
            p = doc.add_heading(line[4:], level=3)
        elif line.startswith('#### '):
            p = doc.add_heading(line[5:], level=4)
        
        # Horizontal rule
        elif line.strip() == '---':
            doc.add_paragraph('_' * 80)
        
        # Tables
        elif '|' in line and line.strip().startswith('|'):
            if not in_table:
                in_table = True
                table_headers = [cell.strip() for cell in line.split('|')[1:-1]]
                table_rows = []
            elif line.strip().replace('|', '').replace('-', '').strip() == '':
                # Skip separator line
                pass
            else:
                row_data = [cell.strip() for cell in line.split('|')[1:-1]]
                table_rows.append(row_data)
        
        # End of table
        elif in_table and not line.strip().startswith('|'):
            # Create table
            if table_headers and table_rows:
                table = doc.add_table(rows=len(table_rows) + 1, cols=len(table_headers))
                table.style = 'Light Grid Accent 1'
                
                # Add headers
                for j, header in enumerate(table_headers):
                    cell = table.rows[0].cells[j]
                    cell.text = header
                    # Bold headers
                    for paragraph in cell.paragraphs:
                        for run in paragraph.runs:
                            run.font.bold = True
                
                # Add rows
                for row_idx, row_data in enumerate(table_rows):
                    for col_idx, cell_data in enumerate(row_data):
                        if col_idx < len(table.rows[row_idx + 1].cells):
                            table.rows[row_idx + 1].cells[col_idx].text = cell_data
                
                doc.add_paragraph()  # Add spacing after table
            
            in_table = False
            table_headers = []
            table_rows = []
            continue
        
        # Bold text
        elif '**' in line:
            p = doc.add_paragraph()
            parts = re.split(r'(\*\*.*?\*\*)', line)
            for part in parts:
                if part.startswith('**') and part.endswith('**'):
                    run = p.add_run(part[2:-2])
                    run.font.bold = True
                else:
                    p.add_run(part)
        
        # Bullet lists
        elif line.strip().startswith('- ') or line.strip().startswith('* '):
            text = line.strip()[2:]
            # Remove markdown formatting
            text = text.replace('**', '')
            text = text.replace('✅ ', '✓ ')
            text = text.replace('❌ ', '✗ ')
            p = doc.add_paragraph(text, style='List Bullet')
        
        # Numbered lists
        elif re.match(r'^\d+\.\s', line.strip()):
            text = re.sub(r'^\d+\.\s', '', line.strip())
            text = text.replace('**', '')
            p = doc.add_paragraph(text, style='List Number')
        
        # Regular paragraph
        elif line.strip() and not line.strip().startswith('#'):
            # Skip metadata lines
            if ':' in line and line.count(':') == 1 and len(line) < 100:
                parts = line.split(':', 1)
                if parts[0].strip() in ['Project', 'Module', 'Version', 'Date', 'Author']:
                    p = doc.add_paragraph()
                    p.add_run(parts[0] + ': ').font.bold = True
                    p.add_run(parts[1].strip())
                    i += 1
                    continue
            
            doc.add_paragraph(line.strip())
        
        i += 1
    
    # Save document
    doc.save(docx_file)
    print(f"✅ Successfully converted to {docx_file}")

if __name__ == '__main__':
    md_file = r'd:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.md'
    docx_file = r'd:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.docx'
    
    try:
        parse_markdown_to_docx(md_file, docx_file)
    except Exception as e:
        print(f"❌ Error: {e}")
        import traceback
        traceback.print_exc()
