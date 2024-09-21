#!/bin/bash

VERSION=$1
OWNER=$2
PROVIDER_NAME=$3

if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  CHECKSUM_CMD="shasum -a 256"
else
  # Assume Linux (and other systems that support sha256sum)
  CHECKSUM_CMD="sha256sum"
fi

echo "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_amd64.zip | awk '{ print $1 }')"

# Create the JSON structure
cat <<EOF > terraform-registry-manifest.json
{
  "providers": {
    "github.com/$OWNER/$PROVIDER_NAME": {
      "versions": {
        "$VERSION": {
          "platforms": {
            "darwin_amd64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_amd64.zip",
              "shasum": "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_amd64.zip | awk '{ print $1 }')"
            },
            "darwin_arm64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_arm64.zip",
              "shasum": "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_arm64.zip | awk '{ print $1 }')"
            },
            "linux_amd64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_linux_amd64.zip",
              "shasum": "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_linux_amd64.zip | awk '{ print $1 }')"
            },
            "linux_arm64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_linux_arm64.zip",
              "shasum": "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_linux_arm64.zip | awk '{ print $1 }')"
            },
            "windows_amd64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_windows_amd64.zip",
              "shasum": "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_windows_amd64.zip | awk '{ print $1 }')"
            },
            "windows_arm64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_windows_arm64.zip",
              "shasum": "$($CHECKSUM_CMD dist/terraform-provider-$PROVIDER_NAME_${VERSION}_windows_arm64.zip | awk '{ print $1 }')"
            }
          }
        }
      }
    }
  }
}
EOF