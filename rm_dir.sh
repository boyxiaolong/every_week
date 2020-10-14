dir_re="2020_"
max_month=9
add="0"
end="_*"
for i in $(seq 1 $max_month)
do
   str=$dir_re$add$i
   if [ $i -gt 9 ]
   then
       str=$dir_re$i
   fi
   str=$str$end
   rm -rf $str&
   echo "delete $str"
done
wait
echo "delete finish"