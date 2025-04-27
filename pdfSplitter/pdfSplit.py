import fitz
import os
import glob

downloads_folder = os.path.expanduser("~/Downloads")

files = glob.glob(os.path.join(downloads_folder, "*"))

files.sort(key=os.path.getmtime, reverse=True)


# Print the 10 most recent files
n = 0
files_dict = {}

for file in files[:10]:
    print(str(n) + " " + os.path.basename(file))
    n += 1
    files_dict[n] = file



file = input("Enter selection (as num) should be a .pdf \n: ")


output_folder = "outputs"
os.makedirs(output_folder, exist_ok=True)

if not os.path.exists(output_folder):
    os.makedirs(output_folder)

file = int(file)

#print(files_dict[file])




doc = fitz.open(files_dict[file])
i = 0
page_text = ""

for page_num, page in enumerate(doc, start = 1):
    page_text += page.get_text()

    if page_num % 50 == 0:
        filename = os.path.join(output_folder, f"output_{i}.txt")
        with open(filename, "w") as f:
            f.write(page_text)
        page_text = ""
        i += 1

if page_text:
    filename = os.path.join(output_folder, f"output_{i}.txt")

    with open(filename, "w") as f:
        f.write(page_text)


doc.close