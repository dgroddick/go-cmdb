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

const (
	servicenowInstance = ""
	servicenowUser     = ""
	servicenowPassword = ""
)

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
	updateServerFlag := flag.String("update-server", "", "Update details of a server by sys_id")
	updateNameFlag := flag.String("name", "", "New name for the server")
	updateDescriptionFlag := flag.String("description", "", "New description for the server")
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
	} else if *updateServerFlag != "" {
		if *updateNameFlag == "" && *updateDescriptionFlag == "" {
			log.Fatalf("Please provide a new name or description to update the server")
		}

		updateData := UpdatePayload{
			Name:             *updateNameFlag,
			ShortDescription: *updateDescriptionFlag,
		}

		err := updateCMDBServer(client, *updateServerFlag, updateData)
		if err != nil {
			log.Fatalf("Error updating server: %v", err)
		} else {
			fmt.Println("Server updated successfully.")
		}
	} else {
		flag.Usage()
	}
}

func getCMDBServers(client *resty.Client) ([]Server, error) {
	resp, err := client.R().
		SetBasicAuth(servicenowUser, servicenowPassword).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("https://%s.service-now.com/api/now/table/cmdb_ci_linux_server", servicenowInstance))

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
		SetBasicAuth(servicenowUser, servicenowPassword).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("https://%s.service-now.com/api/now/table/cmdb_ci_linux_server/%s", servicenowInstance, sysID))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status())
	}

	// Debugging: Print response body
	fmt.Printf("Response Body: %s\n", resp.Body())

	var server Server
	err = json.Unmarshal(resp.Body(), &server)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &server, nil
}

func updateCMDBServer(client *resty.Client, sysID string, payload UpdatePayload) error {
	resp, err := client.R().
		SetBasicAuth(servicenowUser, servicenowPassword).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Put(fmt.Sprintf("https://%s.service-now.com/api/now/table/cmdb_ci_linux_server/%s", servicenowInstance, sysID))

	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status())
	}

	return nil
}
