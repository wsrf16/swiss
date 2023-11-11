package sdk

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"time"
)

func Login(cli *client.Client, host, username, password string) {
	log.Println("Logging in docker registry...")

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	ok, err := cli.RegistryLogin(ctx, types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: host,
	})
	if err != nil {
		log.Fatalf("Error while logging in docker registry! %s", err.Error())
	}
	log.Printf("%s --- Token: %s\n", ok.Status, ok.IdentityToken)
}

func NewClientWithOpts(ops ...client.Opt) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(ops...)
	if err != nil {
		panic(err)
	}
	return cli, err
}

func ImagePull(cli *client.Client, image string) (io.ReadCloser, error) {
	reader, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	return reader, err
}

func ImageList(cli *client.Client) ([]types.ImageSummary, error) {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	return images, err
}

func ContainerPS(cli *client.Client) ([]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	return containers, err
}

func ContainerCreate(cli *client.Client, config *container.Config, hostConfig *container.HostConfig, containerName string) (container.CreateResponse, error) {
	reader, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, containerName)
	return reader, err
}

func ImageSave(cli *client.Client, images []string, output string) error {
	resp, err := cli.ImageSave(context.Background(), images)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp)
	if err != nil {
		return err
	}

	return nil
}

// func PS(cli *client.Client) ([]types.Container, error) {
//     containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
//     return containers, err
// }
