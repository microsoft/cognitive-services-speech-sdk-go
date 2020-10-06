# Overview

This project contains Golang binding for the Microsoft Cognitive Service Speech SDK.

# Getting Started

This project requires Go 1.13

## Linux

### Compiling

First, you new need to set the environment variables for Go to find the SDK:


```bash
export CGO_CFLAGS="-I/path/to/SDKs/include/c_api"
export CGO_LDFLAGS="-L/path/to/SDKs/library -lMicrosoft.CognitiveServices.Speech.core"

```

After that we are ready to compile the package

### Running

To run applications consuming this project, we need to add the library to the path (if it is not there already)

```bash
export LD_LIBRARY_PATH="/path/to/SDKs/library:$LD_LIBRARY_PATH"
```

### Running Tests

In addition to the environment variables needed to run applications, running tests requires setting the following variables:

```bash
export TEST_SUBSCRIPTION_KEY="your_subscription_key"
export TEST_SUBSCRIPTION_REGION="your_region"
```

# Reference

Reference documentation for these packages is available at http://aka.ms/csspeech/goref

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
