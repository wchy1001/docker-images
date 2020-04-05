package test

import (
	"github.com/wchy1001/docker-images/images"
	"testing"
)

func TestDockerpull(t *testing.T) {
	var a []string
	a = append(a, "nginx", "docker.io/nginx", "nginx:latest", "docker.io/nginx:latest")
	for _, i := range a {
		image := images.Newimage(i)
		image.Do("127.0.0.1:4000", "wchy1001")
		if image.OldImage != "docker.io/nginx:latest" {
			t.Errorf("images func test failed，input: %s, but output: %s", "nginx", image.OldImage)
		} else if image.NewImage != "127.0.0.1:4000/wchy1001/nginx:latest" {
			t.Errorf("the output image name is wrong，oldinput: %s, but outputimage: %s", image.OldImage, image.NewImage)
		}
	}
}