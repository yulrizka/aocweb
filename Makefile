server:
	AOC_TEMPLATES=cmd/aocweb/templates gin -i -d "cmd/aocweb" --path "../../" main.go;rm -f gin-bin
