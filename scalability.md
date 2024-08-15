## Service scalability and performance

I will make following changes for scalability and performance.

Divide the service in to 4 main components

    - Image Upload service 
        This will expose the Image upload APIs behind a loadbalancer. It is stateless so it can be horizntally scalled. 
        It will be deployed on K8s (EKS) with scaling baseed on application metrics.
        For storing images we will use blob storage (AWS S3) and all replicas will write images to this shared storage.
        It will also publish updates to Kafka topic about new images being uploaded and persisted in storage.

    - Image Processing Service
        This will read from Kafka topic and do the real time processing.
        Here we will use stream processing engine (Apache spark or Flink) for aggregations over time, data cleaning or any other similar use cases.
        It will write the data to its data store (Cassandra) which support high write performance.
        This also publish data using change data capture using Kafka.  

    - Image Statistics service
        This will read CDC data from image processing service and keep a copy in Postgress DB (AWS RDS).
        This will expose the APIs behind a loadbalancer. It is stateless so it can be horizntally scalled. 
        It will be deployed on K8s (EKS) with scaling baseed on application metrics.

    - Batch Processing Service
        This will read from Kafka topic and use Apache Spark for batch processing tasks.
        It will be deployed on AWS EMR.
        
## Further Questions

### On top of what is mentioned above, what would you use the batch/real-time pipeline for?

Batch pipeline can be used for:

    - create thumbnails
    - format conversion
    - pre compute different sizes
    - image storage reorganisation and cleanup
    - garbage collection

Real-time pipeline can be used for:

    - image statistics
    - service usage statistics
    - duplicate detection


### What kind of questions would you need to have answered in order to solve this task more concretely or in order to make better decisions on architecture and technologies?

- What is scale application need to handle in terms of DAUs and RPS.
- How long do we need to store image data.
- What are the storage requirements for real-time versus historical data.
- How many concurrent user do we need to support.
- What is the time delay between data being ingested and data being queried by user.
- What are the trade-offs between cost and performance.
