# Enron Email Search App

Download the dataset [here](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz).

For deploying in AWS, you need to have the Terraform and AWS CLI installed, also set up credentials for AWS.

## Terraform Commands
(I recommend using WSL if you are on Windows)
* Initialize the terraform directory
```bash
terraform init
```

* Plan the deployment
```bash
terraform plan
```

* Apply the deployment
```bash
terraform apply
```

This should deploy all the resources in AWS and output the public IP of the frontend instance.

## Access to the bastion host
If you want to access the servers, you can use the provided bastion host to access the private instances.

First, export the private key of the bastion host
```bash
terraform output -raw private_key > private_key.pem
```

Start the ssh-add agent
```bash
eval "$(ssh-agent -s)""
```

Limit the permissions of the private key
```bash
chmod 600 private_key.pem
```

Then, add the key to your ssh agent 
```bash
ssh-add private_key.pem
```

Add the bastion host to your ssh config.
Create a file named `config` in the `~/.ssh` directory and add the following content
```bash
Host bastion-host
    HostName <bastion_public_ip>
    User ec2-user
    IdentityFile ~/.ssh/private_key.pem
    ForwardAgent yes
```

Finally, ssh into the bastion host
```bash
ssh bastion-host
```

From the bastion host, you can ssh into the private instances
```bash
ssh -A ec2-user@<backend_private_ip>
```

## Commands to update frontend image
```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 418272755608.dkr.ecr.us-east-1.amazonaws.com
docker tag enron-email-search-frontend:latest 418272755608.dkr.ecr.us-east-1.amazonaws.com/frontend:latest
docker push 418272755608.dkr.ecr.us-east-1.amazonaws.com/frontend:latest
```
