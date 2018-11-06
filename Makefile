operator: build/operator push/operator deploy/operator

build/operator:
	docker build -t gelidus/krane-operator ./krane-operator

push/operator:
	docker push gelidus/krane-operator

deploy/operator:
	helm upgrade --install operator --namespace=krane ./krane-operator/chart

clean/operator:
	helm delete operator --purge
