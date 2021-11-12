set ASAN_SAVE_DUMPS=server.dmp
set ASAN_OPTIONS=new_delete_type_mismatch=0
cd activityserver
start activityserver.exe
cd ..

cd databaseserver
start databaseserver.exe
cd ..

cd dispatcherserver
start dispatcherserver.exe
cd ..


cd globalserver
start globalserver.exe
cd ..


cd gameserver
start gameserver.exe
cd ..

cd guildsearchserver
start guildsearchserver.exe
cd ..

cd guildserver
start guildserver.exe
cd ..

cd instanceserver
start instanceserver.exe
cd ..

cd kvkserver
start kvkserver.exe
cd ..

cd localuserserver
start localuserserver.exe
cd ..

cd loginserver
start loginserver.exe
cd ..

cd mailserver
start mailserver.exe
cd ..

cd mapserver
start mapserver.exe
cd ..

cd pathfindserver
start pathfindserver.exe
cd ..

cd transferserver
start transferserver.exe
cd ..

cd userserver
start userserver.exe
cd ..


cd globalactivityserver
start globalactivityserver.exe
cd ..

cd guildderbyserver
start guildderbyserver.exe
cd ..

cd reportserver
start reportserver.exe
cd ..
exit
