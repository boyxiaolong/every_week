import os
from xml.dom.minidom import parse
import xml.etree.ElementTree as ET
import socket

add_index = 0

def cat_xml(bin_dir, res_dir):
    local_ip = socket.gethostbyname(socket.gethostname())
    for root, dirs, files in os.walk(rootdir):
        if 'server' not in root:
            continue
        #print('now dir', root)
        for file in files:
            if 'config.xml' not in file:
                continue
            xml_path = os.path.join(root, file)
            print('xml path', root, file)
            #xmldoc = ET.ElementTree(xml_path)
            xmldoc = parse(xml_path)
            inet_c_list = xmldoc.getElementsByTagName('inet_c')
            if len(inet_c_list) > 0:
                c_ip = inet_c_list[0].attributes['IP'].value
                c_port = inet_c_list[0].attributes['Port'].value
                inet_c_list[0].attributes['Port'].value = str(int(c_port) + add_index)
                inet_c_list[0].attributes['IP'].value = local_ip

            inet_s_list = xmldoc.getElementsByTagName('inet_s')
            if len(inet_s_list) > 0:
                s_ip = inet_s_list[0].attributes['IP'].value
                s_port = inet_s_list[0].attributes['Port'].value
                inet_s_list[0].attributes['Port'].value = str(int(s_port) + add_index)
                inet_s_list[0].attributes['IP'].value = local_ip
            
            if 'loginserver' in root or 'gameserver' in root:
                enet_list = xmldoc.getElementsByTagName('enet')
                if len(enet_list) > 0:
                    enet_list[0].attributes['IP'].value = local_ip
                    print('set login eip', local_ip)
                
                if 'gameserver' in root:
                    GSAddr_list = xmldoc.getElementsByTagName('GSAddr')
                    if len(GSAddr_list) > 0:
                        GSAddr_list[0].attributes['domain'].value = local_ip
                        print('set game domain', local_ip)


            new_str = xmldoc.toxml()
            #xmldoc.close()
            #print(new_str)
            #new_file_name = file
            sub_dir = root[root.rfind('\\')+1:len(root)]
            print('sub_dir', sub_dir)
            new_file_path = os.path.join(res_dir, sub_dir)
            new_file_path = os.path.join(new_file_path, file)
            print('new_file_path', new_file_path)
            myfile = open(new_file_path, "w")
            myfile.write(new_str)
            myfile.close()

rootdir = r'D:\test_config\server_config'
res_dir = r'D:\test_config\server_config_1'

cat_xml(rootdir, res_dir)