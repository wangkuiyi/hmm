{
  is_job = $10;
  is_edu = $11;

  if (is_job) {
    count[$3]++;
  }
}

END {
  for (x in count) {
    printf("%s\t%s\n", x, count[x]);
  }
}
