# Khoa chính:
- Table User:
edge. To ("Articles", Articles. Type),
# Khoá phu:
- Table Article:
edge. From("Users", Users.Type).Ref("Articles").Field("user_id").Unique().Required(),
