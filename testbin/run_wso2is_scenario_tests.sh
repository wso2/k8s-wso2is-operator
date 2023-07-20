#!/bin/bash

# Default values
ISHttpsUrl="https://localhost:9443"
input_dir=""
output_dir=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --is-https-url)
      ISHttpsUrl="$2"
      shift
      shift
      ;;
    --input-dir)
      input_dir="$2"
      shift
      shift
      ;;
    --output-dir)
      output_dir="$2"
      shift
      shift
      ;;
    *)
      shift
      ;;
  esac
done

# Creating the testbin directory if it doesn't exist
mkdir -p testbin

# Cloning the repository into the testbin/temp/product-is directory
echo -e "🔍 Cloning the repository into testbin/temp/product-is...\n"
git clone https://github.com/wso2/product-is.git testbin/temp/product-is

# Checking out the necessary branch/tag
echo -e "\n🌱 Checking out the necessary branch/tag...\n"
cd testbin/temp/product-is
git checkout master

# Navigating to the product-scenarios directory
echo -e "\n📂 Navigating to the product-scenarios directory...\n"
cd product-scenarios

# Creating the infrastructure.properties file if it doesn't exist
echo -e "\n🛠️ Creating the infrastructure.properties file...\n"
infrastructure_file="infrastructure.properties"
if [[ ! -f "$infrastructure_file" ]]; then
  echo "ISHttpsUrl=$ISHttpsUrl" > "$infrastructure_file"
  echo "ISSamplesHttpUrl=http://localhost:8080" >> "$infrastructure_file"
fi

# Modifying the infrastructure.properties file
echo -e "\n🛠️ Modifying the infrastructure.properties file...\n"
sed -i "s#ISHttpsUrl=.*#ISHttpsUrl=${ISHttpsUrl}#" "$infrastructure_file"

# Deploying test samples
echo -e "\n🚀 Deploying test samples...\n"
cd test-resources

if [[ "$OSTYPE" == "darwin"* ]]; then
  yes "A" | sh deploy-samples-mac.sh $ISHttpsUrl
elif [[ "$OSTYPE" == "linux-gnu" ]]; then
  yes "A" | sh deploy-samples-linux.sh $ISHttpsUrl
else
  echo -e "\n❌ Unsupported operating system: $OSTYPE\n"
  exit 1
fi

# Running the test script
echo -e "\n🔬 Running the test script...\n"
cd ..

# Check if input directory and output directory are provided
if [[ -z "$input_dir" ]] || [[ -z "$output_dir" ]]; then
  echo "❗ Input directory and output directory not provided. Using current directory as default."
  input_dir=$PWD
  output_dir=$PWD
fi

yes "A" | ./test.sh --input-dir "$input_dir" --output-dir "$output_dir"

# Removing the cloned repository
echo -e "\n🗑️ Removing the cloned repository...\n"
cd ../..
rm -rf testbin/temp/product-is
