package main

import (
	"os";
	"log";
	"fmt";
	"bufio";
	"strings";
	"strconv";
	"compress/gzip";
	"gopkg.in/amz.v3/s3";
	"gopkg.in/amz.v3/aws";
	"time";
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Usage: ec2-spot-bill from to")
		os.Exit(1)
	}

	bucketName := os.Getenv("S3_BUCKET")
	if bucketName == "" {
		log.Println("S3_BUCKET system variable should be set to s3 bucket name")
		os.Exit(1)
	}
	
	awsAccountId := os.Getenv("AWS_ACCOUNT")
	if awsAccountId == "" {
		log.Println("AWS_ACCOUNT system variable should be set to aws account id")
		os.Exit(1)
	}

	from, err := time.Parse("2006-01-02", os.Args[1])
    if err != nil {
		log.Println("Can not parse from time:" + err.Error())
		os.Exit(1)
    }
	to, err := time.Parse("2006-01-02", os.Args[2])
    if err != nil {
		log.Println("Can not parse to time:" + err.Error())
		os.Exit(1)
    }

    awsAtuh, err := aws.EnvAuth()
    if err != nil {
    	panic(err)
    }
    s3 := s3.New(awsAtuh, aws.USEast)
    if err != nil {
    	panic(err)
    }
    bucket, err := s3.Bucket(os.Getenv("S3_BUCKET"))
    if err != nil {
    	panic(err)
    }

    for date := from; date.Before(to); date = date.AddDate(0, 0, 1) {
    	chargeMap := map[string]float64{}
    	marker := ""
    	for {
    		prefix := awsAccountId + "." + date.Format("2006-01-02") + "-"
	    	listResp, err := bucket.List(prefix, "", marker, 10000)
	    	if err != nil {
	    		panic(err)
	    	}

	    	for _, key := range listResp.Contents {
	    		reader,err := bucket.GetReader(key.Key)
	    		if err != nil {
	    			panic(err)
	    		}
	    		uncompressingReader, err := gzip.NewReader(reader)
	    		if err != nil {
	    			panic(err)
	    		}
	    		scanner := bufio.NewScanner(uncompressingReader)
				for scanner.Scan() {
					line := scanner.Text()
					if line[0] == '#' {
						continue
					}
					columns := strings.Split(line, "\t")
					sir := columns[4]
					charge := strings.Split(columns[7], " ")[0]
					chargeFloat, err := strconv.ParseFloat(charge, 64)
					if err != nil {
						panic(err)
					}
					chargeMap[sir] += chargeFloat
				}
	    	}

	    	if listResp.IsTruncated {
	    		marker = listResp.Marker
    		} else {
    			break
    		}
		}

		for sir, charge := range chargeMap {
			fmt.Printf("%v\t%v\t%v USD\n", date.Format("2006-01-02"), sir, charge)
		}
    }
}