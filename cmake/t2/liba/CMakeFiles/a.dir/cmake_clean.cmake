FILE(REMOVE_RECURSE
  "CMakeFiles/a.dir/liba.o"
  "liba.pdb"
  "liba.so"
)

# Per-language clean rules from dependency scanning.
FOREACH(lang C)
  INCLUDE(CMakeFiles/a.dir/cmake_clean_${lang}.cmake OPTIONAL)
ENDFOREACH(lang)
