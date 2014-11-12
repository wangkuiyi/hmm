out=$(gawk -f convert.awk test.txt) 
if diff <(echo "$out") \
   <(cat <<EOF
393397482	105509708	5/1/2011	8/1/2011	linkedin	intern	3	engineering				
719403537	105509708	6/1/2009	5/1/2010	uc-berkeley	research-assistant	3	research				
725180647	105509708	6/1/2012	10/30/2014	linkedin	software-engineer	3	engineering				
1091821488	103275950	10/1/2013	10/30/2014	linkedin	senior-software-engineer	4	engineering				
385629752	103275950	5/1/2011	8/1/2011	linkedin	software-engineering-intern	2	engineering				
366804234	103275950	6/1/2010	5/1/2011	cisco	software-engineer	3	engineering				
366802272	103275950	6/1/2009	8/1/2009	cisco	software-engineer	3	engineering				
366815133	103275950	6/1/2008	8/1/2008	cisco	software-engineering-intern	2	engineering				
448433441	103275950	9/1/2011	10/1/2013	linkedin	software-engineer	3	engineering				
388122029	63143482	9/1/2010	5/1/2011	uc-berkeley	reader	4	education				
389807717	63143482	5/1/2009	8/1/2009	big-city-chefs-inc.	web-development-intern	2	information-technology				
1391922035	63143482	10/1/2014	10/30/2014	linkedin	senior-software-engineer	4	engineering				
276943251	63143482	5/1/2010	8/1/2010	blackberry	student-software-developer	2	engineering				
388120632	63143482	5/1/2011	12/1/2011	linkedin	software-development-intern	2	engineering				
738081298	63143482	7/1/2012	9/1/2014	linkedin	software-engineer	3	engineering				
711905352	99944053	5/1/2012	8/1/2012	linkedin	software-engineering-intern	2	engineering				
634092722	99944053	1/1/2012	6/1/2012	coursegain	software-engineer	3	engineering				
1238668086	99944053	1/1/2014	10/30/2014	linkedin	software-engineer	3	engineering				
904419021	99944053	2/1/2013	12/1/2013	linkedin	associate-software-engineer	3	engineering				
302678603	63631912	2/1/2010	4/1/2012	haas-school-of-business	assistant-web-developer	3	information-technology				
302678610	63631912	5/1/2010	8/1/2010	reputationdefender	software-engineering-intern	2	engineering				
387291976	63631912	5/1/2011	8/1/2011	linkedin	software-engineering-intern	2	engineering				
742323568	63631912	7/1/2012	10/30/2014	linkedin	software-engineer	3	engineering				
EOF
) > /dev/null; then
    compare=$((echo "$out") | gawk -F'\t'  '{if(NF!=12 || $9!="" || $10!="" || $11!="" || $12!="") print "abc"}')
    if [[ -z "$compare" ]];then
            echo "PASSED";
    else
            echo "FAILED: last 4 columns should be empty"
    fi
else
    echo "FAILED: output does not match expectation!"
fi
