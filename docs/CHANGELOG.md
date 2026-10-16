# Changelog tài liệu — go-backend

Ghi lại lịch sử thay đổi của [`API.md`](API.md) và [`ARCHITECTURE.md`](ARCHITECTURE.md), theo sát các feature/fix thực tế trong code. Mục nào không rõ ngày commit (do chưa commit vào git tại thời điểm viết) thì lấy theo ngày cập nhật tài liệu.

## 2026-08-12 — Thêm auth middleware + protected route

**Feature/fix trong code:**
- Thêm [`middlewares.AuthMiddleware.Protect`](../internal/common/middlewares/auth.middleware.go): đọc cookie `accessToken`, verify JWT qua `tokenUsecase.VerifyAccessToken`, trả `401` nếu thiếu/không hợp lệ.
- Thêm endpoint **`GET /api/auth/get-info`**, được bảo vệ bởi middleware trên (route đầu tiên trong dự án yêu cầu đăng nhập).
- `POST /api/auth/login` đổi cách trả token: **từ trả JWT trong `data` → set vào cookie `accessToken`** (`SetSameSite(Lax)`, `MaxAge` theo `EXPIRES_AT_ACCESS_TOKEN`); response body giờ chỉ còn `data: true`.
- Thêm `usecase.TokenUsecase.VerifyAccessToken` và `dto.CustomClaim` (tách claim ra file `token.dto.go` riêng thay vì định nghĩa inline trong `token.usecase_impl.go`).
- Thêm `usecase.AuthUsecase.GetInfo` — **hiện là stub, trả `nil, nil`, chưa implement logic thật**.
- **Bug fix**: `env.New()` trước đây đọc `SECRET_ACCESS_TOKEN` từ `.env` nhưng quên gán vào struct `Env` trả về (thiếu dòng `SecretAccessToken: secretAccessToken`) → JWT trước đó luôn được ký/verify bằng secret rỗng. Đã bổ sung.
- `di.go`: `Injection(...)` nhận thêm `env *env.Env`; `AuthHandler` và `authDelivery` nhận thêm dependency tương ứng (`env`, `authMiddleware`) để hỗ trợ set cookie và bảo vệ route.

**Cập nhật tài liệu:**
- `API.md`: viết lại mục `POST /api/auth/login` (cookie thay vì body), thêm mục mới `GET /api/auth/get-info` và mô tả middleware `Protect`; thêm cảnh báo cookie chưa bật `httpOnly`/`secure`; ghi chú bug-fix `SecretAccessToken` trong bảng biến môi trường.
- `ARCHITECTURE.md`: cập nhật §6 (Auth) — xác nhận middleware verify JWT đã tồn tại, mô tả cơ chế cookie, nêu rõ claim chưa được set vào `gin.Context`; cập nhật §9 (Known gaps) — đánh dấu gap "chưa có middleware auth" đã giải quyết, thêm 3 gap mới (`GetInfo` stub, claim chưa set vào context, cookie chưa bật `httpOnly`/`secure`).

**Chưa thay đổi / vẫn còn nguyên (so với lần viết tài liệu trước):**
- `article` API vẫn public, không yêu cầu đăng nhập.
- `user_id` khi tạo article vẫn hardcode = `1`.
- `foods`/`orders` vẫn chưa có API.
- Soft-delete mixin vẫn chưa được enforce ở query/delete.

## 2026-08-12 — Khởi tạo tài liệu

- Tạo mới [`API.md`](API.md) và [`ARCHITECTURE.md`](ARCHITECTURE.md), phản ánh trạng thái backend tại thời điểm đó: 3 resource được implement (`demo`, `article`, `auth` — chỉ `register`/`login`, chưa có route bảo vệ), JWT được cấp nhưng chưa được verify ở bất kỳ đâu, `foods`/`orders` chỉ có schema chưa có API.
