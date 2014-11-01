BEGIN {
  firstLine = 1
  input = 0
  correct = 0
  error = 0
}

{
  input++

  if (firstLine) {
    firstLine = 0;
  } else {
    entry = $1;
	member = $2;
	begin = $11;
	end = $12;
    is_job = $13;
	is_edu = $14;

    if ($4 == "-9" || $4 == "unknown" || $4 == "?") {
      company_or_school = "";
    } else {
      company_or_school = $4;
    }

    if ($7 == "-9" || $7 == "unknown" || $7 == "?") {
      title_or_degree = "";
    } else {
      title_or_degree = $7;
    }

    if ($8 == "-9" || $8 == "unknown" || $8 == "?") {
      seniority_or_degree_rank = "";
    }      else {
      seniority_or_degree_rank = $8;
    }

    if ($10 == "-9" || $10 == "unknown" || $10 == "?") {
      function_or_field = "";
    } else {
      function_or_field = $10;
    }

    split(begin, begin_segs, "/")
    split(end, end_segs, "/")

    if (length(begin_segs) != 3) {
      print "Error", entry, " failed parse begin" >> "/dev/stderr"
    } else if (length(end_segs) != 3) {
      print "Error", entry, " failed parse end" >> "/dev/stderr"
    } else if (begin_segs[3] > end_segs[3]) {
      print "Error", entry, " end year earlier than begin year" >> "/dev/stderr"
    } else if (is_job == is_edu) {
      print "Error", entry, " is_job == is_edu." >> "/dev/stderr"
      error++;
    } else if (company_or_school == "" && title_or_degree == "" && seniority_or_degree_rank == "" && function_or_field == "") {
      print "Error", entry, " all properties are empty." >> "/dev/stderr"
      error++;
    } else if (begin == "-9" || begin == "unknown" || begin == "?") {
      print "Error", entry, " unknow begin time" >> "/dev/stderr"
      error++;
    } else if (end == "-9" || end == "unknown" || end == "?") {
      print "Error", entry, " unknow end time" >> "/dev/stderr"
      error++;
    } else {
      company = "";
      title = "";
      seniority = "";
      function_ = "";
      school = "";
      degree = "";
      degree_rank = "";
      field = "";

      if (is_job) {
        company = company_or_school;
        title = title_or_degree;
        seniority = seniority_or_degree_rank;
        function_ = function_or_field;
      } else {
        school = company_or_school;
        degree = title_or_degree;
        degree_rank = seniority_or_degree_rank;
        field = function_or_field;
      }

      printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
             entry, member, begin, end,
             company, title, seniority, function_,
             school, degree, degree_rank, field);
      correct++
    }
  }
}

END {
  print "Summary: input", input, "error", error, "correct", correct >> "/dev/stderr"
}
