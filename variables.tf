variable "postgres_db" {
  description = "PostgreSQL database name"
  type        = string
  default     = "mydb"
}

variable "postgres_user" {
  description = "PostgreSQL user"
  type        = string
  default     = "user"
}

variable "postgres_password" {
  description = "PostgreSQL password"
  type        = string
  default     = "postgres"
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  default     = "supersecretkey"
  sensitive   = true
}

variable "frontend_port" {
  description = "Frontend external port"
  type        = number
  default     = 3000
}

variable "backend_port" {
  description = "Backend external port"
  type        = number
  default     = 8080
}

variable "database_port" {
  description = "Database external port"
  type        = number
  default     = 5432
}