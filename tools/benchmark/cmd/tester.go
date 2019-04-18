package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var batchTestCmd = &cobra.Command{
	Use:   "batch-test",
	Short: "Benchmark batch-test",

	Run: batchTestFunc,
}

var (
	inputFile string
)

type testConfig struct {
	TestType     string `yaml:"type"`
	Conns        uint   `yaml:"conns"`
	Clients      uint   `yaml:"clients"`
	Total        int    `yaml:"total"`
	Rate         int    `yaml:"rate"`
	KeySpaceSize int    `yaml:"keySpaceSize,omitempty"`
	Key          string `yaml:"key,omitempty"`
	KeySize      int    `yaml:"keySize,omitempty"`
	ValSize      int    `yaml:"valSize,omitempty"`
}

func init() {
	RootCmd.AddCommand(batchTestCmd)
	batchTestCmd.Flags().StringVar(&inputFile, "input-file", "", "Input file for batch test")
}

func batchTestFunc(cmd *cobra.Command, args []string) {
	testConfigData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading test config file: %s. Err: %+v", inputFile, err)
	}

	var testConfigs []testConfig
	err = yaml.Unmarshal([]byte(testConfigData), &testConfigs)
	if err != nil {
		log.Fatalf("Error unmarshalling test configs: %+v", err)
	}

	ep1 := os.Getenv("HOST_1")
	ep2 := os.Getenv("HOST_2")
	ep3 := os.Getenv("HOST_3")
	endpoints = []string{ep1, ep2, ep3}

	for _, testConfig := range testConfigs {
		fmt.Printf("\nTesting config: %+v\n", testConfig)
		totalConns = testConfig.Conns
		totalClients = testConfig.Clients
		keySize = testConfig.KeySize
		valSize = testConfig.ValSize
		switch testConfig.TestType {
		case "put":
			keySpaceSize = testConfig.KeySpaceSize
			putTotal = testConfig.Total
			putRate = testConfig.Rate
			putFunc(cmd, args)
		case "get":
			getArgs := []string{string(testConfig.Key)}
			rangeTotal = testConfig.Total
			rangeRate = testConfig.Rate
			rangeFunc(cmd, getArgs)
		}
	}
}
