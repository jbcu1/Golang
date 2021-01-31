package DBElastic

import (
	"PTwork/models"
	"fmt"
	"github.com/olivere/elastic"
	"golang.org/x/net/context"
	"time"
)


var client *elastic.Client
var err error


//Initialise connection to elasticSearch
func init(){

	Ctx := context.Background()

	client,err = elastic.NewClient(
		elastic.SetSniff(true),
		elastic.SetURL(ElasticURL),
		elastic.SetHealthcheckInterval(5*time.Second),
		elastic.SetBasicAuth(ElasticUser,ElasticPassword),
		)

	if err!=nil{
		return
	}

	info, code, err:=client.Ping(ElasticURL).Do(Ctx)

	if err!=nil{
		return
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return
}

func indexExist(index string)(exist bool){

	ctx:=context.Background()
	exist, err:=client.IndexExists(index).Do(ctx)
	if err!=nil{
		_ = fmt.Errorf("some errors %v", err)
	}
	return exist
}


//Bulk method to send posts to elasticSearch index
//Collect 10 thousand post then send it to index
func BulkSendPost(dates []models.DomainData) error {

	index:="pdns-test"
	lengthDomain :=len(dates)/10000
	ctx:=context.Background()
	n:=0

	for i:=0;i< lengthDomain+1;i++{
		bulkRequest:=client.Bulk()

		for j:=0;j<10000;j++{

			req:=elastic.NewBulkIndexRequest().
				Index(index).
				Type("_doc").
				Doc(dates[n])
			bulkRequest=bulkRequest.Add(req)
			n++
			if n>len(dates) -1{
				break

			}
		}

		bulkResponse,err:=bulkRequest.Do(ctx)

		if err!=nil{
			return err
		}

		if bulkResponse!=nil{

		}

		fmt.Printf("Indexed %v posts\n", n)

	}

	fmt.Printf("Indexed all posts to elastic: %v posts.\n", n)

	return nil
}


//Bulk method to tranco to send posts to elasticSearch index
//Collect 10 thousand post then send it to index
func BulkSendPostTranco(dates []models.TrancoDomainData) error {
	index:="pdns-test"
	lengthDomain :=len(dates)/10000
	ctx:=context.Background()
	n:=0
	for i:=0;i< lengthDomain;i++{
		bulkRequest:=client.Bulk()
		for j:=0;j<10000;j++{

			req:=elastic.NewBulkIndexRequest().
				Index(index).
				Type("_doc").
				Doc(dates[n])
			bulkRequest=bulkRequest.Add(req)
			n++
			if n>len(dates) -1{
				break
			}
		}

		bulkResponse,err:=bulkRequest.Do(ctx)
		if err!=nil{
			return err
		}
		if bulkResponse!=nil{

		}
		fmt.Printf("Indexed %v posts\n", n)

	}

	fmt.Printf("Indexed %v posts to elastic\n", n)


	return nil
}

//Delete one post from elasticSearch
//Find document at id post
func DeletePost(id string) error{

	index:="pdns-test"
	ctx:=context.Background()
	_,err:=client.Delete().Index(index).Type("_doc").Id(id).Do(ctx)

	if err!=nil{
		return err
	}

	return nil
}

//Send one post to elasticSearch
func SendPost(data interface{}) (string,error) {
	index:="pdns-test"
	ctx:=context.Background()

	put,err:=client.Index().
		Index(index).
		Type("_doc").
		BodyJson(data).
		Refresh("wait_for").
		Do(ctx)

	if err!=nil{
		_ = fmt.Errorf("some errors %s", err)
		return "",err
	}

	fmt.Printf("Indexed post %s to index %s\n", put.Id,put.Index )
	return put.Id,nil
}

