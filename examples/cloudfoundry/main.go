package main

import (
	"fmt"
	"github.com/Piszmog/cloudconfigclient/v2"
	"log"
	"strings"
)

func main() {
	// ensure you have set the environment variable 'VCAP_SERVICES' that contains the information to connect to the
	// Config Server

	client, err := cloudconfigclient.New(cloudconfigclient.DefaultCFService())
	if err != nil {
		log.Fatalln(err)
	}

	// load a config file
	configuration, err := client.GetConfiguration("test-app", "cloud")
	if err != nil {
		log.Fatalln(err)
	}
	// we can also call configuration.Unmarshal(..) to unmarshal the configuration into a struct
	localProp, err := configuration.GetPropertySource("application-cloud.yml")
	if err != nil {
		log.Fatalln("failed to find cloud property file")
	}
	fmt.Printf("cloud property file: %+v\n", localProp)
	// handle all config properties
	configuration.HandlePropertySources(func(propertySource cloudconfigclient.PropertySource) {
		if strings.HasSuffix(propertySource.Name, "test-app.properties") {
			// TODO save off values
		}
	})

	// load a specific file (e.g. json/txt)
	var f map[string]string
	// if 'fooDir' has been added to 'searchPaths' in SCS v3.x, then pass "" (blank) for directory
	if err = client.GetFile("fooDir", "bar.json", &f); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("file from default branch: %+v\n", f)

	// load a specific file (e.g. json/txt)
	var b map[string]string
	// if 'fooDir' has been added to 'searchPaths' in SCS v3.x, then pass "" (blank) for directory
	if err = client.GetFileFromBranch("develop", "fooDir", "bar.json", &b); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("file from specific branch: %+v\n", b)
}
