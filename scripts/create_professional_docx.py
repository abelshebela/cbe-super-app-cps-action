# BPS Action Service - Professional DOCX Generator
# Converts Markdown to DOCX with PlantUML diagrams as images

import re
import os
import base64
import zlib
import requests
from pathlib import Path

def encode_plantuml(plantuml_text):
    """Encode PlantUML text for URL"""
    # Encode to bytes
    plantuml_bytes = plantuml_text.encode('utf-8')
    
    # Compress
    compressed = zlib.compress(plantuml_bytes)[2:-4]
    
    # Custom base64 encoding for PlantUML
    encoded = base64.b64encode(compressed).decode('utf-8')
    encoded = encoded.translate(str.maketrans('+/', '-_')).rstrip('=')
    
    return encoded

def download_plantuml_diagram(plantuml_code, output_path, diagram_num):
    """Download PlantUML diagram as PNG"""
    try:
        print(f"  📊 Rendering diagram {diagram_num}...")
        
        # Encode PlantUML
        encoded = encode_plantuml(plantuml_code)
        
        # Use PlantUML server
        url = f"http://www.plantuml.com/plantuml/png/{encoded}"
        
        # Download image
        response = requests.get(url, timeout=30)
        response.raise_for_status()
        
        # Save image
        with open(output_path, 'wb') as f:
            f.write(response.content)
        
        print(f"  ✅ Diagram {diagram_num} saved!")
        return True
    except Exception as e:
        print(f"  ⚠️  Failed to render diagram {diagram_num}: {e}")
        return False

def extract_and_render_diagrams(md_file, diagrams_dir):
    """Extract PlantUML diagrams and render them as images"""
    
    # Create diagrams directory
    os.makedirs(diagrams_dir, exist_ok=True)
    
    # Read markdown
    with open(md_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Find all PlantUML blocks
    pattern = r'```plantuml\s*\n(.*?)\n```'
    matches = re.findall(pattern, content, re.DOTALL)
    
    print(f"\n🔍 Found {len(matches)} PlantUML diagrams\n")
    
    # Download each diagram
    diagram_paths = []
    for i, plantuml_code in enumerate(matches, 1):
        diagram_path = os.path.join(diagrams_dir, f'diagram_{i}.png')
        success = download_plantuml_diagram(plantuml_code, diagram_path, i)
        diagram_paths.append(diagram_path if success else None)
    
    return diagram_paths, content

def create_docx_with_images(md_file, output_docx, diagrams_dir):
    """Create DOCX file with embedded diagram images"""
    
    try:
        from docx import Document
        from docx.shared import Inches, Pt, RGBColor
        from docx.enum.text import WD_ALIGN_PARAGRAPH
    except ImportError:
        print("\n❌ python-docx not installed!")
        print("   Run: pip install python-docx")
        return False
    
    # Extract and render diagrams
    diagram_paths, content = extract_and_render_diagrams(md_file, diagrams_dir)
    
    # Create document
    doc = Document()
    
    # Set margins
    for section in doc.sections:
        section.top_margin = Inches(1)
        section.bottom_margin = Inches(1)
        section.left_margin = Inches(1)
        section.right_margin = Inches(1)
    
    print("\n📝 Creating DOCX document...\n")
    
    # Replace PlantUML blocks with image placeholders
    diagram_counter = 0
    lines = content.split('\n')
    
    i = 0
    in_plantuml = False
    in_code_block = False
    in_table = False
    table_data = []
    
    while i < len(lines):
        line = lines[i]
        
        # Handle PlantUML blocks
        if line.strip().startswith('```plantuml'):
            in_plantuml = True
            i += 1
            continue
        
        if in_plantuml:
            if line.strip() == '```':
                in_plantuml = False
                # Add diagram image
                if diagram_counter < len(diagram_paths) and diagram_paths[diagram_counter]:
                    try:
                        p = doc.add_paragraph()
                        p.alignment = WD_ALIGN_PARAGRAPH.CENTER
                        run = p.add_run()
                        run.add_picture(diagram_paths[diagram_counter], width=Inches(6))
                        doc.add_paragraph()  # Add spacing
                        print(f"  ✅ Inserted diagram {diagram_counter + 1}")
                    except Exception as e:
                        print(f"  ⚠️  Could not insert diagram {diagram_counter + 1}: {e}")
                diagram_counter += 1
            i += 1
            continue
        
        # Handle code blocks
        if line.strip().startswith('```') and not in_code_block:
            in_code_block = True
            i += 1
            continue
        
        if line.strip() == '```' and in_code_block:
            in_code_block = False
            doc.add_paragraph()  # Add spacing
            i += 1
            continue
        
        if in_code_block:
            p = doc.add_paragraph(line)
            for run in p.runs:
                run.font.name = 'Courier New'
                run.font.size = Pt(9)
            i += 1
            continue
        
        # Handle headings
        if line.startswith('# '):
            doc.add_heading(line[2:].strip(), level=1)
        elif line.startswith('## '):
            doc.add_heading(line[3:].strip(), level=2)
        elif line.startswith('### '):
            doc.add_heading(line[4:].strip(), level=3)
        elif line.startswith('#### '):
            doc.add_heading(line[5:].strip(), level=4)
        
        # Handle horizontal rules
        elif line.strip() == '---':
            p = doc.add_paragraph('_' * 80)
            p.runs[0].font.color.rgb = RGBColor(200, 200, 200)
        
        # Handle tables
        elif '|' in line and line.strip().startswith('|'):
            if not in_table:
                in_table = True
                table_data = []
            
            # Skip separator line
            if not line.strip().replace('|', '').replace('-', '').replace(' ', ''):
                i += 1
                continue
            
            # Add row data
            row = [cell.strip() for cell in line.split('|')[1:-1]]
            table_data.append(row)
            
        elif in_table and not line.strip().startswith('|'):
            # Create table
            if len(table_data) > 1:
                table = doc.add_table(rows=len(table_data), cols=len(table_data[0]))
                table.style = 'Light Grid Accent 1'
                
                # Fill table
                for row_idx, row_data in enumerate(table_data):
                    for col_idx, cell_text in enumerate(row_data):
                        if col_idx < len(table.rows[row_idx].cells):
                            cell = table.rows[row_idx].cells[col_idx]
                            cell.text = cell_text
                            # Bold first row
                            if row_idx == 0:
                                for paragraph in cell.paragraphs:
                                    for run in paragraph.runs:
                                        run.font.bold = True
                
                doc.add_paragraph()  # Add spacing
            
            in_table = False
            table_data = []
            continue
        
        # Handle lists
        elif line.strip().startswith('- ') or line.strip().startswith('* '):
            text = line.strip()[2:].replace('**', '').replace('✅', '✓').replace('❌', '✗')
            doc.add_paragraph(text, style='List Bullet')
        
        elif re.match(r'^\d+\.', line.strip()):
            text = re.sub(r'^\d+\.\s*', '', line.strip()).replace('**', '')
            doc.add_paragraph(text, style='List Number')
        
        # Handle regular paragraphs
        elif line.strip() and not line.strip().startswith('#'):
            # Handle bold text
            if '**' in line:
                p = doc.add_paragraph()
                parts = re.split(r'(\*\*.*?\*\*)', line)
                for part in parts:
                    if part.startswith('**') and part.endswith('**'):
                        run = p.add_run(part[2:-2])
                        run.font.bold = True
                    else:
                        p.add_run(part)
            else:
                doc.add_paragraph(line.strip())
        
        i += 1
    
    # Save document
    doc.save(output_docx)
    print(f"\n✅ DOCX file created: {output_docx}")
    return True

if __name__ == '__main__':
    md_file = r'd:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.md'
    output_docx = r'd:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.docx'
    diagrams_dir = r'd:\Project\Refactor\cbe-supper-app-cps-action\docs\diagrams'
    
    print("🔄 BPS Action Service Documentation Converter")
    print("=" * 50)
    
    success = create_docx_with_images(md_file, output_docx, diagrams_dir)
    
    if success:
        print("\n✨ Conversion complete!")
        print(f"📄 DOCX file: {output_docx}")
        print(f"📊 Diagrams: {diagrams_dir}")
    else:
        print("\n❌ Conversion failed!")
