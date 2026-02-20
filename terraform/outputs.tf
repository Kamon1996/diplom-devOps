output "app_server_public_ip" {
  description = "Public IP of the app server"
  value       = aws_eip.app.public_ip
}

output "monitoring_server_public_ip" {
  description = "Public IP of the monitoring server"
  value       = aws_eip.monitoring.public_ip
}

output "ssh_app" {
  description = "SSH command for app server"
  value       = "ssh -i ~/.ssh/id_rsa ubuntu@${aws_eip.app.public_ip}"
}

output "ssh_monitoring" {
  description = "SSH command for monitoring server"
  value       = "ssh -i ~/.ssh/id_rsa ubuntu@${aws_eip.monitoring.public_ip}"
}

output "app_url" {
  description = "Application URL"
  value       = "http://${aws_eip.app.public_ip}"
}

output "grafana_url" {
  description = "Grafana URL"
  value       = "http://${aws_eip.monitoring.public_ip}:3000"
}

output "prometheus_url" {
  description = "Prometheus URL"
  value       = "http://${aws_eip.monitoring.public_ip}:9090"
}
