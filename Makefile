.SILENT:
run:
	go run main.go

build:
	fyne package -os windows --app-build 1 --release

build-linux:
	fyne package -os linux --app-build 1 --release

build-macos:
	fyne package -os darwin --app-build 1 --release
