{
  is_job = $10;
  if (is_job) {
    count[$4]++;
  }
}

END {
  for (x in count) {
    printf("%s\t%s\n", x, count[x]);
  }
}
