package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const schemaRegistryUrl = "https://schema-registry.host.xxx"

func main() {
	schemaFlag := flag.String("schema", "", "specify the schema to generate dto")
	versionFlag := flag.Int("version", -1, "specify the schema version to generate dto")
	packageFlag := flag.String("package", "", "specify the package name for generated schema")

	flag.Parse()
	if (schemaFlag == nil || *schemaFlag == "") || (versionFlag == nil || *versionFlag == -1) {
		fmt.Println("Please specify the schema/version you want to generate using arguments")
		if schemaFlag == nil || *schemaFlag == "" {
			fmt.Println("invalid schema: ", *schemaFlag)
		}
		if versionFlag == nil || *versionFlag == -1 {
			fmt.Println("invalid version: ", *versionFlag)
		}
		return
	}

	fmt.Println("Schema to be pulled:", *schemaFlag)
	fmt.Println("Version to be used:", *versionFlag)
	fmt.Println("Using schema registry instance:" + schemaRegistryUrl)

	// schema registry client
	srClient := NewSchemaRegistryClient(schemaRegistryUrl, 300)
	schemaId, err := srClient.GetIdByVersion(*schemaFlag, *versionFlag)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Schema ID:", schemaId)

	schema, err := srClient.GetFullSchemaById(schemaId)
	if err != nil {
		log.Fatal(err)
	}

	schemaFile := "models/avro-dto/schemas/" + *schemaFlag + ".avsc"
	generationFolder := "models/avro-dto/" + *schemaFlag

	fmt.Println("schemaFile:", schemaFile)
	fmt.Println("generationFolder:", generationFolder)

	fmt.Println("Remvoing folder and schema file, if exists")

	err = os.RemoveAll(schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(generationFolder)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Writing schema into a file")

	f, err := os.Create(schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(schema)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating generation folder")

	err = os.Mkdir(generationFolder, 0o700)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Executing gogen-avro, package:", *packageFlag)
	var cmd *exec.Cmd
	if *packageFlag != "" {
		cmd = exec.Command("gogen-avro", "--short-unions=true", "--package="+*packageFlag, generationFolder, schemaFile)
	} else {
		cmd = exec.Command("gogen-avro", "--short-unions=true", generationFolder, schemaFile)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
