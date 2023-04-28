blog:
	@ go run -C example/blog -mod=mod main.go serve

blog.watch:
	@ watch --clear -- go run -C example/blog -mod=mod main.go serve