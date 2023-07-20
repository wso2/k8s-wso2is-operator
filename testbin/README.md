# Identity Server Product Scenario Tests

This README explains how to run the default scenario tests locally for the WSO2 Identity Server product, deployed to a k8s cluster, using the `k8s-wso2is-operator` k8s operator.

The scenario tests help validate the functionality and performance of the product under different scenarios

## Automatically Run Scenario Tests

The whole process of running the WSO2 IS scenario tests could be automatically run by following the instructions in this section. This script will handle,

- Cloning the `product-is` repo
- Creating/Updating the `infrastructure.properties` file
- Deploying the sample tests via tomcat
- Running the tests
- Cleanup

The steps are,

### 1. Run the script from the root folder

```
./testbin/temp/run_scenario_tests.sh --is-https-url <ISHttpsURL> --input-dir <INPUT_DIR> --output-dir <OUTPUT_DIR>
```

For example:

```
./testbin/temp/run_scenario_tests.sh --is-https-url "https://dev.wso2is.com" --input-dir $PWD --output-dir $PWD
```

## Manually Scenario Tests

### 1. Setting Up

1. Clone the repository

```
git clone https://github.com/wso2/product-is.git
```

2. Checkout the necessary branch/tag for the Identity Server version you want to test.

```
cd product-is
git checkout master
```

3. Navigate to the `product-scenarios` directory:

```
cd product-scenarios
```

4. Open the existing `infrastructure.properties` file in a text editor.

5. Add or update the `ISHttpsUrl` configuration to point to the hostname of the ingress used for the k8s deployment of the WSO2 Identity Server instance. For example:

```
ISHttpsUrl=https://dev.wso2is.com
```

Note: Replace `dev.wso2is.com` with the actual hostname of the ingress.

6. Save the `infrastructure.properties` file.

7. Start the WSO2 Identity Server instance that you want to test, which should be deployed using the `k8s-wso2is-operator` and accessible through the specified hostname.

Note: The test suite assumes the deployed WSO2 Identity Server instances are with the default username and password `admin` and `admin` respectively.

### 2. Deploying Test Samples

To deploy the test samples, follow these steps:

1. Navigate to the `test-resources` directory:

```
cd test-resources
```

2. Run the appropriate deployment script based on your operating system:

- For macOS:

  ```
  sh deploy-samples-mac.sh <ISHttpsURL>
  ```

- For Linux:
  ```
  sh deploy-samples-linux.sh <ISHttpsURL>
  ```

Replace `<ISHttpsURL>` with the URL of your WSO2 Identity Server instance (e.g., `https://dev.wso2is.com`).

### 3. Running Tests

To execute the product scenario tests, follow these steps:

1. Navigate back to the `product-scenarios` directory:

```
cd ..
```

2. Run the test script with the following command, replacing `<INPUT_DIR>` and `<OUTPUT_DIR>` with the desired file paths on your machine:

```
./test.sh --input-dir <INPUT_DIR> --output-dir <OUTPUT_DIR>
```

For example:

```
./test.sh --input-dir $PWD --output-dir $PWD
```

Note: The above example assumes that you want to use the current directory as both the input and output directory.

### Notes

- By default, tests are executed against H2 DB. If you want to execute the tests locally against other DB types, you need to do the necessary datasource configurations in the `deployment.toml` file to point to another DB type in the WSO2 Identity Server pack. TG execution will happen against the JDBC userstore and different DB types.

- You can run only a part of the tests by commenting out unnecessary components from the root `pom.xml` file.
