API 정리
=======

Login/Logout 요청
----------------

|Method | URL     | 내용        |
|------|---------|------------|
| POST | /api/login  | 로그인 요청   |
| POST | /api/logout | 로그아웃 요청 |


*로그인 관련 참고*
| 상황 | Code |
|-----|------|
| 정상 로그인                   | statusCode = 200 |
| 로그인 후 POST,PUT,DELETE 호출 | Access-Token cookie 내용 확인 후 정상이면 StatusCode = 200 |
| 로그아웃 후 POST,PUT,DELETE 호출 | Access-Token cookie 내용 제거 후 접속 하면 StatusCode = 401 |


Potofolio
---------
|Method | URL     | 내용        |
|------|-------------|------------|
| GET  | /api/potofolio  | 포토폴리오 목록 요청   |
| GET  | /api/potofolio/:id | :id 해당하는 포토폴리오 내용 요청 |
| DELETE | /api/potofolio/:id | :id 해당하는 포토폴리오 전체 제거 |
| PUT  | /api/potofolio/:id | :id 해당하는 포토폴리오 내용 수정  |
| POST | /api/potofolio | 포토폴리오 추가 |

Essay
---------
|Method | URL     | 내용        |
|------|-------------|------------|
| GET  | /api/essay  | 에세이 목록 요청   |
| GET  | /api/essay/:id | :id 해당하는 에세이 내용 요청 |
| DELETE | /api/essay/:id | :id 해당하는 에세이 전체 제거 |
| PUT | /api/essay/:id | :id 해당하는 에세이 내용 수정  |
| POST | /api/essay | 에세이 내용 추가

About
---------
|Method | URL     | 내용        |
|------|-------------|------------|
| GET  | /api/about  | 내 소개 전체 내용 요청   |
| POST  | /api/about | 내 소개 내용 수정 |
| POST  | /api/about-history | 내 소개  경력 내용 수정 |



