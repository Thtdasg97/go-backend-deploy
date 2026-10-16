.PHONY: ent ent-new

# Chạy generate lại code ent sau khi đổi schema
ent: 
	go generate ./ent


ent-new: 
	go run -mod=mod entgo.io/ent/cmd/ent new ${name}
