/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package views

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/ernestio/api-gateway/models"
	"github.com/r3labs/graph"
)

// BuildRender : Build representation to be rendered on the frontend
type BuildRender struct {
	ID              string              `json:"id"`
	EnvironmentID   int                 `json:"environment_id"`
	Name            string              `json:"name"`
	Project         string              `json:"project"`
	Provider        string              `json:"provider"`
	Status          string              `json:"status"`
	UserID          int                 `json:"user_id"`
	UserName        string              `json:"user_name"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
	Vpcs            []map[string]string `json:"vpcs,omitempty"`
	Networks        []map[string]string `json:"networks,omitempty"`
	Instances       []map[string]string `json:"instances,omitempty"`
	Nats            []map[string]string `json:"nats,omitempty"`
	SecurityGroups  []map[string]string `json:"security_groups,omitempty"`
	Elbs            []map[string]string `json:"elbs,omitempty"`
	RDSClusters     []map[string]string `json:"rds_clusters,omitempty"`
	RDSInstances    []map[string]string `json:"rds_instances,omitempty"`
	EBSVolumes      []map[string]string `json:"ebs_volumes,omitempty"`
	LoadBalancers   []map[string]string `json:"load_balancers,omitempty"`
	SQLDatabases    []map[string]string `json:"sql_databases,omitempty"`
	VirtualMachines []map[string]string `json:"virtual_machines,omitempty"`
}

// Render : Map a Build to a BuildRender
func (o *BuildRender) Render(b models.Build) (err error) {
	o.ID = b.ID
	o.EnvironmentID = b.EnvironmentID
	o.CreatedAt = b.CreatedAt.String()
	o.UpdatedAt = b.UpdatedAt.String()
	o.Status = b.Status
	o.UserID = b.UserID
	o.UserName = b.Username

	g, err := b.GetMapping()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	o.Vpcs = RenderVpcs(g)
	o.Networks = RenderNetworks(g)
	o.SecurityGroups = RenderSecurityGroups(g)
	o.Nats = RenderNats(g)
	o.Instances = RenderInstances(g)
	o.Elbs = RenderELBs(g)
	o.RDSClusters = RenderRDSClusters(g)
	o.RDSInstances = RenderRDSInstances(g)
	o.EBSVolumes = RenderEBSVolumes(g)
	o.LoadBalancers = RenderLoadBalancers(g)
	o.SQLDatabases = RenderSQLDatabases(g)
	o.VirtualMachines = RenderVirtualMachines(g)

	return err
}

// RenderVpcs : renders a builds vpcs
func RenderVpcs(g *graph.Graph) []map[string]string {
	var vpcs []map[string]string

	for _, n := range g.GetComponents().ByType("vpc") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["vpc_aws_id"].(string)
		subnet, _ := (*gc)["subnet"].(string)
		vpcs = append(vpcs, map[string]string{
			"name":       name,
			"vpc_id":     id,
			"vpc_subnet": subnet,
		})
	}

	return vpcs
}

// RenderNetworks : renders a builds networks
func RenderNetworks(g *graph.Graph) []map[string]string {
	var networks []map[string]string

	for _, n := range g.GetComponents().ByType("network") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["network_aws_id"].(string)
		az, _ := (*gc)["availability_zone"].(string)
		networks = append(networks, map[string]string{
			"name":              name,
			"network_aws_id":    id,
			"availability_zone": az,
		})
	}

	return networks
}

// RenderSecurityGroups : renders a builds security groups
func RenderSecurityGroups(g *graph.Graph) []map[string]string {
	var sgs []map[string]string

	for _, n := range g.GetComponents().ByType("firewall") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["security_group_aws_id"].(string)
		sgs = append(sgs, map[string]string{
			"name":                  name,
			"security_group_aws_id": id,
		})
	}

	return sgs
}

// RenderNats : renders a builds nat gateways
func RenderNats(g *graph.Graph) []map[string]string {
	var nats []map[string]string

	for _, n := range g.GetComponents().ByType("nat") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["nat_gateway_aws_id"].(string)
		pubIP, _ := (*gc)["nat_gateway_allocation_ip"].(string)
		nats = append(nats, map[string]string{
			"name":               name,
			"nat_gateway_aws_id": id,
			"public_ip":          pubIP,
		})
	}

	return nats
}

// RenderELBs : renders a builds elbs
func RenderELBs(g *graph.Graph) []map[string]string {
	var elbs []map[string]string

	for _, n := range g.GetComponents().ByType("elb") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		dns, _ := (*gc)["dns_name"].(string)
		elbs = append(elbs, map[string]string{
			"name":     name,
			"dns_name": dns,
		})
	}

	return elbs
}

// RenderInstances : renders a builds instances
func RenderInstances(g *graph.Graph) []map[string]string {
	var instances []map[string]string

	for _, n := range g.GetComponents().ByType("instance") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["instance_aws_id"].(string)
		pip, _ := (*gc)["public_ip"].(string)
		ip, _ := (*gc)["ip"].(string)
		instances = append(instances, map[string]string{
			"name":            name,
			"instance_aws_id": id,
			"public_ip":       pip,
			"ip":              ip,
		})
	}

	return instances
}

// RenderRDSClusters : renders a builds rds clusters
func RenderRDSClusters(g *graph.Graph) []map[string]string {
	var rdss []map[string]string

	for _, n := range g.GetComponents().ByType("rds_cluster") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		endpoint, _ := (*gc)["endpoint"].(string)
		rdss = append(rdss, map[string]string{
			"name":     name,
			"endpoint": endpoint,
		})
	}

	return rdss
}

// RenderRDSInstances : renders a builds rds instances
func RenderRDSInstances(g *graph.Graph) []map[string]string {
	var rdss []map[string]string

	for _, n := range g.GetComponents().ByType("rds_instance") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		endpoint, _ := (*gc)["endpoint"].(string)
		rdss = append(rdss, map[string]string{
			"name":     name,
			"endpoint": endpoint,
		})
	}

	return rdss
}

// RenderEBSVolumes : renders a builds ebs volumes
func RenderEBSVolumes(g *graph.Graph) []map[string]string {
	var rdss []map[string]string

	for _, n := range g.GetComponents().ByType("ebs_volume") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["volume_aws_id"].(string)
		rdss = append(rdss, map[string]string{
			"name":          name,
			"volume_aws_id": id,
		})
	}

	return rdss
}

// RenderLoadBalancers : renders load balancers
func RenderLoadBalancers(g *graph.Graph) []map[string]string {
	var lbs []map[string]string
	ips := listIPAddresses(g)

	for _, n := range g.GetComponents().ByType("lb") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["id"].(string)
		configs, _ := (*gc)["frontend_ip_configurations"].([]interface{})
		ip := ""
		if len(configs) > 0 {
			cfg, _ := configs[0].(map[string]interface{})
			ipID, _ := cfg["public_ip_address_id"].(string)
			ip, _ = ips[ipID]
		}

		lbs = append(lbs, map[string]string{
			"name":      name,
			"id":        id,
			"public_ip": ip,
		})
	}

	return lbs
}

func listIPAddresses(g *graph.Graph) map[string]string {
	existingIPs := make(map[string]string, 0)

	for _, ip := range g.GetComponents().ByType("public_ip") {
		gc := ip.(*graph.GenericComponent)
		id, _ := (*gc)["id"].(string)
		ipAddress, _ := (*gc)["ip_address"].(string)
		existingIPs[id] = ipAddress
	}

	return existingIPs
}

// RenderVirtualMachines : renders virtual machines
func RenderVirtualMachines(g *graph.Graph) []map[string]string {
	var resources []map[string]string
	mappedIPs := make(map[string]interface{}, 0)
	existingIPs := listIPAddresses(g)

	for _, ni := range g.GetComponents().ByType("network_interface") {
		var public []string
		var private []string

		gc := ni.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		ips := make(map[string][]string)

		configs, _ := (*gc)["ip_configuration"].([]interface{})
		for _, cfg := range configs {
			c, _ := cfg.(map[string]interface{})
			pubID, _ := c["public_ip_address_id"].(string)
			pri, _ := c["private_ip_address"].(string)
			if pub, ok := existingIPs[pubID]; ok {
				public = append(public, pub)
			}
			private = append(private, pri)
		}

		ips["public"] = public
		ips["private"] = private
		mappedIPs[name] = make(map[string][]string, 0)
		mappedIPs[name] = ips
	}

	for _, n := range g.GetComponents().ByType("virtual_machine") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["id"].(string)
		networks, _ := (*gc)["network_interfaces"].([]interface{})
		publicIPs := make([]string, 0)
		privateIPs := make([]string, 0)
		for _, ni := range networks {
			netName := ni.(string)
			if val, ok := mappedIPs[netName]; ok {
				ips, _ := val.(map[string][]string)
				publicIPs = append(publicIPs, ips["public"]...)
				privateIPs = append(privateIPs, ips["private"]...)
			}
		}

		resources = append(resources, map[string]string{
			"name":       name,
			"id":         id,
			"public_ip":  strings.Join(publicIPs, ", "),
			"private_ip": strings.Join(privateIPs, ", "),
		})
	}

	return resources
}

// RenderSQLDatabases : renders sql databases
func RenderSQLDatabases(g *graph.Graph) []map[string]string {
	return renderResources(g, "sql_database", func(gc *graph.GenericComponent) map[string]string {
		name, _ := (*gc)["name"].(string)
		server, _ := (*gc)["server_name"].(string)
		id, _ := (*gc)["id"].(string)

		return map[string]string{
			"name":        name,
			"server_name": server + ".database.windows.net",
			"id":          id,
		}
	})
}

type convert func(*graph.GenericComponent) map[string]string

func renderResources(g *graph.Graph, resourceType string, f convert) (resources []map[string]string) {
	for _, n := range g.GetComponents().ByType(resourceType) {
		gc := n.(*graph.GenericComponent)
		resources = append(resources, f(gc))
	}

	return
}

// RenderCollection : Maps a collection of Builds on a collection of BuildRender
func (o *BuildRender) RenderCollection(builds []models.Build) (list []BuildRender, err error) {
	for _, b := range builds {
		var output BuildRender
		if err := output.Render(b); err == nil {
			list = append(list, output)
		}
	}

	return list, nil
}

// ToJSON : Converts a ServiceRender to json string
func (o *BuildRender) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}

// RenderChanges : renders build definition steps
func RenderChanges(mapping map[string]interface{}) (result []byte, err error) {
	var lines []string
	var actions = map[string]string{"create": "Create", "update": "Update", "delete": "Delete"}

	if mapping["changes"] != nil {
		for _, change := range mapping["changes"].([]interface{}) {
			component := change.(map[string]interface{})
			c := component["_component"].(string)
			c = strings.Replace(c, "_", " ", -1)
			n := component["name"].(string)
			a := component["_action"].(string)
			line := actions[a] + " a " + c + " named " + n
			lines = append(lines, line)
		}
	}

	result, err = json.Marshal(lines)
	if err != nil {
		return result, err
	}

	return result, nil
}
