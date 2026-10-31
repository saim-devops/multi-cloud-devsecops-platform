output "cluster_endpoint" {
  value       = aws_eks_cluster.main.endpoint
  description = "EKS Cluster API server endpoint"
}

output "cluster_name" {
  value       = aws_eks_cluster.main.name
  description = "EKS Cluster Name"
}

output "ecr_repository_url" {
  value       = aws_ecr_repository.core_service.repository_url
  description = "AWS ECR Repository URL"
}

output "irsa_role_arn" {
  value       = aws_iam_role.core_service_irsa.arn
  description = "IAM Role for Service Account ARN"
}

output "kubeconfig_command" {
  value       = "aws eks --region ${var.aws_region} update-kubeconfig --name ${aws_eks_cluster.main.name}"
  description = "CLI command to update local kubeconfig"
}
