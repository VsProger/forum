# Переменные
IMAGE_NAME = golang-web-app
CONTAINER_NAME = golang-web-container
PORT = 8081
DOCKER_FILE_PATH = ./Dockerfile

# Сборка Docker-образа
build:
	docker build -f $(DOCKER_FILE_PATH) -t $(IMAGE_NAME) .

# Запуск контейнера
run: build
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):8080 $(IMAGE_NAME)

# Остановка контейнера
stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

# Очистка неиспользуемых образов и контейнеров
clean:
	docker system prune -af

# Запуск тестов (если есть)
test:
	go test ./...

# Перезапуск контейнера (остановить и снова запустить)
restart: stop run

# Открытие оболочки в контейнере
exec:
	docker exec -it $(CONTAINER_NAME) /bin/sh
