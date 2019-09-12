package kind

import (
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/yaml.v2"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"kind": resourceKind(),
		},
	}
}

func resourceKind() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the 'kind' cluster",
				Required:    true,
			},
			"config": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Configuration file for the 'kind' cluster",
				ForceNew:    true,
				Optional:    true,
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Description: "Base64 encoded public certificate used by clients to authenticate to the cluster endpoint.",
				Computed:    true,
			},
			"client_key": {
				Type:        schema.TypeString,
				Description: "Base64 encoded private key used by clients to authenticate to the cluster endpoint.",
				Computed:    true,
			},
			"cluster_ca_certificate": {
				Type:        schema.TypeString,
				Description: "Base64 encoded public certificate that is the root of trust for the cluster.",
				Computed:    true,
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Endpoint that can be used to reach API server",
				Computed:    true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)

	args := []string{"create", "cluster", "--name", name}

	if config, ok := d.GetOk("config"); ok == true {
		args = append(args, "--config", config.(string))
	}

	if out, err := exec.Command("kind", args...).Output(); err != nil {
		return fmt.Errorf("could not create cluster: %s", out)
	}

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	var out []byte
	var err error

	if out, err = exec.Command("kind", "get", "kubeconfig", "--name", name).Output(); err != nil {
		return fmt.Errorf("could not get cluster config: %s", out)
	}

	var config KubeConfig
	if err = yaml.Unmarshal([]byte(out), &config); err != nil {
		return err
	}

	host := config.Clusters[0].Cluster.Server
	var clusterCaCertificate, clientCertificate, clientKey []byte

	if clusterCaCertificate, err = base64.StdEncoding.DecodeString(config.Clusters[0].Cluster.CertificateAuthorityData); err != nil {
		return err
	}
	if clientCertificate, err = base64.StdEncoding.DecodeString(config.Users[0].User.ClientCertificateData); err != nil {
		return err
	}
	if clientKey, _ = base64.StdEncoding.DecodeString(config.Users[0].User.ClientKeyData); err != nil {
		return err
	}

	d.Set("host", host)
	d.Set("cluster_ca_certificate", string(clusterCaCertificate))
	d.Set("client_certificate", string(clientCertificate))
	d.Set("client_key", string(clientKey))

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	if out, err := exec.Command("kind", "delete", "cluster", "--name", name).Output(); err != nil {
		return fmt.Errorf("could not delete cluster: %s", out)
	}
	d.SetId("")
	return nil
}

type KubeConfig struct {
	APIVersion     string           `yaml:"apiVersion"`
	Clusters       []ClusterElement `yaml:"clusters"`
	Contexts       []ContextElement `yaml:"contexts"`
	CurrentContext string           `yaml:"current-context"`
	Kind           string           `yaml:"kind"`
	Preferences    Preferences      `yaml:"preferences"`
	Users          []UserElement    `yaml:"users"`
}

type ClusterElement struct {
	Cluster ClusterCluster `yaml:"cluster"`
	Name    string         `yaml:"name"`
}

type ClusterCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type ContextElement struct {
	Context ContextContext `yaml:"context"`
	Name    string         `yaml:"name"`
}

type ContextContext struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type Preferences struct {
}

type UserElement struct {
	Name string   `yaml:"name"`
	User UserUser `yaml:"user"`
}

type UserUser struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}
