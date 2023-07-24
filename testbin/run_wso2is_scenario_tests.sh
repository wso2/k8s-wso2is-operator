#!/bin/bash

# Default values
ISHttpsUrl="https://localhost:9443"
input_dir=""
output_dir=""
hostname="localhost"

# KeyStore
keystore_name="wso2carbon.jks"
keystore_path="scenarios-commons/src/main/resources/keystores/products/$keystore_name"
keystore_password="wso2carbon"
key_password="wso2carbon"

# Function to extract hostname from URL without https:// prefix
extract_hostname() {
  local url="$1"
  local hostname=$(echo "$url" | awk -F[/:] '{print $4}')
  echo "$hostname"
}

# Parse named arguments
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --is-https-url)
      ISHttpsUrl="$2"
      hostname=$(extract_hostname "$ISHttpsUrl")
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

if [ ! -d "testbin/temp/product-is/product-scenarios" ]; then
  echo -e "\n🔍 Downloading the 'product-scenarios'...\n"
  mkdir -p testbin/temp/product-is
  cd testbin/temp/product-is

  wget --progress=bar:force -O product-scenarios.zip https://github.com/wso2/product-is/archive/master.zip
  echo -e "\n📥 Downloading finished\n"

  unzip -q product-scenarios.zip
  mv product-is-master/product-scenarios ./
  rm -rf product-is-master product-scenarios.zip

  cd ../../../
fi

# Navigating to the product-scenarios directory
echo -e "\n📂 Navigating to the product-scenarios directory...\n"
cd testbin/temp/product-is/product-scenarios

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

# Disable ingress addon
echo -e "\n🌑 Disabling the 'ingress' addon...\n"
minikube addons disable ingress

# Check if 'mkcert' secret exists and delete it if it does
echo -e "\n🔒 Checking for existing 'mkcert' secret...\n"
if kubectl -n kube-system get secret mkcert &>/dev/null; then
  echo -e "Deleting existing 'mkcert' secret...\n"
  kubectl -n kube-system delete secret mkcert
fi

# Generate and install the certificate
echo -e "\n🔒 Generating and installing the certificate...\n"
mkcert $hostname
kubectl -n kube-system create secret tls mkcert --key ./$hostname-key.pem --cert ./$hostname.pem

# Configure ingress addon
echo -e "\n✅ Configuring the 'ingress' addon...\n"
echo "kube-system/mkcert" | minikube addons configure ingress
minikube addons enable ingress

# Install local CA certificates
echo -e "\n⚡️ Installing local CA certificates...\n"
sudo mkdir -p /usr/local/share/ca-certificates
sudo cp "$(mkcert -CAROOT)"/rootCA.pem /usr/local/share/ca-certificates/mkcert.crt
sudo update-ca-certificates

# Convert the certificate and key to DER format
echo -e "\n🔐 Converting the certificate and key to DER format...\n"
openssl x509 -outform der -in ./$hostname.pem -out ./$hostname.der

cd ..

# Add the certificate to the JKS store
echo -e "\n⚙️ Adding the certificate to the JKS store...\n"
keytool -importcert -alias $hostname -file ./test-resources/$hostname.der -keystore "$keystore_path" -storepass "$keystore_password" -noprompt

# Clean up temporary files
rm ./test-resources/$hostname.der
rm ./test-resources/$hostname-key.pem
rm ./test-resources/$hostname.pem

# Running the test script
echo -e "\n🔬 Running the test script...\n"

# Check if input directory and output directory are provided
if [[ -z "$input_dir" ]] || [[ -z "$output_dir" ]]; then
  echo "❗ Input directory and output directory not provided. Using current directory as default."
  input_dir=$PWD
  output_dir=$PWD
fi

yes "A" | ./test.sh --input-dir "$input_dir" --output-dir "$output_dir"

# Removing the cloned repository
echo -e "\n🗑️ Removing the cloned repository...\n"
cd ../../../../
rm -rf testbin/temp
