

all: vet test demo

test:
	go test ./...

vet:
	go vet ./...

demo:
	make -C demo

# use quicktime player to capture screen, then, save as .mov
# and place here as 'demo-in.mov', then run this target
demo-out.gif:
	ffmpeg -i demo-in.mov -pix_fmt rgb24 -r 10 -f gif - \
		| gifsicle --optimize=3 --delay=7 \
		> $@

.PHONY: demo
