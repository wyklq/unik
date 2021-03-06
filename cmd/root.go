package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/emc-advanced-dev/unik/pkg/config"
	"github.com/emc-advanced-dev/unik/pkg/types"
)

var clientConfigFile, host string
var port int

var RootCmd = &cobra.Command{
	Use:   "unik",
	Short: "The unikernel compilation, deployment, and management tool",
	Long: `Unik is a tool for compiling application source code
into bootable disk images. Unik also runs and manages unikernel
instances across infrastructures.

Create a client configuration file with 'unik target'

You may set a custom client configuration file w
ith the global flag --client-config=<path>`,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&clientConfigFile, "client-config", os.Getenv("HOME")+"/.unik/client-config.yaml", "client config file (default is $HOME/.unik/client-config.yaml)")
	RootCmd.PersistentFlags().StringVar(&host, "host", "", "<string, optional>: host/ip address of the host running the unik daemon")
	targetCmd.Flags().IntVar(&port, "port", 3000, "<int, optional>: port the daemon is running on (default: 3000)")
}

var clientConfig config.ClientConfig

func readClientConfig() error {
	data, err := ioutil.ReadFile(clientConfigFile)
	if err != nil {
		logrus.WithError(err).Errorf("failed to read client configuration file at " + clientConfigFile + `\n
Try setting your config with 'unik target --host HOST_URL'`)
		return err
	}
	data = bytes.Replace(data, []byte("\n"), []byte{}, -1)
	if err := yaml.Unmarshal(data, &clientConfig); err != nil {
		logrus.WithError(err).Errorf("failed to parse client configuration yaml at " + clientConfigFile + `\n
Please ensure config file contains valid yaml.'\n
Try setting your config with 'unik target --host HOST_URL'`)
		return err
	}
	return nil
}

func printImages(images ...*types.Image) {
	fmt.Printf("%-20s %-20s %-15s %-30s %-6s %-20s\n", "NAME", "ID", "INFRASTRUCTURE", "CREATED", "SIZE(MB)", "MOUNTPOINTS")
	for _, image := range images {
		printImage(image)
	}
}

func printImage(image *types.Image) {
	for i, deviceMapping := range image.RunSpec.DeviceMappings {
		//ignore root device mount point
		if deviceMapping.MountPoint == "/" {
			image.RunSpec.DeviceMappings = append(image.RunSpec.DeviceMappings[:i], image.RunSpec.DeviceMappings[i+1:]...)
		}
	}
	if len(image.RunSpec.DeviceMappings) == 0 {
		fmt.Printf("%-20.20s %-20.20s %-15.15s %-30.30s %-8.0d \n", image.Name, image.Id, image.Infrastructure, image.Created.String(), image.SizeMb)
	} else if len(image.RunSpec.DeviceMappings) > 0 {
		fmt.Printf("%-20.20s %-20.20s %-15.15s %-30.30s %-8.0d %-20.20s\n", image.Name, image.Id, image.Infrastructure, image.Created.String(), image.SizeMb, image.RunSpec.DeviceMappings[0].MountPoint)
		if len(image.RunSpec.DeviceMappings) > 1 {
			for i := 1; i < len(image.RunSpec.DeviceMappings); i++ {
				fmt.Printf("%102s\n", image.RunSpec.DeviceMappings[i].MountPoint)
			}
		}
	}
}

func printInstances(instance ...*types.Instance) {
	fmt.Printf("%-15s %-20s %-14s %-30s %-20s %-15s %-12s\n",
		"NAME", "ID", "INFRASTRUCTURE", "CREATED", "IMAGE", "IPADDRESS", "STATE")
	for _, instance := range instance {
		printInstance(instance)
	}
}

func printInstance(instance *types.Instance) {
	fmt.Printf("%-15.15s %-20.20s %-14.14s %-30.30s %-20.20v %-15.15s %-12.12s\n",
		instance.Name, instance.Id, instance.Infrastructure, instance.Created.String(), instance.ImageId, instance.IpAddress, instance.State)
}

func printVolumes(volume ...*types.Volume) {
	fmt.Printf("%-15.15s %-15.15s %-14.14s %-30.30s %-20.20v %-12.12s\n",
		"NAME", "ID", "INFRASTRUCTURE", "CREATED", "ATTACHED-INSTANCE", "SIZE(MB)")
	for _, volume := range volume {
		printVolume(volume)
	}
}

func printVolume(volume *types.Volume) {
	fmt.Printf("%-15.15s %-15.15s %-14.14s %-30.30s %-20.20v %-12.12d\n",
		volume.Name, volume.Id, volume.Infrastructure, volume.Created.String(), volume.Attachment, volume.SizeMb)
}
