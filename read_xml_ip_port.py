import os
from xml.dom.minidom import parse

def cat_xml(bin_dir):
    for root, dirs, files in os.walk(rootdir):
        if 'server' not in root:
            continue
        #print('now dir', root)
        for file in files:
            if 'config.xml' not in file:
                continue
            print(root,file)
            xml_path = os.path.join(root, file)
            print('xml', xml_path)
            xmldoc = parse(xml_path)
            inet_c_list = xmldoc.getElementsByTagName('inet_c')
            print('inet_c_list len', len(inet_c_list))
            if len(inet_c_list) > 0:
                print("IP : ", inet_c_list[0].attributes['IP'].value)
                print("Port : ", inet_c_list[0].attributes['Port'].value)

            inet_s_list = xmldoc.getElementsByTagName('inet_s')
            print('inet_s_list len', len(inet_s_list))
            if len(inet_s_list) > 0:
                print("IP : ", inet_s_list[0].attributes['IP'].value)
                print("Port : ", inet_s_list[0].attributes['Port'].value)

rootdir = r'D:\auto_explore_dir\bin\Debug'
cat_xml(rootdir)