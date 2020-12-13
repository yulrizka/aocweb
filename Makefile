server:
	AOC_TEMPLATES=cmd/aocweb/templates gin -i -d "cmd/aocweb" run main.go;rm -f gin-bin
