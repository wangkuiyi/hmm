# This script takes input from the output of convert*.awk.  It outputs
# frequency counts of each column into a .cs file in /tmp/.  The
# result from the i-th column is named /tmp/$i.csv.  You can then
# import them into Spreadsheet software for charting and other
# analysis.
#
# Usage:
#
#  gawk -f ../data/convert_v3.awk ../data/pos_and_edu_for_LI_employees_v4.txt | gawk -f stats.awk
#
# Before import outputs into Spreadsheets, you might want to sort
# them.
#
#  mkdir /tmp/sorted; for i in /tmp/*.tab; do cat $i | sort -k2 -g -r > /tmp/sorted/$(basename $i); done
#
# And then you can import /tmp/sorted*.tab to Google Sheet:
#
#  for i in /tmp/sorted/*.tab; do google docs upload $i; done
#
BEGIN {
  header[1] = "entry";
  header[2] = "member";
  header[3] = "begin";
  header[4] = "end";

  header[5] = "company";
  header[6] = "title";
  header[7] = "seniority";
  header[8] = "function";

  header[9] = "school";
  header[10] = "degree";
  header[11] = "degree_rank";
  header[12] = "field";
}

{
  # We cannot use $1 ~ $NF; instead we need to part lines by splitting
  # by "\t"s.
  split($0, sep, "\t")
  for (i = 1; i <= length(sep); i++) {
    if (sep[i] != "") {  # Does not count empty string values.
      counts[i][sep[i]]++
    }
  }
}

END {
  for (i = 1; i <= length(header); i++) {
    print "Output " length(counts[i]) " unique items to /tmp/" header[i] ".tab" > "/dev/stderr"
    for (x in counts[i]) {
      print x, "\t", counts[i][x]  > ("/tmp/" header[i] ".tab")
    }
  }
}
