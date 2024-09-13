# gpu-mem-for-llm

This command line tool helps you calculate the estimated GPU memory required to serve large language models (LLMs) on a GPU. It takes into account the size of the model parameters, the precision used during training, and an optional overhead percentage. The tool provides a simple interface to input these parameters and outputs the estimated memory requirement.

## Installation

To use this tool, you need to have Go installed on your system. Then, you can install it using the following command:

```bash
go install github.com/ashprao/gpu-mem-for-llm@latest
```

After installation, the tool will be available in your `PATH` assuming you have set up go bin in your path.

## Usage

To calculate the memory requirement for a given model, use the following command format:

```bash
gpu-mem-for-llm --size <model-parameter-size> --fp32|--fp16|--bf16|--int8|--int4 [--overhead <percentage>] [--json]
```

Replace `<model-parameter-size>` with the size of your model parameters in millions (m) or billions (b), and choose the desired precision. If you want to include an overhead percentage, use the `--overhead` flag followed by a percentage value. Use the `--json` flag if you prefer the output in JSON format instead of human-readable text.

For example:

```bash
gpu-mem-for-llm --size 7b --fp16 --overhead 30
gpu-mem-for-llm --size 2b --bf16 --overhead 25 --json
```

## Flags

- `--size`: Specifies the size of the model parameters (e.g., "7b" for 7 billion). this flag is required.
- `--fp32`, `--fp16`, `--bf16`, `--int8`, `--int4`: These flags indicate the precision used during training and determine the memory requirement. Only one of them can be specified at a time.
- `--overhead`: this flag specifies an optional overhead percentage as an integer (e.g., "30" for 30%). The default value is 20% If not provided.
- `--json`: This flag indicates that the output should be in JSON format instead of human-readable text.

## Examples

Here are some examples of how to use the tool with different parameters:

```bash
gpu-mem-for-llm --size 100m --fp32
gpu-mem-for-llm --size 2b --bf16 --overhead 25
gpu-mem-for-llm --size 8b --int8 --overhead 40
```

## Contributing

Contributions are welcome! Please open an issue or create a pull request to share your ideas and improvements.
