{
  if ($10) {
    count[$5]++;
  }
}

END {
  for (x in count) {
    printf("%s\t%s\n", x, count[x]);
  }
}
