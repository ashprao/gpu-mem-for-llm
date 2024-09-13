package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

// getParameterSize parses the parameter size value provided as a string and should be checked
// to be in a form such as 100m for 100 million or 7b for 7 billion. If any other string is provided
// an error is returned. If not, then the number is extracted and returned as an integer
func getParameterSize(param string) (int, error) {
	// Define regex pattern to match strings like "100m", "7b", etc.
	pattern := `^(\d+)([mbMB])$`
	re := regexp.MustCompile(pattern)

	// Match the input string against the pattern
	matches := re.FindStringSubmatch(param)
	if matches == nil {
		return 0, errors.New("invalid format; must be an integer followed by 'm' or 'b'")
	}

	// Extract number and unit from the matches
	numStr := matches[1]
	unit := matches[2]

	number, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %v", err)
	}

	switch unit {
	case "m", "M":
		return number * 1_000_000, nil
	case "b", "B":
		return number * 1_000_000_000, nil
	default:
		return 0, errors.New("invalid unit; must be 'm' or 'b'")
	}
}

// func get precision value from the flags provided
func getPrecision() (float32, error) {
	if fp32 {
		return 4, nil
	} else if fp16 {
		return 2, nil
	} else if bf16 {
		return 2, nil
	} else if int8 {
		return 1, nil
	} else if int4 {
		return 0.5, nil
	} else {
		return 0, errors.New("no precision flag provided")
	}
}

// calculateRequiredMemory returns the gpu memory required for serving llms
func calculateRequiredMemory(parameterSize int, precision float32, overhead float32) int {
	// Calculate the memory required for parameters
	memoryForParams := float32(parameterSize) * precision

	// Convert overhead to a percentage
	overhead /= 100
	overhead = 1 + overhead

	// Add overhead to the calculated memory
	totalMemoryRequired := int(memoryForParams * overhead)

	return totalMemoryRequired
}

// formatMemory takes an integer representing memory in bytes and returns a formatted string
// with the memory in megabytes or gigabytes, depending on the size. The returned string is rounded
// to two decimal places for readability.
func formatMemory(memoryBytes int) string {
	const (
		megabyte = 1_000_000
		gigabyte = 1_000_000_000
	)

	if memoryBytes >= gigabyte {
		return fmt.Sprintf("%.2f GB", float64(memoryBytes)/gigabyte)
	}

	return fmt.Sprintf("%d MB", memoryBytes/megabyte)
}

// checkMutuallyExclusivePrecisionFlags checks if multiple precision flags are provided at the same time.
// It returns an error if more than one flag is set, nil otherwise.
func checkMutuallyExclusivePrecisionFlags(cmd *cobra.Command) error {
	flags := []string{"fp32", "fp16", "bf16", "int8", "int4"}
	var count int

	for _, flag := range flags {
		if cmd.Flag(flag).Changed {
			count++
		}
	}

	if count > 1 {
		return errors.New("only one of --fp32, --fp16, --bf16, --int8, --int4 can be set at a time")
	}

	return nil
}

// rootCmd represents the base command identified by the 'Use' attribute
// when called without any subcommands. This name should be used in any
// build scripts.
var rootCmd = &cobra.Command{
	Use:   "gpu-mem-for-llm",
	Short: "Calculate memory required to serve LLm",
	Long: `Provide your model parameter size, precision and a percentage
overhead to calculate the estimated gpu memory required 
to run the model. 

For example:
./gpu-mem-for-llm --size 7b --fp16 --overhead 30

Flag details:
   --size: Specifies the size of the model parameters (e.g., "7b" for 7 billion). 
           This flag is required.
   --fp32 | --fp16 | --bf16 | --int8 | --int4: These flags indicate the precision 
           used during training and determine the memory requirement. 
           Only one of them can be specified at a time.
   --overhead: This flag specifies an optional overhead percentage as an integer 
           (e.g., "30" for 30%). 
           The default value is 20% if not provided.
`,
	Version: appVersion,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkMutuallyExclusivePrecisionFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		parameterSize, err := getParameterSize(size)
		if err != nil {
			fmt.Println(err)
			return
		}

		precision, err := getPrecision()
		if err != nil {
			fmt.Println(err)
			return
		}

		requiredMemory := calculateRequiredMemory(parameterSize, precision, float32(overhead))

		if jsonOutput {
			output := map[string]string{
				"mem_size": formatMemory(requiredMemory),
			}
			jsonData, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error generating JSON:", err)
				return
			}
			fmt.Println(string(jsonData))
		} else {
			fmt.Printf("Estimated memory required: %s\n", formatMemory(requiredMemory))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	// flags
	fp32       bool
	fp16       bool
	bf16       bool
	int8       bool
	int4       bool
	overhead   int
	size       string
	jsonOutput bool

	// versioning
	appVersion string = "0.1.0"
)

func init() {
	// Define a flag for the parameter size of the model in millions (m) or billions (b)
	rootCmd.Flags().StringVarP(&size, "size", "s", "", "model parameter size (e.g., 7b) - required")
	rootCmd.MarkFlagRequired("size")

	// Define a flag group for all these precision values - fp32, fp16, bf16, int8, int4100M
	// eg. --fp32, --fp16, --bf16, --int8, --int4
	// each of them is a boolean flag
	// only one of them can be provided at any given time.
	rootCmd.Flags().BoolVar(&fp32, "fp32", false, "use fp32 precision")
	rootCmd.Flags().BoolVar(&fp16, "fp16", false, "use fp16 precision")
	rootCmd.Flags().BoolVar(&bf16, "bf16", false, "use bf16 precision")
	rootCmd.Flags().BoolVar(&int8, "int8", false, "use int8 precision")
	rootCmd.Flags().BoolVar(&int4, "int4", false, "use int4 precision")
	rootCmd.MarkFlagsOneRequired("fp32", "fp16", "bf16", "int8", "int4")

	// Define a flag for the overhead
	rootCmd.Flags().IntVarP(&overhead, "overhead", "o", 20, "overhead as a percentage")

	// Define a flag for JSON output
	rootCmd.Flags().BoolVar(&jsonOutput, "json", false, "output results in JSON format")

	// Define a flag for version
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}
