ec2-spot-price aggregates ec2 spot instances charges by spot request and day and prints result in TSV format

usage:
``` bash
AWS_ACCOUNT=... AWS_ACCESS_KEY_ID=... AWS_SECRET_ACCESS_KEY=... S3_BUCKET=... ec2-spot-bill from to
```

**AWS_ACCOUNT** is long number representing account id

**S3_BUCKET** is a S3 bucket where spot instances charge feed is gathered. Details: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-data-feeds.html

**AWS_ACCESS_KEY_ID** and **AWS_SECRET_ACCESS_KEY** are keys with read and list access to this bucket

**from** and **to** are dates in format YYYY-MM-dd, for example 2015-04-21

