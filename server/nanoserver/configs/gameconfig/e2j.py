import os
import glob
import pandas as pd
import json
import sys

# 检查字符串是否为有效的JSON
def is_valid_json(string):
    try:
        json.loads(string)
        return True
    except (json.JSONDecodeError, TypeError, ValueError):
        return False

# 读取Excel文件并转换为JSON格式，保留原始JSON字段格式
def excel_to_json(excel_path, json_path):
    # 先读取前3行获取列标题
    header_df = pd.read_excel(excel_path, nrows=3, header=None)
    # 获取有效列的索引和名称
    valid_cols = []
    col_names = []
    for i in range(len(header_df.columns)):
        if pd.notna(header_df.iloc[2, i]) and str(header_df.iloc[2, i]).strip():
            valid_cols.append(i)
            col_names.append(str(header_df.iloc[2, i]))
    
    # 读取数据，只保留有效列
    df = pd.read_excel(excel_path, skiprows=2, usecols=valid_cols, names=col_names)
    for col in df.columns:
        df[col] = df[col].apply(lambda x: json.loads(x) if is_valid_json(x) else x)
    df.to_json(json_path, orient='records', force_ascii=False, indent=4)
    print(f"转换完成：{json_path}")
    
# 转换单个Excel文件为JSON
def convert_single_file(excel_path, json_folder_path):
    file_name = os.path.splitext(os.path.basename(excel_path))[0]
    json_path = os.path.join(json_folder_path, f"{file_name}.json")
    excel_to_json(excel_path, json_path)
    print("单个文件转换完成。")
    
# 转换文件夹内所有.xlsx文件为JSON
def convert_folder_to_json(folder_path, json_folder_path):
    for excel_file in glob.glob(os.path.join(folder_path, '*.xlsx')):
        file_name = os.path.splitext(os.path.basename(excel_file))[0]
        json_path = os.path.join(json_folder_path, f"{file_name}.json")
        excel_to_json(excel_file, json_path)
    print("所有文件转换完成。")

# 路径设置
folder_path = './'  
json_folder_path = './json4C/'  

# 确保输出文件夹存在
if not os.path.exists(json_folder_path):
    os.makedirs(json_folder_path)

# 执行转换
# convert_folder_to_json(folder_path, json_folder_path)

# 检查命令行参数
if len(sys.argv) > 1:
    # 如果提供了单个文件路径，则只转换该文件
    single_file_path = sys.argv[1]
    convert_single_file(single_file_path, json_folder_path)
else:
    # 没有提供文件路径，转换文件夹下所有文件
    convert_folder_to_json(folder_path, json_folder_path)