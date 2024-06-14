package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

var snowInstance = os.Getenv("SN_INSTANCE")
var snowUser = os.Getenv("SN_USER")
var snowPassword = os.Getenv("SN_PASS")
var ciTable = "cmdb_ci_linux_server"

type Result struct {
	Result []Server `json:"result"`
}

type Server struct {
	SysID            string `json:"sys_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
}

type UpdatePayload struct {
	Name             string `json:"name,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
}

func main() {
	listAllFlag := flag.Bool("list-all", false, "List all CMDB CI servers")
	getServerFlag := flag.String("get-server", "", "Get details of a server by sys_id")
	helpFlag := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	client := resty.New()

	if *listAllFlag {
		servers, err := getCMDBServers(client)
		if err != nil {
			log.Fatalf("Error retrieving servers: %v", err)
		}

		fmt.Println("List of CMDB CI Servers:")
		for _, server := range servers {
			fmt.Printf("SysID: %s, Name: %s, Description: %s\n", server.SysID, server.Name, server.ShortDescription)
		}
	} else if *getServerFlag != "" {
		server, err := getCMDBServer(client, *getServerFlag)
		if err != nil {
			log.Fatalf("Error fetching server: %v", err)
		}

		fmt.Printf("Fetched Server - SysID: %s, Name: %s, Description: %s\n", server.SysID, server.Name, server.ShortDescription)
	} else {
		flag.Usage()
	}
}

func getCMDBServers(client *resty.Client) ([]Server, error) {
	resp, err := client.R().
		SetBasicAuth(snowUser, snowPassword).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("https://%s.service-now.com/api/now/table/%s", snowInstance, ciTable))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status())
	}

	var result Result
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return result.Result, nil
}

func getCMDBServer(client *resty.Client, sysID string) (*Server, error) {
	resp, err := client.R().
		SetBasicAuth(snowUser, snowPassword).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("https://%s.service-now.com/api/now/table/%s/%s", snowInstance, ciTable, sysID))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status())
	}

	// Debug
	fmt.Printf("Response Body: %s\n", resp.Body())

	var server Server
	err = json.Unmarshal(resp.Body(), &server)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &server, nil
}
