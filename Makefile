helloapp:
	@ # rm -rf example/helloapp
	@ # go run main.go new example/helloapp
	@ # go run main.go -C example/helloapp serve
	@ go run -C example/helloapp -mod=mod main.go

helloapp.watch:
	@ watch --clear -- $(MAKE) helloapp

blog:
	@ go run -C example/blog -mod=mod main.go serve

blog.watch:
	@ watch --clear -- go run -C example/blog -mod=mod main.go serve
