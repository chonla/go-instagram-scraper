build:
	docker build -t chonla/instagram-scraper:1.0.0 -t chonla/instagram-scraper:latest .

push:
	docker push chonla/instagram-scraper:1.0.0
	docker push chonla/instagram-scraper:latest