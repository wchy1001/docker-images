package images

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type dockerclient struct {
	cli 	*client.Client
	ctx		context.Context
}

type image struct {
	// docker.io/nginx:latest source image name
	OldImage string
	// 127.0.0.1:4000/org/nginx:latest  new image
	NewImage string
	// Imgname is nginx
	Imgname  string
	// Tag is latest
	Tag		 string
	client	 	dockerclient
}


//生成所有的镜像名字
func (i *image) generateNewImage(dhost, dorg string) error {
	bareimage := strings.Split(string(i.OldImage), "/")
	// eg: nginx:latest
	if len(bareimage) ==1 {
		img :=strings.Split(bareimage[0],":")
		if len(img)>2{
			return errors.New("image format is invalid: "+i.OldImage)
		}
		if len(img) == 1 {
			i.OldImage = "docker.io/" + i.OldImage + ":latest"
			i.Imgname = bareimage[0]
			i.Tag = "latest"
		} else if len(img) == 2 {
			i.OldImage = "docker.io/"+i.OldImage
			i.Imgname = img[0]
			i.Tag = img[1]
		}
	}
	if len(bareimage) >1 {
		last := bareimage[len(bareimage)-1]
		if len(strings.Split(last,":"))==1{
			i.OldImage = i.OldImage+":latest"
			i.Imgname = last
			i.Tag = "latest"
		} else if img :=strings.Split(last,":"); len(img)==2 {
			i.Imgname = img[0]
			i.Tag = img[1]
		} else {
			return errors.New("image format is invalid: "+i.OldImage)
		}

	}
	i.NewImage = dhost+"/"+dorg+"/"+i.Imgname+":"+i.Tag
	return nil
}

func (i *image) GenerateDockerClient(){
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	i.client.cli = cli
	i.client.ctx = ctx
}

func (i *image) pull() error {
	out, err := i.client.cli.ImagePull(i.client.ctx, i.OldImage, types.ImagePullOptions{})
	defer out.Close()
	if err != nil {
		log.Printf("%v image pull failed: %v",i.OldImage,err)
		return err
	}
	io.Copy(os.Stdout, out)
	return nil
}

func (i *image) retag() error{
	images, err := i.client.cli.ImageList(i.client.ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}
	for _, image := range images {
		for _, t := range image.RepoTags{
			if t == i.NewImage{
				command := fmt.Sprintf("docker rmi %v", i.NewImage)
				if err := exec.Command("/bin/bash", "-c", command).Run(); err != nil{
					return err
				}

			}
		}
	}
	if err := i.client.cli.ImageTag(i.client.ctx, i.OldImage,i.NewImage);err != nil{
		return fmt.Errorf("failed to retag images: %v",err.Error())
	}
	return nil
}
func (i *image) push() error {
	log.Printf("push image: %v",i.NewImage)
	//必须要写一个认证，否则报错。认证随便写一个即可，适用于没有认证的请求。
	//_, err := i.client.cli.ImagePush(i.client.ctx, i.NewImage, types.ImagePushOptions{RegistryAuth:"123"})
	// 使用shell来替代这个操作，因为ImagePush 可能是非阻塞的状态，先pass
	command := fmt.Sprintf("docker push %v", i.NewImage)

	if err := exec.Command("/bin/bash", "-c", command).Run(); err != nil{
		return fmt.Errorf("failed to push images: %v",err.Error())
	}
	return nil
}

func (i *image) Do(dhost,dorg string) error{
	if err := i.generateNewImage(dhost, dorg); err !=nil{
		log.Printf("failed to generate new images %v",err)
		return nil
	}
	i.GenerateDockerClient()
	if err := i.pull(); err != nil{
		log.Printf("failed to pull images %v :  %v",i.OldImage, err)
		return err
	}
	if err := i.retag(); err != nil{
		return err
	}
	if err := i.push();err != nil{
		return err
	}
	return nil
}
func Newimage (imagename string) image{
	return image{OldImage: imagename}
}