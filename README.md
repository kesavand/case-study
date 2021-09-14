# case-study

#There are two microservices in this project 
1. evnt-producer that periodically writes user records to kafka
2. evnt-receiver - reads the user record from kafka and writes to maria-db and provides http endpoint (http://hos:port/users)end point to get the user records

The microservices can be deployed by both docker-compose and helm charts

# Deployment using docker compose 
The docker compose file is in deploments/compose 
The microservices can be started/stopped in docker compose by using the below command
  docker-compose -f compose.yml up -d 
  docker-compose -f compose.yml down 
  
  # Verify results by hitting the endpoint 
    http://localhost:8000/users
  
 # Deployment in minikube
 The microservices can be deployed in k8s environment by below steps
# start:
        helm repo add bitnami https://charts.bitnami.com/bitnami
	
        helm upgrade --install case-study-kafka  --namespace case-study --create-namespace bitnami/kafka --wait
	
        helm upgrade --install case-study-db --namespace case-study  bitnami/mariadb  --wait
	
        helm upgrade --install evnt-producer --namespace case-study ./deployment/k8s/charts/evnt-producer/ --wait
	
	ROOT_PASSWORD=$(kubectl get secret --namespace case-study case-study-db-mariadb -o jsonpath="{.data.mariadb-root-password}" | base64 --decode)
	
	echo $ROOT_PASSWORD
      
       (use the <user password got from the above step in the command below to deploy the event receiver microservice)
       
       helm upgrade --install evnt-rcvr  --namespace case-study ./deployment/k8s/charts/evnt-rcvr/ --set services.db.credentials=root:<ROOT_PASSWORD> --wait

# stop:
        helm ls --namespace case-study --short --all | xargs -L1 --no-run-if-empty helm delete --namespace case-study
  
       kubectl delete namespace case-study
        
 # Verify results by following the below commands
 minikube service list
|----------------------|-------------------------------------|--------------|---------------------------|
|      NAMESPACE       |                NAME                 | TARGET PORT  |            URL            |
|----------------------|-------------------------------------|--------------|---------------------------|
| case-study           | case-study-db-mariadb               | No node port |
| case-study           | case-study-kafka                    | No node port |
| case-study           | case-study-kafka-headless           | No node port |
| case-study           | case-study-kafka-zookeeper          | No node port |
| case-study           | case-study-kafka-zookeeper-headless | No node port |
| case-study           | evnt-rcvr                           |         8080 | http://192.168.49.2:30036 |
| default              | kubernetes                          | No node port |
| kube-system          | kube-dns                            | No node port |
| kubernetes-dashboard | dashboard-metrics-scraper           | No node port |
| kubernetes-dashboard | kubernetes-dashboard                | No node port |
|----------------------|-------------------------------------|--------------|---------------------------|

Use the below endpoint to verify results
http://192.168.49.2:30036/users

# Swagger documentation is present in docs/swagger


