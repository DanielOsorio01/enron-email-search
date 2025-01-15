terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
  }
}

# AWS Provider Configuration
provider "aws" {
  region = "us-east-1" # Replace with your desired AWS region
}

# 1. Define the VPC
resource "aws_vpc" "prod-vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name = "enron-main-vpc"
  }
}

# 2. Define Internet Gateway
resource "aws_internet_gateway" "gw" {
  vpc_id     = aws_vpc.prod-vpc.id

  tags = {
    Name = "enron-igw"
  }
}

# 3. Define a public and private subnet
resource "aws_subnet" "public-subnet" {
  vpc_id                  = aws_vpc.prod-vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true # Specify true to indicate that instances launched into the subnet should be assigned a public IP address. Default is false.
  
  depends_on = [ aws_internet_gateway.gw ]
  tags = {
    Name = "enron-public-subnet"
  }
}

resource "aws_subnet" "private-subnet" {
  vpc_id            = aws_vpc.prod-vpc.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-east-1b"
  tags = {
    Name = "enron-private-subnet"
  }
}

# 4. Create a custom route table for the public subnet
resource "aws_route_table" "rt" {
  vpc_id = aws_vpc.prod-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_internet_gateway.gw.id
  }

  tags = {
    Name = "enron-public-route-table"
  }
}

# 5. Associate the public subnet with the route table
resource "aws_route_table_association" "public-subnet-association" {
  subnet_id      = aws_subnet.public-subnet.id
  route_table_id = aws_route_table.rt.id
}

# Variable to store your dynamic IP
data "http" "ip" {
  url = "https://checkip.amazonaws.com"
}

# Output to display the current IP address
output "current_ip" {
  value = chomp(data.http.ip.response_body)
}


# Create a security group rule for allowing SSH access to the bastion host
resource "aws_security_group" "allow_ssh" {
  name        = "bastion_host_sg"
  description = "Allow SSH traffic"
  vpc_id      = aws_vpc.prod-vpc.id
  depends_on = [ data.http.ip ]
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    # Temporarily use "191.156.32.0/21" as the CIDR block for testing purposes
    cidr_blocks = ["191.156.32.0/20"]# ["${chomp(data.http.ip.response_body)}/32"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

}

# 6. Create a security group for the public subnet
resource "aws_security_group" "public-sg" {
  name        = "public-instances-sg"
  description = "Allow inbound traffic on port 80 and port 22"
  vpc_id      = aws_vpc.prod-vpc.id

  # Allow port 22 (SSH) only from the bastion host
  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.allow_ssh.id] # Reference the bastion host security group
  }

  # Allow port 80 (HTTP) from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "public-instances-sg"
  }
}

# 7. Create a key pair for SSH access
resource "tls_private_key" "instance-key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "enron-keypair" {
  key_name   = "enron-keypair"
  public_key = tls_private_key.instance-key.public_key_openssh
}

output "private_key" {
  value     = tls_private_key.instance-key.private_key_openssh
  sensitive = true
}



# Create the Bastion Host
resource "aws_instance" "bastion-host" {
  ami                    = "ami-05576a079321f21f8"
  instance_type          = "t2.micro"
  subnet_id              = aws_subnet.public-subnet.id
  key_name               = "enron-keypair"
  vpc_security_group_ids = [aws_security_group.allow_ssh.id]
  user_data = <<-EOF
              #!/bin/bash
              sudo yum update -y
              echo "ZINCSEARCH_URL=http://${aws_instance.database.private_ip}:4080" >> /etc/environment
              source /etc/environment
              EOF
  
  depends_on             = [aws_key_pair.enron-keypair, data.http.ip, aws_security_group.allow_ssh, aws_instance.database]
  tags = {
    Name = "bastion-host"
  }
}

output "bastion-ip" {
  value = aws_instance.bastion-host.public_ip

}

# Security Group for Private Instances
resource "aws_security_group" "private-sg" {
  name        = "private-instances-sg"
  description = "Allow communication within private network"
  vpc_id      = aws_vpc.prod-vpc.id

  # Allow port 22 (SSH) only from the bastion host
  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.allow_ssh.id] # Reference the bastion host security group
  }

  # Allow other necessary application-specific traffic (e.g., backend communication)
  ingress {
    description = "Allow communication within private network"
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"] # Replace with your VPC CIDR block
  }

  # Allow outbound traffic to the internet via NAT or other routes
  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "private-instances-sg"
  }
}


# 7. Create IAM role for allowing EC2 instances to access ECR
resource "aws_iam_role" "ec2-role" {
  name = "ec2-ecr-access-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "ec2.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Create a NAT Gateway for the private subnet
# Elastic IP for NAT Gateway
resource "aws_eip" "nat-eip" {
  domain = "vpc"

  tags = {
    Name = "nat-eip"
  }
}

# NAT Gateway in a Public Subnet
resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat-eip.id
  subnet_id     = aws_subnet.public-subnet.id

  tags = {
    Name = "nat-gateway"
  }
}

# Route Table for Private Subnet
resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.prod-vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }

  tags = {
    Name = "private-route-table"
  }
}

# Associate Private Subnet with Route Table
resource "aws_route_table_association" "private_association" {
  subnet_id      = aws_subnet.private-subnet.id
  route_table_id = aws_route_table.private_route_table.id
}

# Attach the AmazonEC2ContainerRegistryReadOnly policy to the role
resource "aws_iam_role_policy_attachment" "ecr-readonly" {
  role       = aws_iam_role.ec2-role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

# Create an instance profile
resource "aws_iam_instance_profile" "ec2-instance-profile" {
  name = "ec2-ecr-access-instance-profile"
  role = aws_iam_role.ec2-role.name
}

# 7. Define the instances after creating the key pair

resource "aws_instance" "frontend" {
  ami                    = "ami-05576a079321f21f8"
  instance_type          = "t2.micro"
  iam_instance_profile   = aws_iam_instance_profile.ec2-instance-profile.name
  subnet_id              = aws_subnet.public-subnet.id
  key_name               = "enron-keypair"
  vpc_security_group_ids = [aws_security_group.public-sg.id]
  depends_on             = [aws_key_pair.enron-keypair, aws_instance.backend]
  user_data              = <<-EOF
                          #!/bin/bash
                          sudo yum update -y
                          sudo yum -y install docker
                          sudo service docker start
                          sudo usermod -a -G docker ec2-user
                          sudo chmod 666 /var/run/docker.sock
                          aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 418272755608.dkr.ecr.us-east-1.amazonaws.com
                          docker pull 418272755608.dkr.ecr.us-east-1.amazonaws.com/frontend
                          docker run -d -p 80:80 --name frontend -e VUE_APP_API_URL=http://${aws_instance.backend.private_ip}:3000 418272755608.dkr.ecr.us-east-1.amazonaws.com/frontend:latest
                          EOF
  tags = {
    Name = "enron-web-server"
  }
}

output "frontend_ip" {
  value = aws_instance.frontend.public_ip
}

output "frontend_private_ip" {
  value = aws_instance.frontend.private_ip
}

resource "aws_instance" "backend" {
  ami           = "ami-05576a079321f21f8"
  instance_type = "t2.micro"
  iam_instance_profile = aws_iam_instance_profile.ec2-instance-profile.name
  subnet_id     = aws_subnet.private-subnet.id
  key_name      = "enron-keypair"
  vpc_security_group_ids = [ aws_security_group.private-sg.id ]
    user_data              = <<-EOF
                          #!/bin/bash
                          sudo yum update -y
                          sudo yum -y install docker
                          sudo service docker start
                          sudo usermod -a -G docker ec2-user
                          sudo chmod 666 /var/run/docker.sock
                          aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 418272755608.dkr.ecr.us-east-1.amazonaws.com
                          docker pull 418272755608.dkr.ecr.us-east-1.amazonaws.com/frontend
                          docker run -d -p 3000:3000 --name backend -e DB_HOST=http://${aws_instance.database.private_ip}:4080 -e DB_USER=admin -e DB_PASSWORD=Complexpass#123 -e SERVER_PORT=3000 418272755608.dkr.ecr.us-east-1.amazonaws.com/backend:latest
                          EOF

  
  depends_on    = [aws_key_pair.enron-keypair, aws_instance.database]

  tags = {
    Name = "enron-backend-server"
  }
}

output "backend_private_ip" {
  value = aws_instance.backend.private_ip
}

resource "aws_instance" "database" {
  ami           = "ami-05576a079321f21f8"
  instance_type = "t2.micro"
  iam_instance_profile = aws_iam_instance_profile.ec2-instance-profile.name
  subnet_id     = aws_subnet.private-subnet.id
  key_name      = "enron-keypair"
  vpc_security_group_ids = [ aws_security_group.private-sg.id ]
  user_data              = <<-EOF
                          #!/bin/bash
                          sudo yum update -y
                          sudo yum -y install docker
                          sudo service docker start
                          sudo usermod -a -G docker ec2-user
                          sudo chmod 666 /var/run/docker.sock
                          mkdir data
                          docker run -d -v ~/data -e ZINC_DATA_PATH="/data" -p 4080:4080 -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 -e GIN_MODE=release --name zincsearch public.ecr.aws/zinclabs/zincsearch:latest
                          EOF

  depends_on    = [aws_key_pair.enron-keypair]
  tags = {
    Name = "enron-database-server"
  }
}

output "database_private_ip" {
  value = aws_instance.database.private_ip
}

