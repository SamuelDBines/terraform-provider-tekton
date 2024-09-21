#!/bin/bash

VERSION=$1
OWNER=$2
PROVIDER_NAME=$3

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
              "shasum": "$(sha256sum dist/terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_amd64.zip | awk '{ print $1 }')"
            },
            "darwin_arm64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_arm64.zip",
              "shasum": "$(sha256sum dist/terraform-provider-$PROVIDER_NAME_${VERSION}_darwin_arm64.zip | awk '{ print $1 }')"
            },
            "linux_amd64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_linux_amd64.zip",
              "shasum": "$(sha256sum dist/terraform-provider-$PROVIDER_NAME_${VERSION}_linux_amd64.zip | awk '{ print $1 }')"
            },
            "linux_arm64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_linux_arm64.zip",
              "shasum": "$(sha256sum dist/terraform-provider-$PROVIDER_NAME_${VERSION}_linux_arm64.zip | awk '{ print $1 }')"
            },
            "windows_amd64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_windows_amd64.zip",
              "shasum": "$(sha256sum dist/terraform-provider-$PROVIDER_NAME_${VERSION}_windows_amd64.zip | awk '{ print $1 }')"
            },
            "windows_arm64": {
              "filename": "terraform-provider-$PROVIDER_NAME_${VERSION}_windows_arm64.zip",
              "shasum": "$(sha256sum dist/terraform-provider-$PROVIDER_NAME_${VERSION}_windows_arm64.zip | awk '{ print $1 }')"
            }
          }
        }
      }
    }
  }
}
EOF