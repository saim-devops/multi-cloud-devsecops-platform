resource "aws_ecr_repository" "core_service" {
  name                 = "core-service"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }

  tags = {
    Name = "core-service-repository"
  }
}

resource "aws_ecr_lifecycle_policy" "retention" {
  repository = aws_ecr_repository.core_service.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Retain latest 10 images to prevent storage bloat"
        selection = {
          tagStatus   = "any"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 14
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
