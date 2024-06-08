#!/bin/bash

# Navigate to each service directory and generate Swagger docs
(cd services/account_service && swag init --output docs)
(cd services/policy_service && swag init --output docs)
(cd services/transaction_service && swag init --output docs)
(cd services/rating_service && swag init --output docs)
(cd services/workflow_service && swag init --output docs)


# Combine the Swagger JSON files
# jq -s '.[0] * .[1] * .[2] * .[3]' \
#   services/account_service/docs/swagger.json \
#   services/policy_service/docs/swagger.json \
#   services/transaction_service/docs/swagger.json \
#   services/rating_service/docs/swagger.json > combined_swagger.json

# echo "Combined Swagger JSON created at combined_swagger.json"
