package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/melbahja/goph"
	"gopkg.in/yaml.v3"
	"math/big"
	"os"
	"strconv"
)

type config struct {
	Port               int    `yaml:"port"`
	FingerTableEnabled bool   `yaml:"finger_table_enabled"`
	MaxNumOfNodes      int64  `yaml:"max_num_of_nodes" default:"16"`
	RequestTimeout     int64  `yaml:"request_timeout_sec" default:"5"`
	TTLDuration        int64  `yaml:"ttl_min" default:"180"`
	NeighbourCheck     bool   `yaml:"neighbour_check" default:"false"`
	Predecessor        string `yaml:"predecessor"`
	PredecessorPort    string `yaml:"predecessor_port"`
	Successor          string `yaml:"successor"`
	SuccessorPort      string `yaml:"successor_port"`
}

func main() {
	hosts := []string{
		`compute-3-28`,
		`compute-8-7`,
		`compute-6-17`,
		`compute-6-22`,
	}

	if len(os.Args) != 1 {
		panic(`missing arguments (eg: ./initiator <max nodes> <num of nodes> )`)
	}

	maxNodes, err := strconv.Atoi(os.Args[0])
	if err != nil {
		panic(`incorrect type for arguments (int required)`)
	}

	numOfNodes, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(`incorrect type for arguments (int required)`)
	}

	nodes := make(map[int]*goph.Client)
	var startedClients []*goph.Client
	for _, h := range hosts {
		bucket, err := bucketId(h, maxNodes)
		if err != nil {
			panic(fmt.Sprintf(`getting bucket id failed for node %s`, h))
		}

		_, exists := nodes[bucket]
		if !exists {
			nodes[bucket] = &goph.Client{}
		}
	}


}

func startNode(host, path, localBinaryPath, binaryName string) (*goph.Client, error) {
	client, err := goph.New(`ywi006`, host, nil)
	if err != nil {
		//panic(fmt.Sprintf(`failed to start node %s`, host))
		return nil, err
	}
	defer client.Close()

	sftp, err := client.NewSftp()
	if err != nil {
		return nil, err
	}

	if err = sftp.Mkdir(path); err != nil {
		// Do not consider it an error if the directory existed
		remoteFi, fiErr := sftp.Lstat(path)
		if fiErr != nil || !remoteFi.IsDir() {
			return nil, err
		}
	}

	conf := config{
		Port:               52520,
		FingerTableEnabled: false,
		MaxNumOfNodes:      4,
		RequestTimeout:     2,
		TTLDuration:        180,
		NeighbourCheck:     false,
		Predecessor:        "compute=6-18",
		PredecessorPort:    "",
		Successor:          "compute-6-18",
		SuccessorPort:      "",
	}

	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		return nil, err
	}

	file, err := sftp.Create(path + `/configs.yaml`)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Write(yamlData)
	if err != nil {
		return nil, err
	}

	err = client.Upload(localBinaryPath, path)
	if err != nil {
		return nil, err
	}

	_, err = client.Run(path + `/` + binaryName + ` &`)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func bucketId(key string, maxNodes int) (int, error) {
	hexVal := sha256.Sum256([]byte(key))
	n := new(big.Int)
	n.SetString(hex.EncodeToString(hexVal[:]), 16)
	return int(n.Uint64() % uint64(maxNodes)), nil
}
