# Tài liệu API — go-backend

Base URL: `http://{HOST}:{PORT}/api` (mặc định khi chạy local: `http://localhost:8080/api`)

Toàn bộ route được gắn dưới group `/api` (xem [`root.delivery.go`](../internal/delivery/root.delivery.go)).

> 📜 Lịch sử cập nhật tài liệu này: xem [CHANGELOG.md](CHANGELOG.md).

## Quy ước chung

### Response thành công

Mọi response thành công đều được bọc bởi [`response.Success`](../internal/common/response/response.go):

```json
{
  "status": "success",
  "statusCode": 200,
  "message": "OK",
  "data": { },
  "doc": ""
}
```

- `statusCode` mặc định là `200`, `message` mặc định là text chuẩn của HTTP status nếu handler không truyền giá trị khác.
- `data` chứa payload riêng của từng endpoint, được mô tả chi tiết bên dưới.

### Response lỗi

Lỗi được đẩy vào context qua `ctx.Error(...)` và được xử lý tập trung bởi [`ErrorHandler`](../internal/common/middlewares/error.middleware.go):

```json
{
  "status": "error",
  "statusCode": 400,
  "message": "Body required",
  "doc": ""
}
```

| Hàm khởi tạo Exception | HTTP status |
|---|---|
| `NewBadRequestException` | 400 |
| `NewUnauthorizedErrorException` | 401 |
| `NewForbiddenErrorException` | 403 |
| `NewNotFoundErrorException` | 404 |
| `NewInternalServerErrorException` | 500 |

Mọi panic không được xử lý sẽ được `gin.CustomRecovery` bắt lại và trả về `500 Internal Server Error`.

### CORS

Chỉ cho phép origin `http://localhost:3000` (frontend Next.js) gọi API (`cmd/app.go`).

---

## Auth API — `/api/auth`

### POST `/api/auth/register`

Tạo tài khoản người dùng mới. Mật khẩu được băm bằng bcrypt trước khi lưu (`authUsecase.Register`).

**Request body**

```json
{
  "email": "user@example.com",
  "fullName": "Nguyen Van A",
  "password": "plaintext-password"
}
```

| Field | Kiểu | Ràng buộc |
|---|---|---|
| `email` | string | bắt buộc, phải đúng định dạng email |
| `fullName` | string | bắt buộc |
| `password` | string | bắt buộc |

**Response**

- `200` — trả về (trong `data`) entity user vừa tạo, không có field `password` (do được đánh dấu `Sensitive()` trong ent schema).
- `400` — body thiếu/không hợp lệ, hoặc email đã tồn tại (`"Email đã có vui lòng đăng nhập"`).

### POST `/api/auth/login`

Xác thực người dùng và set JWT access token vào cookie.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "plaintext-password"
}
```

| Field | Kiểu | Ràng buộc |
|---|---|---|
| `email` | string | bắt buộc, phải đúng định dạng email |
| `password` | string | bắt buộc |

**Response**

- `200` — `data: true`. Access token **không** trả trong body nữa mà được set vào cookie `accessToken` (`Set-Cookie`), `SameSite=Lax`, `MaxAge` = `EXPIRES_AT_ACCESS_TOKEN`, path `/`. *(Cập nhật 2026-08-12 — trước đây token được trả trực tiếp trong `data`, xem [CHANGELOG.md](CHANGELOG.md).)*
- `400`:
  - không tìm thấy user theo email
  - user chưa có password (tài khoản đăng nhập bằng Google) — `"Vui lòng đăng nhập bằng google để cập nhật"`
  - sai mật khẩu — `"Mật khẩu không chính xác"`

Claims JWT gồm: `UserId` + `ExpiresAt` chuẩn (thời hạn cấu hình qua env `EXPIRES_AT_ACCESS_TOKEN`), ký bằng HS256 với secret `SECRET_ACCESS_TOKEN`.

> ⚠️ Cookie hiện được set với `secure: false` và `httpOnly: false` (xem [`auth.handler.go`](../internal/handler/auth.handler.go)) — token có thể bị đọc bởi JavaScript phía client và không bắt buộc HTTPS. Nên bật `httpOnly: true` (chống XSS đọc cookie) và `secure: true` khi lên production.

### GET `/api/auth/get-info` 🔒

**Mới thêm.** Trả về thông tin user hiện tại dựa trên access token trong cookie. Route này được bảo vệ bởi `AuthMiddleware.Protect`.

**Yêu cầu**: cookie `accessToken` hợp lệ (được set từ bước login).

**Response**

- `200` — hiện tại `data: null`. **`authUsecase.GetInfo` mới chỉ là stub (`return nil, nil`), chưa thực sự truy vấn user theo `UserId` trong claim** — cần hoàn thiện trước khi FE dùng được endpoint này.
- `401` — thiếu cookie `accessToken`, hoặc token không hợp lệ/hết hạn.
- `400` — lỗi đọc cookie khác (không phải do thiếu cookie).

### Middleware xác thực — `AuthMiddleware.Protect`

**Mới thêm.** Áp dụng cho mọi route được `.Use(d.authMiddleware.Protect)` gắn vào (hiện chỉ có `GET /api/auth/get-info`).

Luồng xử lý ([`auth.middleware.go`](../internal/common/middlewares/auth.middleware.go)):
1. Đọc cookie `accessToken` từ request.
2. Không có cookie → `401 Unauthorized`.
3. Lỗi khác khi đọc cookie → `400 Bad Request`.
4. Có cookie → verify bằng `tokenUsecase.VerifyAccessToken` (parse + kiểm tra chữ ký + hạn dùng); verify thất bại → `401 Unauthorized`.
5. Verify thành công → cho request đi tiếp (**lưu ý**: claim hiện *chưa* được `ctx.Set(...)` vào context, nên handler phía sau không lấy lại được `UserId` từ middleware — đây là lý do `GetInfo` chưa implement được).

> Route `article` và các route khác **chưa** được gắn middleware này — vẫn truy cập công khai, không cần đăng nhập.

---

## Article API — `/api/article`

Thao tác trên bảng `articles` (xem [Data model](#data-model)).

### POST `/api/article`

Tạo bài viết mới.

**Request body**

```json
{
  "title": "My article",
  "content": "Optional content",
  "image_url": "https://example.com/image.png"
}
```

| Field | Kiểu | Ràng buộc |
|---|---|---|
| `title` | string | bắt buộc |
| `content` | string \| null | optional |
| `image_url` | string \| null | optional |

**Response**

- `200` — `data` là entity article vừa tạo.
- `400` — body thiếu/không hợp lệ.

> Lưu ý implementation: `user_id` hiện đang bị hardcode = `1` khi tạo (xem [`article.repository_impl.go`](../internal/repository/repository_impl/article.repository_impl.go)) — do chưa có context user đã xác thực được truyền vào.

### GET `/api/article`

Lấy danh sách bài viết có phân trang và filter.

**Query params**

| Param | Kiểu | Mặc định | Mô tả |
|---|---|---|---|
| `page` | int | `1` | Số trang (bắt đầu từ 1); giá trị không hợp lệ/≤0 sẽ dùng mặc định |
| `pageSize` | int | `3` | Số item mỗi trang; giá trị không hợp lệ/≤0 sẽ dùng mặc định |
| `filters` | chuỗi JSON | `{}` | Xem bên dưới |

`filters` là một chuỗi JSON được encode trong URL, khớp với:

```json
{
  "id": 0,
  "content": "",
  "views": null
}
```

- `id` — match chính xác, bỏ qua nếu `<= 0`
- `content` — match theo substring, không phân biệt hoa/thường trên field `content`, bỏ qua nếu rỗng
- `views` — match chính xác, bỏ qua nếu `null`

Ví dụ: `GET /api/article?page=1&pageSize=10&filters={"content":"golang"}`

**Response**

```json
{
  "status": "success",
  "statusCode": 200,
  "message": "OK",
  "data": {
    "items": [ /* mảng article entity, mỗi item có kèm edge Users */ ],
    "page": 1,
    "pageSize": 10,
    "totalItem": 42,
    "totalPage": 5
  },
  "doc": ""
}
```

### DELETE `/api/article/:id`

Xóa bài viết theo ID (hard delete qua ent, dù schema có mixin soft-delete).

**Path params**

| Param | Kiểu | Mô tả |
|---|---|---|
| `id` | int | ID bài viết |

**Response**

- `200` — `data: true`
- `400` — `id` không phải số hợp lệ, hoặc xóa thất bại (ví dụ không tìm thấy).

---

## Demo API — `/api/demo`

Các endpoint mẫu/tham khảo dùng để minh họa cách bind request trong Gin. Không thuộc phạm vi tính năng sản phẩm.

| Method | Path | Mục đích |
|---|---|---|
| `GET` | `/api/demo/query` | Đọc query param `page`/`pageSize`, trả lại qua `demoUsecase.Query`. |
| `DELETE` | `/api/demo/param/:id` | Đọc path param `:id`; minh họa cách truyền dữ liệu giữa các middleware qua `gin.Context`. |
| `POST` | `/api/demo/body` | Bind JSON body theo `{ "email": string (required, email), "password": string (required) }`. |
| `PATCH` | `/api/demo/header` | Đọc header `API_KEY` từ request. |

Cả 4 endpoint đều trả kết quả theo envelope chuẩn `response.Success`.

---

## Data model

Các entity được định nghĩa bằng [ent](https://entgo.io) trong [`ent/schema`](../ent/schema) và ánh xạ vào MySQL qua GORM. Tất cả các bảng đều có mixin soft-delete (`deleted_at`), tuy nhiên không phải mọi hàm ở repository đều đang tôn trọng field này.

### `users`

| Field | Kiểu | Ghi chú |
|---|---|---|
| `id` | int | PK |
| `email` | string | unique |
| `full_name` | string, nullable | |
| `avatar` | string, nullable | |
| `password` | string, nullable, sensitive | hash bcrypt; không bao giờ serialize ra JSON |
| `totp_secret` | string, nullable, sensitive | |
| `google_id` | string, nullable, sensitive | |
| `created_at` / `updated_at` | time | |

Edges: có nhiều `articles`, có nhiều `orders`.

### `articles`

| Field | Kiểu | Ghi chú |
|---|---|---|
| `id` | int | PK |
| `title` | string | tối đa 255 ký tự |
| `content` | string, nullable | |
| `image_url` | string, nullable | JSON key là `imageUrl` |
| `like_count` | int | mặc định `0`, JSON key là `likeCount` |
| `views` | int | mặc định `0` |
| `user_id` | int | FK → `users.id`, bắt buộc |
| `created_at` / `updated_at` | time | |

### `foods`

| Field | Kiểu | Ghi chú |
|---|---|---|
| `id` | int | PK |
| `name` | string | |
| `description` | string, optional | |
| `created_at` / `updated_at` | time | |

Edges: có nhiều `orders`.

### `orders`

| Field | Kiểu | Ghi chú |
|---|---|---|
| `id` | int | PK |
| `user_id` | int | FK → `users.id`, bắt buộc |
| `food_id` | int | FK → `foods.id`, bắt buộc |
| `created_at` / `updated_at` | time | |

> Bảng `foods` và `orders` đã có trong schema/migration nhưng hiện **chưa có** delivery/handler/usecase/repository nào tương ứng — chưa có API nào expose các bảng này.

---

## Cấu hình môi trường

Được đọc từ file `.env` lúc khởi động ([`internal/common/env/env.go`](../internal/common/env/env.go)):

| Biến | Mục đích |
|---|---|
| `IS_PRODUCTION` | `"true"` để đánh dấu môi trường production |
| `HOST` | Host mà Gin server bind vào |
| `PORT` | Port mà Gin server bind vào |
| `DATABASE_URL` | DSN kết nối MySQL, dùng chung cho cả ent client và GORM client |
| `EXPIRES_AT_ACCESS_TOKEN` | Thời hạn JWT access token (chuỗi duration của Go, ví dụ `1h`) |
| `SECRET_ACCESS_TOKEN` | Secret dùng để ký JWT (HMAC) |

> 🐛 **Đã sửa (2026-08-12)**: trước đây `env.New()` có đọc `SECRET_ACCESS_TOKEN` nhưng quên gán vào struct `Env` trả về, khiến `SecretAccessToken` luôn rỗng — JWT bị ký/verify bằng secret rỗng. Đã fix trong lần cập nhật này.
