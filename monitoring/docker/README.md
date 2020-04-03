## Heimdall Monitoring

Monitor your Heimdall server with Prometheus and Grafana

## Follow below steps to set up monitoring on heimdall

step1:
```
enable prometheus falg to true on heimdall config.toml file
File path = $HOME/.heimdalld/config
```
step2:
```
restart heimdall
```
step3:
```
start both prometheus and grafana containers
usage:
docker-compose up -d
```
step4:
```
open grafana at url: http://host_ip:3000
Login to grafana dashboard and edit the datasource HTTP to url http://host_ip:9090 and save
```
# grafan default login details 
```
username: admin
password: admin
These login credentials can be reset according to user preferences
```
# Grafana datasource configuration and navigation snapshots
open grafana at url: http://host_ip:3000. Hover on the setting icon in the left pane and selet Data Sources 
![Screenshot 2020-04-03 at 4 49 47 PM](https://user-images.githubusercontent.com/31979627/78356085-8bf3a480-75cc-11ea-9ed0-635edd495c96.png)

Notice that Prometheus datasource is added and clik on the same
![Screenshot 2020-04-03 at 4 50 14 PM](https://user-images.githubusercontent.com/31979627/78356289-e856c400-75cc-11ea-86da-e94d742a07f7.png)

change the HTTP url to http://host_ip:9090 and save. Post sucess message, go to grafana home.
![Screenshot 2020-04-03 at 5 14 53 PM](https://user-images.githubusercontent.com/31979627/78357564-4dabb480-75cf-11ea-9c9c-f6e8daadec47.png)

Click on the Home button on the left top and select Heimdall-Dashboard 
![Screenshot 2020-04-03 at 5 39 36 PM](https://user-images.githubusercontent.com/31979627/78359766-543c2b00-75d3-11ea-8b62-d8e8ee422191.png)

Notice Heimdall-Dashboard loaded as below
![Screenshot 2020-04-03 at 5 46 49 PM](https://user-images.githubusercontent.com/31979627/78359855-78980780-75d3-11ea-8cdf-8db0cb5ac4cc.png)



