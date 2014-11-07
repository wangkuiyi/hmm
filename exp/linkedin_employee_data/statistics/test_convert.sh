gawk -f convert.awk test.txt > 1.csv
if diff 1.csv test.result.csv > /tmp/null; then
	S=$( grep linkedin 1.csv )
	if [ "$S" != "" ]; then 
		echo "company dominaint linkedin is still there";
		break;
	fi		
	T=$(grep sales 1.csv )
	if [ "$T" != "" ]; then
		echo "function dominaint sales is still there";
		break;
	fi		
	if gawk -F'\t'  '{if($7!="" || $11!="") print $7, $11}' 1.csv > /tmp/null; then
        	echo "test passed";
	else
		echo "$7 and $11 should be empty"
	fi
else
	echo "output does not match";
fi

rm 1.csv



