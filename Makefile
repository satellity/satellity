default: start          # Start local development server	

image: output/image     # Build docker image
bin: output/bin         # Build cross-platform excutable binary
clean: cleanall         # Clean side effects of command compose up, for example: volumns etc.

PLATFORM=local
NAME=satellity-backend

.PHONY: start
start: output/image
	@docker-compose -f deploy/docker-compose.yml up

.PHONY: output/image
output/image:
	@docker build . \
	-f build/backend/Dockerfile \
	--target dev \
	-t ${NAME}

.PHONY: output/bin
output/bin:
	@docker build . \
	-f build/backend/Dockerfile \
	--target bin \
	--output output/ \
	--platform ${PLATFORM} 

.PHONY: cleanall
cleanall:
	@docker-compose -f deploy/docker-compose.yml down --remove-orphans