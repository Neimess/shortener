
doc:
	cd ${WORKDIR}
	swag init -g "./cmd/shortener/main.go" -o ./docs

run:
	bash scripts/run.sh