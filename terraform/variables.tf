variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name used for resource naming and tagging"
  type        = string
  default     = "habits-tracker"
}

# --- Network ---

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
  description = "Public subnet CIDR block"
  type        = string
  default     = "10.0.1.0/24"
}

# --- EC2 ---

variable "app_instance_type" {
  description = "EC2 instance type for app server"
  type        = string
  default     = "t3.small"
}

variable "monitoring_instance_type" {
  description = "EC2 instance type for monitoring server"
  type        = string
  default     = "t3.medium"
}

variable "app_volume_size" {
  description = "Root volume size (GB) for app server"
  type        = number
  default     = 20
}

variable "monitoring_volume_size" {
  description = "Root volume size (GB) for monitoring server"
  type        = number
  default     = 30
}

# --- SSH ---

variable "ssh_public_key_path" {
  description = "Path to the public SSH key file"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "allowed_ssh_cidr" {
  description = "CIDR block allowed for SSH access (your IP). Get it: curl -s ifconfig.me"
  type        = string
}
