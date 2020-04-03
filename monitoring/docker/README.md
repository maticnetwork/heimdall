# Heimdall Monitoring

Monitor your Heimdall node with Prometheus and Grafana

## Follow below steps to set up monitoring on heimdall

**Step 1:**

Enable prometheus flag to true on heimdall `config.toml` file.

```
vim $HOME/.heimdalld/config.toml

# change prometheus flag to true

prometheus = true
```

**Step 2:**

You need to restart the heimdall node since config change 

**Step 3:**

Start both prometheus and grafana containers using following command:

```
docker-compose up -d
```

**Step 4:**

Open grafana at following URL

```
http://host_ip:3000
```

Login to grafana dashboard and edit the datasource HTTP to url http://host_ip:9090 and save.

## Grafana default login details

These login credentials can be changed according to user preferences once logged in:

```
username: admin
password: admin
```

## Grafana datasource configuration and navigation snapshots

Grafana uses web based APIs to connect to prometheus server for indexed data. For that, it needs prometheus host name.


1. Open grafana at url: http://host_ip:3000. Hover on the setting icon in the left pane and selet Data Sources: 



    ![Screenshot 2020-04-03 at 4 49 47 PM](https://user-images.githubusercontent.com/31979627/78356085-8bf3a480-75cc-11ea-9ed0-635edd495c96.png)


2. Notice that Prometheus sample datasource is added and click on the same:


     ![Screenshot 2020-04-03 at 4 50 14 PM](https://user-images.githubusercontent.com/31979627/78356289-e856c400-75cc-11ea-86da-e94d742a07f7.png)


3. Change the HTTP url to http://host_ip:9090 and save. After the success message, go to Grafana home:


     ![Screenshot 2020-04-03 at 5 14 53 PM](https://user-images.githubusercontent.com/31979627/78357564-4dabb480-75cf-11ea-9c9c-f6e8daadec47.png)


4. Click on the Home button on the left top and select Heimdall-Dashboard 


     ![Screenshot 2020-04-03 at 5 39 36 PM](https://user-images.githubusercontent.com/31979627/78359766-543c2b00-75d3-11ea-8b62-d8e8ee422191.png)

5. Notice Heimdall-Dashboard loaded as below

     ![Screenshot 2020-04-03 at 5 46 49 PM](https://user-images.githubusercontent.com/31979627/78359855-78980780-75d3-11ea-8cdf-8db0cb5ac4cc.png)
