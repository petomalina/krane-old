
build/operator:
	docker build -t gelidus/krane-operator ./operator

push/operator:
	docker push gelidus/krane-operator