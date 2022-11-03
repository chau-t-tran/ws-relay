.PHONY: start build deploy start test cover

export image := `aws lightsail get-container-images --service-name ws-to-me | jq -r '.containerImages[0].image'`

start:
	go run main.go

build:
	docker build -t ws-to-me .

deploy:
	aws lightsail push-container-image --service-name ws-to-me --label app --image ws-to-me 
	aws lightsail create-container-service-deployment --service-name ws-to-me \
		--containers '{"app":{"image":"'$(image)'","environment":{"HOST":"","PORT":"5000","LOG_ENV":"production"},"ports":{"5000":"HTTP"}}}' \
		--public-endpoint '{"containerName":"app","containerPort":5000,"healthCheck":{"path":"/health"}}'

test:
	go test -coverprofile=cover.out ./...

cover:
	make test
	go tool cover -html=cover.out -o cover.html
	open cover.html
