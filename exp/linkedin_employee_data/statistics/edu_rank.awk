{
  is_edu = $11;
  if (is_edu) {
    count[$4]++;
  }
}

END {
  for (x in count) {
    printf("%s\t%s\n", x, count[x]);
  }
}
