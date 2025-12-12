# BPS Action Service - DOCX Conversion Guide

## 🎯 Objective
Convert the BPS Action Service Low-Level Design documentation from Markdown to a professional Microsoft Word DOCX format with all PlantUML diagrams rendered as images.

---

## ✅ Recommended Solution: Use Pandoc (Best Quality)

### Step 1: Install Pandoc
```powershell
# Install using winget (Windows Package Manager)
winget install --id=JohnMacFarlane.Pandoc -e
```

### Step 2: Convert to DOCX
```powershell
# Navigate to project directory
cd "d:\Project\Refactor\cbe-supper-app-cps-action"

# Convert Markdown to DOCX
pandoc docs/BPS-ACTION-SERVICE-LLD.md -o docs/BPS-ACTION-SERVICE-LLD.docx --toc --reference-doc=custom-reference.docx
```

### Step 3: Render PlantUML Diagrams
Since PlantUML diagrams are code blocks, you'll need to:

1. **Option A: Use Online PlantUML Editor**
   - Visit: https://www.plantuml.com/plantuml/uml/
   - Copy each PlantUML code block
   - Generate PNG image
   - Insert into Word document

2. **Option B: Use PlantUML CLI**
   ```powershell
   # Install PlantUML
   winget install --id=PlantUML.PlantUML -e
   
   # Generate diagrams (after extracting PlantUML blocks to .puml files)
   java -jar plantuml.jar diagram.puml
   ```

---

## 🔄 Alternative Solution 1: Microsoft Word Direct Import

### Steps:
1. Open Microsoft Word
2. Go to **File → Open**
3. Select `docs/BPS-ACTION-SERVICE-LLD.md`
4. Word will automatically convert the Markdown
5. Manually insert PlantUML diagrams:
   - Extract PlantUML code blocks
   - Generate images at https://www.plantuml.com/plantuml/
   - Insert images into document

---

## 🔄 Alternative Solution 2: Use Python Script (Requires Python)

### Step 1: Install Python
```powershell
winget install --id=Python.Python.3.12 -e
```

### Step 2: Install Required Packages
```powershell
pip install python-docx markdown pillow requests
```

### Step 3: Run Conversion Script
```powershell
python scripts/create_professional_docx.py
```

This will:
- ✅ Extract all PlantUML diagrams
- ✅ Render them as PNG images using PlantUML online service
- ✅ Create a formatted DOCX with embedded images
- ✅ Apply proper styles (headings, tables, lists, code blocks)

---

## 📊 Manual Diagram Extraction

If you prefer to manually handle diagrams, here are the 7 PlantUML diagrams in the document:

### Diagram 1: Service Boundaries
**Location**: Section 2.2  
**Purpose**: Shows external systems, CPS Action Service components, supporting services, and data layer

### Diagram 2: Clean Architecture Layers
**Location**: Section 3.1  
**Purpose**: Illustrates the 4-layer architecture (Presentation, Application, Domain, Infrastructure)

### Diagram 3: Component Interaction Flow
**Location**: Section 3.2  
**Purpose**: Sequence diagram showing maker-checker workflow (creation, approval, rejection)

### Diagram 4: Account Management Flow
**Location**: Section 5.1  
**Purpose**: Link account sequence with Core Banking integration

### Diagram 5: Security Operations Flow
**Location**: Section 5.2  
**Purpose**: User blocking workflow

### Diagram 6: Profile Management Flow
**Location**: Section 5.3  
**Purpose**: Phone number change with dual OTP verification

### Diagram 7: Bank Management Flow
**Location**: Section 5.6  
**Purpose**: Bank creation with MinIO integration

### Diagram 8: Complete Action Lifecycle
**Location**: Section 5.10  
**Purpose**: End-to-end activity diagram

---

## 🎨 Formatting Guidelines for Word

When creating the DOCX manually, apply these styles:

### Headings
- **Heading 1**: Main sections (e.g., "1. Executive Summary")
- **Heading 2**: Subsections (e.g., "1.1 Purpose")
- **Heading 3**: Sub-subsections (e.g., "Business Rules")
- **Heading 4**: Minor headings

### Tables
- Use **Table Grid** or **Light Grid Accent 1** style
- Bold the header row
- Center-align numeric columns

### Code Blocks
- Font: **Courier New** or **Consolas**
- Size: **9-10pt**
- Background: Light gray (#F5F5F5)

### Lists
- Use built-in **Bullet** and **Numbered** list styles
- Maintain consistent indentation

### Images (Diagrams)
- Center-align all diagrams
- Width: 6-6.5 inches
- Add caption below each diagram
- Format: "Figure X: [Description]"

---

## 📝 Quick Start (Easiest Method)

**If you just want a DOCX file quickly:**

1. Install Pandoc:
   ```powershell
   winget install --id=JohnMacFarlane.Pandoc -e
   ```

2. Convert:
   ```powershell
   cd "d:\Project\Refactor\cbe-supper-app-cps-action"
   pandoc docs/BPS-ACTION-SERVICE-LLD.md -o docs/BPS-ACTION-SERVICE-LLD.docx
   ```

3. Open the DOCX in Word

4. For each PlantUML code block:
   - Copy the code
   - Go to https://www.plantuml.com/plantuml/
   - Paste and generate PNG
   - Download and insert into Word document

---

## ✨ Expected Output

The final DOCX should contain:
- ✅ Professional formatting with proper heading hierarchy
- ✅ Well-formatted tables with borders and shading
- ✅ All 8 PlantUML diagrams as high-quality PNG images
- ✅ Proper code block formatting
- ✅ Bulleted and numbered lists
- ✅ Table of contents (auto-generated by Pandoc)
- ✅ Page numbers and headers/footers

---

## 🔧 Troubleshooting

### Issue: Pandoc not found
**Solution**: Restart terminal after installation or add to PATH manually

### Issue: PlantUML diagrams not rendering
**Solution**: Use online service at https://www.plantuml.com/plantuml/

### Issue: Python not found
**Solution**: Install Python 3.12+ and restart terminal

### Issue: Tables not formatting correctly
**Solution**: Use Pandoc with `--reference-doc` option and a custom template

---

## 📞 Support

For issues with:
- **Pandoc**: https://pandoc.org/installing.html
- **PlantUML**: https://plantuml.com/
- **Python**: https://www.python.org/downloads/

---

**Created**: November 29, 2025  
**Status**: Ready for conversion
