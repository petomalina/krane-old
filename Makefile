timestamp := $(shell /bin/date '+%s')

operator: build/operator push/operator deploy/operator

build/operator:
	docker build -t gelidus/krane-operator:$(timestamp) ./krane-operator

push/operator:
	docker push gelidus/krane-operator:$(timestamp)

deploy/operator:
	helm upgrade --install operator --set image=gelidus/krane-operator:$(timestamp) --namespace=krane ./krane-operator/chart

clean/operator:
	helm delete operator --purge
