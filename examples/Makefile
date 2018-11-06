simple/deploy:
	helm upgrade --install --namespace=krane krane-demo-simple ./simple

simple/clean:
	helm delete krane-demo-simple --purge

simple/redeploy: simple/clean simple/deploy