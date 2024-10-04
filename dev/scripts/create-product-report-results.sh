#!/bin/bash

# Function to create a product and return the product ID
create_product() {
  for i in 1 2 3; do
    response=$(curl -s -X POST http://localhost:8081/db/product \
      -H "Content-Type: application/json" \
      -d '{
        "fullName": "OpenShift",
        "shortName": "OCP",
        "contactEmail": "openshift@redhat.com"
      }')
    
    if [ -n "$response" ]; then
      product_id=$(echo $response | jq -r '.id')
      if [ "$product_id" != "null" ]; then
        echo $product_id
        return 0
      fi
    fi
    
    sleep 1
  done
  
  echo "Failed to create product" >&2
  return 1
}

# Function to report results
report_results() {
  local product_id=$1
  curl -X POST http://localhost:8081/report/results \
    -F "productId=$product_id" \
    -F "file=@./dev/junit-example.xml"
}

# Main script execution
product_id=$(create_product)
if [ $? -eq 0 ]; then
  echo "Created product with ID: $product_id"
  report_results $product_id
else
  echo "Error: Could not create product"
  exit 1
fi