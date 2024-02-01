provider "aws" {
  region = "eu-west-1"
}

resource "aws_dynamodb_table" "products" {
  name         = "products"
  hash_key     = "id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "id"
    type = "S"
  }
}
