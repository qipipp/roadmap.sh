# Weather API (Go) + Redis Cache

간단한 날씨 API 래퍼 서버입니다.
- 외부 API: Open-Meteo (API 키 없음)
- 캐시: Redis (TTL 적용)
- 같은 요청을 다시 호출하면 `X-Cache: HIT`로 빠르게 응답

## 준비물
- Go 설치
- Redis 실행 (로컬 127.0.0.1:6379)

## Redis 켜기 (Docker)
docker run -d --name redis -p 6379:6379 redis

## 실행
프로젝트 폴더에서:
go get github.com/redis/go-redis/v9
go run .

## 테스트 (MISS -> HIT 확인)
curl.exe -i "http://localhost:8080/weather?lat=37.5665&lon=126.9780&days=7"
curl.exe -i "http://localhost:8080/weather?lat=37.5665&lon=126.9780&days=7"

## API
GET /weather
- lat: 필수 (float)
- lon: 필수 (float)
- days: 선택 (int, 기본 7, 1~16)

예:
http://localhost:8080/weather?lat=37.5665&lon=126.9780&days=7
