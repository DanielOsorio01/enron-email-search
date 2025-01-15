# Enron Email Search App

Download the dataset [here](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz).

For deploying in AWS, you need to have the Terraform and AWS CLI installed, also set up credentials for AWS.

## Terraform Commands
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

Then, add the key to your ssh agent (I recommend using WSL if you are on Windows)
```bash
ssh-add private_key.pem
```

Add the bastion host to your ssh config
```bash
echo "Host bastion-host\n  HostName <bastion_ip>\n  User ec2-user\n  IdentityFile private_key.pem" >> ~/.ssh/config
```

Finally, ssh into the bastion host
```bash
ssh bastion-host
```

From the bastion host, you can ssh into the private instances
```bash
ssh -A ec2-user@<backend_private_ip>
```
