package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mangonet-labs/mgo-go-sdk/transaction"
)

func main() {
	var (
		format     = flag.String("format", "auto", "Input format: auto, hex, base64, json")
		output     = flag.String("output", "pretty", "Output format: pretty, json")
		inputFile  = flag.String("file", "", "Read from file instead of stdin")
		rawMessage = flag.String("message", "", "Raw message to decode")
	)
	flag.Parse()

	var input string
	var err error

	// Get input from various sources
	if *rawMessage != "" {
		input = *rawMessage
	} else if *inputFile != "" {
		content, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		input = string(content)
	} else {
		// Read from stdin
		fmt.Println("Enter raw transfer message (press Ctrl+D when done):")
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		input = strings.Join(lines, "\n")
	}

	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Fprintf(os.Stderr, "No input provided\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create decoder and decode the message
	decoder := transaction.NewRawMessageDecoder()
	var decoded *transaction.DecodedTransferMessage

	switch *format {
	case "auto":
		decoded, err = decoder.DecodeRawMessage(input)
	case "hex":
		decoded, err = transaction.DecodeTransactionHex(input)
	case "base64":
		decoded, err = transaction.DecodeTransactionBase64(input)
	case "json":
		decoded, err = decoder.DecodeJSONMessage(input)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n", *format)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding message: %v\n", err)
		os.Exit(1)
	}

	// Output the result
	switch *output {
	case "pretty":
		fmt.Println(decoded.PrettyPrint())
	case "json":
		jsonOutput, err := decoded.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting to JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(jsonOutput)
	default:
		fmt.Fprintf(os.Stderr, "Unknown output format: %s\n", *output)
		os.Exit(1)
	}
}
