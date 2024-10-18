# Github Actions
## Lambda layer
```yaml
name: Lambda Layer deployment
on:
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest

    permissions:
      id-token: write
      contents: read

      steps:
        - name: Code checkout
          uses: actions/checkout@v4
          
        - name: Install python runtime
          uses: actions/setup-python@v5
          with:
            python-version: "3.x"

        - name: Configure AWS credentials
          uses: aws-actions/configure-aws-credentials@v4
          with:
            role-to-assume: ${{ secrets.AWS_ROLE_ARN_DEV }}
            aws-region: ${{ secrets.AWS_REGION }}
            
        - name: Create build package
            run: |
              chmod +x script/build_package.sh
              script/build_package.sh
              
        - name: Deploy Lambda Layer
            run: |
              aws lambda publish-layer-version \
              --layer-name layer-name \
              --description "Updated Lambda layer" \
              --zip-file fileb://lambda-layer.zip \
              --compatible-runtimes python3.12 
       
```