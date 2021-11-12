import os
from xml.dom.minidom import parse
import xml.etree.ElementTree as ET
add_index = 2

def cat_xml(bin_dir):
    for root, dirs, files in os.walk(rootdir):
        if 'server' not in root:
            continue
        #print('now dir', root)
        for file in files:
            if 'config.xml' not in file:
                continue
            xml_path = os.path.join(root, file)
            print('xml', xml_path)
            xmldoc = parse(xml_path)
            inet_c_list = xmldoc.getElementsByTagName('inet_c')
            if len(inet_c_list) > 0:
                c_ip = inet_c_list[0].attributes['IP'].value
                c_port = inet_c_list[0].attributes['Port'].value
                print("cIP : ", c_ip)
                print("cPort : ", c_port)
                inet_c_list[0].attributes['Port'].value = str(int(c_port) + add_index)

            inet_s_list = xmldoc.getElementsByTagName('inet_s')
            if len(inet_s_list) > 0:
                s_ip = inet_s_list[0].attributes['IP'].value
                s_port = inet_s_list[0].attributes['Port'].value
                print("sIP : ", s_ip)
                print("sPort : ", s_port)
                inet_s_list[0].attributes['Port'].value = str(int(s_port) + add_index)
            #myfile = open('test_'+file, "w")
            #myfile(xmldoc.toxml())

rootdir = r'D:\auto_explore_dir\bin\Release_1'
cat_xml(rootdir)