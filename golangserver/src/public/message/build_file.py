import sys
import os 
import shutil
import subprocess

def GetFileTile(filename):
  names = filename.split('.', -1)
  if len(names) <= 0:
    return ''

  return names[0]


def CreateProtmsgFileText(file_list):
  print('CreateProtmsgFileText')
  cmd = 'protoc  -I="protomsg/" --go_out="protomsg/" --error_format=msvs '
  for filename in file_list:
   cmd = cmd + ' protomsg/' + filename
  #os.system(cmd)
  
  return cmd
  
  
def BuildProtomsgFile(cmd):
  result = subprocess.call(cmd, shell=True)
  if result != 0 :
    print("Error: failed in compiling : " + cmd)
    return False

  return True


  
def BuildErrorCodeFile():
  print('compling error_code')

  #编译指定文件
  cmd = 'protoc  -I="error_code/" --go_out="error_code/" --error_format=msvs error_code/error_code.proto'
  #os.system(cmd)

  result = subprocess.call(cmd, shell=True)
  if result != 0 :
    print("Error: failed in compiling : error_code.proto")
    return False

  return True
  
def BuildMsgTypeFile():
  print('compling msgtype')

  #编译指定文件
  cmd = 'protoc  -I="msgtype/" --go_out="msgtype/" --error_format=msvs msgtype/msg_type.proto'
  #os.system(cmd)

  result = subprocess.call(cmd, shell=True)
  if result != 0 :
    print("Error: failed in compiling : msgtype.proto")
    return False

  return True
  
def DelProtomsgFiles():
  files = os.listdir('protomsg/')
  for file in files:
    if os.path.isdir(file):
      continue
      
    if file.isdigit():
      continue
    
    filename = file.title().lower()
    
    if filename.endswith('.proto'):
      file_path = 'protomsg/' + filename
      os.remove(file_path);
      continue
      
    if filename.endswith('.go'):
      file_path = 'protomsg/' + filename
      os.remove(file_path);

    
  
def GetAllProtoFiles(file_list):
  files = os.listdir('../../../../../../src/libzeroproto/proto/message/')
  for file in files:
    if os.path.isdir(file):
      continue
      
    if file.isdigit():
      continue
    
    filename = file.title().lower()

    if not filename.endswith('.proto'):
      continue

    if (filename == 'reserved_msg_type.proto'):
      continue
      
    if (filename == 'msg_type.proto'):
      continue
      
    if (filename == 'error_code.proto'):
      continue
     
    if (filename == 'msgdefine_server.proto'):
      continue
    
    if (filename == 'chat_bridge_msgs.proto'):
      continue

    if (filename == 'chat_common_msgs.proto'):
      continue


    file_list.append(filename)
    
def DeployErrorCode():
  error_code = 'error_code.proto'
  print('deploying ' + error_code)

  folder_path = '../../../../../../src/libzeroproto/proto/message/'
  if (os.path.exists(folder_path)):
    shutil.copy(folder_path + error_code,  "error_code/" + error_code)
    
def DeployMsgType():
  msgtype = 'msg_type.proto'
  print('deploying ' + msgtype)

  folder_path = '../../../../../../src/libzeroproto/proto/message/'
  if (os.path.exists(folder_path)):
    shutil.copy(folder_path + msgtype,  "msgtype/" + msgtype)
    
def DeployProtoMsg(file_list):
  print('deploying protomsg')
  folder_path = '../../../../../../src/libzeroproto/proto/message/'
  
  for filename in file_list:
   shutil.copy(folder_path + filename,  "protomsg/" + filename)

def main():
  DelProtomsgFiles()
  DeployErrorCode()
  BuildErrorCodeFile()
  DeployMsgType()
  BuildMsgTypeFile()
  
  
  file_list = []
  GetAllProtoFiles(file_list)
  DeployProtoMsg(file_list)
  cmd = CreateProtmsgFileText(file_list)
  if not BuildProtomsgFile(cmd):
      print('Error: failed in building protomsg')
  
if __name__ == '__main__':
  main()
  

  

