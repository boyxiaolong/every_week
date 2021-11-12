
#!bin/sh
for file in ./*
do
    if test -d $file
    then
		rm -rf $file/*.exe
		rm -rf $file/*.xml
		rm -rf $file/*.lib
		rm -rf $file/*.dll
		rm -rf $file/*.pdb
		rm -rf $file/log
		rm -rf $file/*dll*
		rm -rf $file/start.bat
		rm -rf $file/*.txt
		rm -rf $file/behaviorlog
        echo 删除 $file/log
    fi
done