#!/bin/bash

# First API call
first_api_response=$(curl -s -X GET "http://localhost:8080/v1/link/create?expiration_time=5s" -H "Accept: application/json"  -H "Authorization:secret")
echo "$first_api_response"

# Extract the URL from the first API response
url=$(echo $first_api_response | jq -r '.link')

# Check if URL was extracted successfully
if [[ -z "$url" ]]; then
  echo "Failed to extract URL from the first API response."
  exit 1
fi

echo "Upload URL"
echo $url

# Second API call using the extracted URL
second_api_response=$(curl -v -X POST "$url" \
--form 'images=@"images/Canon_40D.jpg"' \
--form 'images=@"images/Olympus_C8080WZ.jpg"' \
--form 'images=@"images/Sony_HDR-HC3.jpg"' \
)
    
# Output the response from the upload API call
echo 
echo "Response from the second API call:"
echo "$second_api_response"

image_id=$(echo $second_api_response | jq -r '.["images"] | to_entries | .[0].value' )

echo
echo "First Image Id"
echo $image_id

# Download one image from image id
echo "Download one image"
curl -s -JLO "http://localhost:8080/v1/image/$image_id"

# Show statistics
echo "Get statistics"
curl  -H "Authorization:secret" http://localhost:8080/v1/statistics | jq