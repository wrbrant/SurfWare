import json
import os
import sys

f = sys.argv[1] if len(sys.argv) > 1 else None



def parse_json_file(filename):
    everything = ""
   # Check if file exists
    if not os.path.isfile(filename):
        return(f"ERROR: {filename} does not exist.")

    with open(filename, 'r') as file:
        data = json.load(file)
        if not 'properties' in data:
            if 'allOf' in data:
                for elem in data['allOf']:
                    for e in elem:
                        if e == 'properties':
                            props = elem[e]
            else:
                return "No properties in allOf"
        else:
            props = data['properties']
            
        
        for prop in props:                
            if 'description' in props[prop]:
                desc = props[prop]['description']
                everything += f'        {prop}: "{desc}",\n'
            else:
                if prop != "id":
                    everything += f'        {prop}: ATTN: NO DESCRIPTION\n'

    return everything

    
def compare_paths():
    # Get the names of every file in a directory that ends in.json
    json_files = [file for file in os.listdir(os.getcwd()) if file.endswith('.json')]
    for file in json_files:
        old = parse_json_file(os.getcwd()+"/"+file)
        new = parse_json_file("/home/branwr/projects/CSR/small-typo-fix/schemas/"+file)
        if old != new:
            print(f"**{file}**\nold: {old}\nnew: {new}\n")

if __name__ == '__main__':    
    # parse_json_file(f)
    # print(compare_paths())
    print(":{\n"+parse_json_file(f)+"    },", end="")

