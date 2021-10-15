# DynamoDB Development

To test the dynamo db implementation, you have to run a [local instance of dynamo db](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html).

You can download dynamodb [here](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html).

## Run command

To test the database locally after download, open the command prompt and execute

```cmd
java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
```
