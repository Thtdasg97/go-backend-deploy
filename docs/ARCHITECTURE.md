# Tài liệu Kiến trúc — go-backend

## 1. Tổng quan

`go-backend` là một REST API viết bằng Go, dùng [Gin](https://github.com/gin-gonic/gin), theo kiến trúc phân lớp (layered, kiểu "clean architecture") với dependency injection thủ công qua constructor. Dữ liệu được lưu vào MySQL thông qua hai ORM song song: [ent](https://entgo.io) (schema-first, sinh code tự động) và [GORM](https://gorm.io) (dùng ở một số chỗ, ví dụ nhánh `CreateGorm` trong article repository).

Repo này nằm trong một workspace lớn hơn phục vụ học tập/capstone (`/Golang`), bên trong còn có frontend Next.js ([`src_FE/go-frontend`](../../../src_FE/go-frontend)) và các bài tập Go độc lập ([`src/GO_1/Lesson_List`](../../Lesson_List)). Tài liệu này chỉ tập trung vào backend service; xem [§7](#7-frontend--bối-cảnh-toàn-workspace) để hiểu backend khớp với phần còn lại của workspace như thế nào.

> 📜 Lịch sử cập nhật tài liệu này: xem [CHANGELOG.md](CHANGELOG.md).

## 2. Tech stack

| Thành phần | Lựa chọn |
|---|---|
| Ngôn ngữ | Go 1.26 |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM chính | [ent](https://entgo.io) — schema-first, sinh code tự động |
| ORM phụ | [GORM](https://gorm.io) — dùng cho một phần nhánh ghi dữ liệu |
| Database | MySQL (`gorm.io/driver/mysql`, dialect MySQL của ent) |
| Auth | JWT (`golang-jwt/jwt/v5`), băm mật khẩu bằng bcrypt |
| Config | File `.env` qua `joho/godotenv` |
| CORS | `gin-contrib/cors` |
| Hot reload (dev) | [`air`](https://github.com/air-verse/air) (`.air.toml`, build ra `./tmp/main`) |

## 3. Kiến trúc phân lớp

Request đi qua 5 lớp, mỗi lớp chỉ đảm nhiệm một trách nhiệm duy nhất. Các lớp giao tiếp với nhau qua interface, và toàn bộ implementation cụ thể được nối (wire) tại một chỗ duy nhất ([`internal/di/di.go`](../internal/di/di.go)):

```
HTTP request
     │
     ▼
┌─────────────┐   đăng ký route Gin trên *gin.RouterGroup,
│  delivery   │   không chứa business logic (article.delivery.go, auth.delivery.go, ...)
└─────────────┘
     │
     ▼
┌─────────────┐   parse/validate request (query, body, param, header) vào DTO,
│   handler    │   map lỗi sang response.Exception, gọi usecase, ghi response.Success
└─────────────┘
     │
     ▼
┌─────────────┐   business logic: quy tắc validate, điều phối giữa các repository,
│   usecase    │   băm mật khẩu, cấp token, tính toán phân trang
└─────────────┘
     │
     ▼
┌─────────────┐   chỉ truy xuất dữ liệu — build query ent/GORM, không chứa business rule
│ repository   │   (interface trong internal/repository, implementation trong repository_impl)
└─────────────┘
     │
     ▼
┌─────────────┐   ent client / GORM client → MySQL
│  ent / gorm  │
└─────────────┘
```

Mỗi lớp được khởi tạo bằng hàm constructor thuần Go (`NewXxx(...)`) nhận dependency dưới dạng interface — không dùng framework/container DI nào; `internal/di/di.go` chính là "composition root", tự tay wire mọi thứ lúc khởi động rồi gọi `RegisterRouter` trên `rootDelivery` đã được build hoàn chỉnh.

### Sơ đồ thư mục

```
cmd/                     entrypoint (main.go) + khởi tạo App (app.go)
internal/
  di/                    wiring dependency thủ công (composition root)
  delivery/              đăng ký route theo từng resource (*.delivery.go)
  handler/                xử lý request/response ở tầng HTTP (*.handler.go)
  usecase/               interface business logic
    usecase_impl/        implementation business logic
  repository/            interface truy xuất dữ liệu
    repository_impl/     implementation dựa trên ent/GORM
  dto/                   struct định nghĩa hình dạng request/response, bind từ JSON/query
  common/
    env/                 load .env vào struct Env có kiểu
    ent_client/           constructor cho ent.Client
    gorm_client/          constructor cho gorm.DB
    middlewares/          middleware của Gin (xử lý lỗi tập trung, demo truyền context)
    response/             envelope success/error + các exception HTTP có kiểu
    pagination/            parse page/pageSize/offset + hình dạng response phân trang
    models/                 struct map với GORM (song song với model do ent sinh ra)
ent/                     code do ent sinh ra + schema/ (nguồn chân lý của DB schema)
docs/                    tài liệu này + các ghi chú riêng về ent
```

## 4. Vòng đời một request

Lấy `POST /api/article` làm ví dụ cụ thể:

1. **`cmd/main.go`** gọi `NewApp()` — hàm này build Gin engine, đăng ký middleware toàn cục (logger, `ErrorHandler`, custom panic recovery, CORS), tạo client ent + GORM, rồi gọi `dependency.Injection(...)` để wire và đăng ký toàn bộ route. Sau đó `app.Start()` gọi `ginEngine.Run(host:port)`.
2. **`articleDelivery.RegisterRouter`** đã gắn `POST /api/article` với `articleHandler.Create`.
3. **`ArticleHandler.Create`** bind JSON body vào `dto.ArticleCreateReq` (`ShouldBindJSON`, các tag `binding:"required"` thực hiện validate). Nếu bind lỗi, handler gọi `ctx.Error(response.NewBadRequestException(...))` rồi `return`.
4. **`articleUsecase.Create`** nhận DTO, chuyển thẳng xuống repository (việc tạo article hiện chưa có business rule nào khác ngoài việc bọc lỗi).
5. **`articleRepository.Create`** build mutation `Articles.Create()` của ent, set field một cách có điều kiện cho các pointer optional, rồi gọi `.Save(ctx)` để ghi vào MySQL qua ent client.
6. Luồng điều khiển quay ngược lên; handler gọi `response.Success(result, "", 0, ctx)` để ghi ra JSON envelope chuẩn qua `ctx.JSON`.
7. Nếu ở bất kỳ lớp nào gọi `ctx.Error(err)` thay vì trả kết quả, response sẽ **không** được ghi trực tiếp — middleware **`ErrorHandler`** chạy sau khi `c.Next()` trả về, kiểm tra `c.Errors`, type-assert sang `*response.Exception` để lấy đúng HTTP status, rồi ghi error envelope. Nhờ vậy toàn bộ logic format lỗi chỉ tồn tại ở đúng một nơi, bất kể lỗi phát sinh từ lớp nào.

## 5. Mô hình xử lý lỗi

- `internal/common/response/exception.response.go` định nghĩa một type `Exception` duy nhất implement interface `error`, cùng các hàm khởi tạo tương ứng từng HTTP status (`NewBadRequestException`, `NewUnauthorizedErrorException`, `NewForbiddenErrorException`, `NewNotFoundErrorException`, `NewInternalServerErrorException`).
- Handler/usecase không bao giờ tự ghi response lỗi trực tiếp — chúng gọi `ctx.Error(err)` rồi `return`. Việc ghi response được nhường lại cho middleware `ErrorHandler`, giúp format response đồng nhất trên toàn bộ API.
- Panic không bắt được sẽ được `gin.CustomRecovery` chuyển thành lỗi `500`, đi qua cùng đường `ctx.Error(...)` nên envelope vẫn nhất quán kể cả với lỗi không lường trước.

## 6. Các mối quan tâm xuyên suốt (cross-cutting concerns)

- **CORS**: hiện hardcode chỉ cho phép `http://localhost:3000` (dev server của frontend) — xem [`cmd/app.go`](../cmd/app.go). Chưa có danh sách origin cấu hình theo môi trường.
- **Auth**: đã có việc cấp JWT (`tokenUsecase.CreateAccessToken`, thuật toán HS256, secret/thời hạn lấy từ env `SECRET_ACCESS_TOKEN`/`EXPIRES_AT_ACCESS_TOKEN`) và **giờ đã có middleware verify JWT** (`middlewares.AuthMiddleware.Protect`, mới thêm 2026-08-12). Cơ chế: sau khi login, access token được set vào cookie `accessToken` (không còn trả trong response body); `Protect` đọc cookie này, verify bằng `tokenUsecase.VerifyAccessToken`, trả `401` nếu thiếu/không hợp lệ. Middleware này mới chỉ được gắn cho route `GET /api/auth/get-info` — các route khác (article, demo) vẫn công khai. Việc băm mật khẩu (bcrypt) đã được áp dụng trong `authUsecase.Register`/`Login`.
  - Claim sau khi verify **chưa** được đưa vào `gin.Context` (không có `ctx.Set(...)` trong `Protect`), nên các handler phía sau middleware chưa lấy lại được `UserId` — đây là lý do `authUsecase.GetInfo` hiện vẫn là stub trả `nil, nil`.
  - Cookie `accessToken` hiện set `httpOnly: false, secure: false` — nên siết lại (`httpOnly: true`, `secure: true` khi production) để giảm rủi ro XSS đánh cắp token.
- **Soft delete**: `ent/soft-delete` định nghĩa một mixin được áp cho cả 4 schema (thêm cột kiểu `deleted_at`), nhưng các implementation ở repository (ví dụ `articleRepository.Delete`) hiện đang hard delete qua `DeleteOneID(...).Exec(ctx)` — mixin này chưa được enforce ở tầng query.
- **Dùng song song 2 ORM**: ent là đường chính cho mọi thao tác đọc/xóa đã wire xong và phần lớn thao tác ghi; GORM tồn tại như một nhánh ghi thay thế (`CreateGorm`) và được dùng qua `internal/common/models` — nhìn giống một sự so sánh đang thử nghiệm/mang tính học tập giữa hai ORM hơn là một quy ước đã chốt. Ai mở rộng tầng repository nên xác nhận rõ resource đó dự định dùng ORM nào.
- **Phân trang**: tập trung trong `internal/common/pagination` — `Get(pageString, pageSizeString)` chuẩn hóa giá trị thiếu/không hợp lệ về `page=1, pageSize=3` và tính `Offset`. `PaginationRes` là hình dạng response phân trang dùng chung (`items`, `page`, `pageSize`, `totalItem`, `totalPage`).
- **Cấu hình**: toàn bộ config runtime đi qua một struct `env.Env` có kiểu, load một lần duy nhất lúc khởi động từ `.env` (xem danh sách biến ở [tài liệu API §Cấu hình môi trường](API.md#cấu-hình-môi-trường)).

## 7. Frontend & bối cảnh toàn workspace

Frontend Next.js tại [`src_FE/go-frontend`](../../../src_FE/go-frontend) là một app/repo riêng (có `.git` riêng, `package.json` riêng), giao tiếp với backend này qua HTTP (đó là lý do CORS chỉ cho phép `localhost:3000`). Frontend dùng React Query (`@tanstack/react-query`) theo pattern `api/actions` + `api/tantask` (TanStack) cho từng resource.

Đáng chú ý: frontend đã scaffold sẵn route/module API cho khá nhiều resource — `chat`, `friend`, `role`, `permission`, `totp`, `comment`, `reaction`, `elastic-search`, `google-authenticator` — mà **backend chưa có** (hiện chỉ mới implement `demo`, `article`, và `auth` phía server, theo [API.md](API.md)). Nên coi phạm vi của frontend là mục tiêu/kế hoạch, còn backend là trạng thái đã implement hiện tại; khoảng cách này là bình thường với một dự án học tập đang làm dở, không phải bug.

## 8. Chạy dự án local

```bash
# chạy dev server hot-reload (dùng .air.toml, build ./cmd → ./tmp/main)
air

# hoặc chạy trực tiếp
go run ./cmd

# sinh lại code ent sau khi sửa schema
make ent
```

Cần có MySQL đang chạy, truy cập được qua `DATABASE_URL`, và file `.env` đã điền đủ (xem [API.md — Cấu hình môi trường](API.md#cấu-hình-môi-trường)).

## 9. Các điểm còn thiếu / cần kiểm tra trước khi mở rộng

- **`authUsecase.GetInfo` mới chỉ là stub** (`return nil, nil`) — chưa thực sự đọc `UserId` từ claim và query user tương ứng. *(mới phát hiện 2026-08-12)*
- **`AuthMiddleware.Protect` chưa lưu claim vào `gin.Context`** — cần `ctx.Set("userId", claim.UserId)` (hoặc tương đương) để handler phía sau dùng được, đây là điều kiện tiên quyết để hoàn thiện `GetInfo`. *(mới phát hiện 2026-08-12)*
- **Cookie `accessToken` chưa bật `httpOnly`/`secure`** — nên bật trước khi lên production để giảm rủi ro XSS đánh cắp token. *(mới phát hiện 2026-08-12)*
- Middleware xác thực mới chỉ bảo vệ 1 route (`get-info`) — route `article` (create/delete) vẫn public dù về nghiệp vụ nên yêu cầu đăng nhập.
- Bảng `foods` và `orders` đã được model hóa (schema ent + migration) nhưng chưa có repository/usecase/handler/delivery — chưa expose qua API.
- Mixin soft-delete có mặt ở mọi schema nhưng chưa được enforce nhất quán ở các query/delete.
- Hai ORM (ent + GORM) đang cùng tồn tại cho resource `articles`; các resource mới nên chọn dứt khoát một ORM thay vì copy theo pattern gần nhất.
- `article.repository_impl.go` phần `Create` đang hardcode `user_id = 1` — cần thay bằng ID của user đã xác thực (giờ đã có middleware `Protect` để dựa vào, chỉ còn thiếu bước set claim vào context ở trên).

✅ Đã giải quyết: mục "chưa có middleware verify JWT" (nêu ở lần viết tài liệu trước) — xem [CHANGELOG.md](CHANGELOG.md).
