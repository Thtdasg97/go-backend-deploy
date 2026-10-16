# So sánh folder schema và folder migration 

make ent 
đảm bảo file schema ổn

atlas migrate diff create_chatgroup_chatmessage_chatgroupmember \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://mysql/8/ent"


## Tạo lại atlas.sum: khi muốn tự thay đối file sql migration
atlas migrate hash  \
  --dir "file://ent/migrate/migrations" \


## Áp dụng migration vào db
make ent

atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "mysql://root:12345@localhost:3307/go-backend"

Sau khi chạy xong thì sẽ migrate db lên prd, trên prd se có 1 table lưu version

# DB đã có sẵn
- tạo init
- bỏ qua không chạy init

atlas migrate diff init \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://mysql/8/ent"

atlas migrate set 20260801041350 \
--dir "file://ent/migrate/migrations" \
--url "mysql://root:12345@localhost:3307/go-back-prd"

## Flow thay đổi 
1. Modify column
2. Diff - tên nội dung
3. Apply 


atlas migrate diff del_food_price \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://mysql/8/ent"
