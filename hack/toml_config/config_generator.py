import os
import pandas as pd
from collections import deque

structs=[]
output_dir="../../api/v1beta1"
output_file="wso2is_structs.go"

class TreeNode:
    def __init__(self, name):
        self.name = name
        self.children = []
        self.path=""

    def add_child(self, child):
        self.children.append(child)

def create_tree(data):
    root = TreeNode('')
    for item in data:
        node = root
        for name in item.split('.'):
            existing_child = next((child for child in node.children if child.name == name), None)
            if existing_child:
                node = existing_child
            else:
                new_child = TreeNode(name)
                node.add_child(new_child)
                node = new_child
    return root

def generate_go_struct(struct_data):
    struct_name = struct_data['struct_name']
    fields = struct_data['fields']
    
    go_struct = f"type {struct_name} struct {{\n"

    for field in fields:
        field_name = field['name']
        field_type = field['type']
        field_default_val = str(field['default_val'])
        json_key=field['json_key']
        toml_key=field['toml_key']

        if pd.isna(toml_key):
            tail_part=f"\t{field_name} {field_type} `json:\"{json_key},omitempty\"`\n\n"
        else:
            tail_part=f"\t{field_name} {field_type} `json:\"{json_key},omitempty\" toml:\"{toml_key},omitempty\"`\n\n"

        if pd.isna(field['default_val']):
            go_struct += tail_part
        else:
            field_default_val = "\t// +kubebuilder:default:=" + field_default_val
            go_struct += field_default_val + "\n" + tail_part

    go_struct += "}\n\n"
    return go_struct

def write_go_file(go_structs, file_path):
    os.makedirs(output_dir, exist_ok=True)
    with open(file_path, 'w') as file:
        for go_struct in go_structs:
            file.write(go_struct)

def create_new_struct(struct_name):
    struct_name = "Configurations" if struct_name == "" else struct_name
    struct_dict = {
        "struct_name": struct_name,
        "fields": []
    }
    
    return struct_dict

def add_field(struct_dict, field_name, field_type,default_val,json_key, toml_key):
    field_dict = {
        "name": field_name,
        "type": field_type,
        "default_val": default_val,
        "json_key":json_key,
        "toml_key":toml_key
    }
    struct_dict["fields"].append(field_dict)
    
    return struct_dict

def process_content(structs):
    prepend_text="package v1beta1\n\n"
    return prepend_text + "".join(structs)

def traverse_tree(root):
    queue = deque([root])  
    while queue:
        node = queue.popleft()
        filtered_df = df[df['yaml_key'] == node.path]
        go_var = filtered_df['go_var'].values[0] if not filtered_df.empty else None
        go_type = filtered_df['go_type'].values[0] if not filtered_df.empty else None

        if node.path=="":
            go_type="Configurations"

        # The structs themselves
        current_struct = create_new_struct(go_type)
        
        for child in node.children:  
            if (node.name!=""):
                child.path=node.path+"."+child.name   
            else:
                child.path=child.name

            filtered_df = df[df['yaml_key'] == child.path]
            go_var = filtered_df['go_var'].values[0] if not filtered_df.empty else None
            go_type = filtered_df['go_type'].values[0] if not filtered_df.empty else None
            default_val = filtered_df['default_val'].values[0] if not filtered_df.empty else None

            # Struct fields
            json_key = filtered_df['yaml_key'].values[0].split('.')[-1] if not filtered_df.empty else None            
            toml_key = filtered_df['toml_key'].values[0]

            if not (pd.isna(toml_key)):
                toml_key = filtered_df['toml_key'].values[0].split('.')[-1] if not filtered_df.empty else None                
            
            current_struct = add_field(current_struct, go_var, go_type,default_val,json_key, toml_key)

            if (child.children):
                # Add each child to the queue for processing
                queue.append(child)

        go_struct = generate_go_struct(current_struct)      
        structs.append(go_struct)

    content=process_content(structs)

    write_go_file(content, output_dir+"/"+output_file)

csv_file = "./configs.csv"
df = pd.read_csv(csv_file) 

sample_data=df['yaml_key'].values.tolist()

# Create the tree
tree = create_tree(sample_data)

# Traverse the tree
traverse_tree(tree)

