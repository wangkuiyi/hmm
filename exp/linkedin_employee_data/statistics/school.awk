{
  is_edu = $11;
  if (is_edu) {
    count[$3]++;
  }
}

END {
  for (x in count) {
    printf("%s\t%s\n", x, count[x]);
  }
}
