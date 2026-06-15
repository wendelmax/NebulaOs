// NebulaOS Export Tool — Export all configuration for migration/backup
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var apiBase string

func get(path string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", apiBase+path, nil)
	if err != nil {
		return nil, err
	}
	token := os.Getenv("NEBULA_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	if resp.StatusCode != 200 {
		return data, fmt.Errorf("GET %s: status %d", path, resp.StatusCode)
	}
	return data, nil
}

func getList(path string) ([]interface{}, error) {
	req, err := http.NewRequest("GET", apiBase+path, nil)
	if err != nil {
		return nil, err
	}
	token := os.Getenv("NEBULA_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data []interface{}
	json.Unmarshal(body, &data)
	if resp.StatusCode != 200 {
		return data, fmt.Errorf("GET %s: status %d", path, resp.StatusCode)
	}
	return data, nil
}

func main() {
	apiBase = os.Getenv("NEBULA_API_URL")
	if apiBase == "" {
		apiBase = "http://localhost:8000"
	}

	fmt.Fprintf(os.Stderr, "Exporting NebulaOS from %s\n", apiBase)

	export := map[string]interface{}{
		"exported_at":  time.Now().UTC().Format(time.RFC3339),
		"nebulaos_url": apiBase,
	}

	tenants, _ := getList("/tenants")
	export["tenants"] = tenants

	projects, _ := getList("/projects")
	export["projects"] = projects

	resources, _ := getList("/resources")
	export["resources"] = resources

	volumes, _ := getList("/storage/volumes")
	export["volumes"] = volumes

	buckets, _ := getList("/storage/buckets")
	export["buckets"] = buckets

	providers, _ := getList("/api/providers")
	export["providers"] = providers

	organizations, _ := getList("/api/organizations")
	export["organizations"] = organizations

	departments, _ := getList("/api/departments")
	export["departments"] = departments

	nodes, _ := getList("/api/baremetal/nodes")
	export["bare_metal_nodes"] = nodes

	regions, _ := getList("/cloud/regions")
	export["regions"] = regions

	zones, _ := getList("/cloud/zones")
	export["zones"] = zones

	blueprints, _ := getList("/marketplace/blueprints")
	export["blueprints"] = blueprints

	policy, _ := get("/governance/policy")
	export["sovereignty_policy"] = policy

	securityGroups, _ := getList("/security-groups")
	export["security_groups"] = securityGroups

	out, _ := json.MarshalIndent(export, "", "  ")
	fmt.Println(string(out))
}
