package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.32.192:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.32.192:8021/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check
	client.Agent().ServiceRegister(registration)

	return nil
}

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.32.192:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}

	for key, _ := range data {
		fmt.Println(key)
	}

}

func FilterSerivice() {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.32.192:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter("Service == `user-web`")
	if err != nil {
		panic(err)
	}

	for key, _ := range data {
		fmt.Println(key)
	}
}

func main() {
	_ = Register("192.168.32.192", 8021, "user-web", []string{"shop", "bobby"}, "user-web")
	//AllServices()
	FilterSerivice()
	fmt.Println(fmt.Sprintf(`Service == "%s"`, "user-srv"))
}
