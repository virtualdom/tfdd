# tfdd (Terraform Drift Detector)

tfdd is a CLI that explores a Terraform repository and reports if "drift" is detected. Drift is the phenomenon that occurs when a Terraform configuration does not match resources that exist in actuality. When this happens, `terraform plan` would typically propose changes.

Currently, tfdd only works with AWS but may work with other providers in the future.

## Todo

This is currently still a work in progress. Things to do before considering it GA

- Write tests and set up continuous integration
- Parse `plan` output to report what resources have drifted
- Publish to Homebrew
- Add log verbosity option

## Installation

Clone this repo and run `make`.

```bash
git clone https://github.com/virtualdom/tfdd.git
cd tfdd
make
```

Move `bin/tfdd` into your `$PATH`.

## Usage

To run a drift detection, run

```bash
tfdd detect /path/to/terraform/repo
```

This will use your default credentials as they are referenced in your Terraform configuration files. See [Terraform documentation](https://www.terraform.io/docs/providers/aws/index.html) for help setting up authentication correctly.

### Assuming an Audit Role

tfdd supports assuming an audit role before executing a drift detection. This might be helpful if `tfdd detect` is running as a dedicated cron and needs to fetch temporary AWS credentials before and throughout a drift detection. **This only works if your Terraform configuration files authenticate using AWS profiles.**

**Configure tfdd to be aware of concerned AWS accounts.** Run `tfdd configure` and follow the prompts for AWS account ID and account name. The account name determines the name of the AWS profile that temporary assumed role credentials are written as. You can also configure tfdd by directly editing the configuration file (defaults to `~/.tfdd/config` but can be customized using `--config-file` or the `TFDD_CONFIG_FILE` env variable with any tfdd command). A sample tfdd configuration file might look like the following

```json
{
 "profiles": [
  {
   "name": "account1",
   "account_id": "123456789012"
  },
  {
   "name": "account2",
   "account_id": "987654321098"
  },
  ...
 ]
}
```

**Create an auditing IAM role** in each concerned AWS account, and attach the AWS-managed [SecurityAudit](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_job-functions.html#jf_security-auditor) policy. Additionally, if your Terraform state files are stored in S3, ensure that this role has `ListBucket` and `GetObject` permissions for the proper state buckets and files. This is the role that will be assumed before and during execution, so ensure this role is assumable by whatever IAM entity tfdd will be running as.

Be sure to name them all the exact same role name. Any name will do and can be customized using the `--audit-role-name` flag in `tfdd detect`. By default, tfdd assumes the name is `audit-role`.

To run a drift detection while assuming an `audit-role` in each concerned AWS account every 30 minutes, run

```bash
tfdd detect /path/to/terraform/repo \
	--auth \
	--interval 1800 \
	--audit-role-name audit-role \
```

To customize what IAM entity tfdd uses before assuming the audit role, set either `--profile` or `--access-key` and `--secret-key` (and optionally, `--session-token`). Keep in mind that these credentials are not the ones that are used during the drift detection, but rather, they are used to assume the audit roles in all AWS accounts, so ensure that they correspond to a profile that is able to assume `audit-role` in all AWS accounts.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)
