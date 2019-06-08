# The Welfare Calculator of Applying Job
Competition Name: 2018 104hackathon - Mobile App/Web Service - Silver Award
Group Name: border-bottom solid 1000px #000

## 2018 104hackathon - AppWebService    
We use the data that include user searching keyword to evaluate the welfare is bad or good.     
Different keyword would get a different result.   
Such as if you searching "React" or "Front-end".   
Maybe we will get a completely different result of welfare score.   

## Demo
![border-bottom solid 1000px #000](https://github.com/Hao1995/104hackathon/blob/master/104hackathon.gif "border-bottom solid 1000px #000")

## Data Source
[2018-104Hackathon-AppWebService](https://github.com/104corp/2018-104Hackathon-AppWebService)

## Install
### Create DataBase
* Import ata that I have been completed.   
    [[Data Link]](https://drive.google.com/open?id=15TetTLofxwuY7VzjKVlj-Vx61mKKopFz)
* Create a schema and name it '104hackathon-welfare'. Then produce the data by yourself.   
    *Not yet complete the detail.*
### Download the API service
```
go get -u -v github.com/Hao1995/104hackathon
```
### Enter your DB info
Copy your app.conf from example.conf

## Manual
More detail  
### country_id
\# `country` will be change to `city`  
Enter the city id to searching data.   
| ID | City |
|---|---|
| 6001001 | 台北市|
| 6001002 | 新北市|
| 6001003 | 宜蘭縣|
| 6001004 | 基隆市|
| 6001005 | 桃園市|
| 6001006 | 新竹縣市|
| 6001007 | 苗栗縣|
| 6001008 | 台中市|
| 6001009 | 台中市(原台中縣)|
| 6001010 | 彰化縣|
| 6001011 | 南投縣|
| 6001012 | 雲林縣|
| 6001013 | 嘉義縣市|
| 6001014 | 台南市|
| 6001015 | 台南市(原台南縣)|
| 6001016 | 高雄市|
| 6001017 | 高雄市(原高雄縣)|
| 6001018 | 屏東縣|
| 6001019 | 台東縣|
| 6001020 | 花蓮縣|
| 6001021 | 澎湖縣|
| 6001022 | 金門縣|
| 6001023 | 連江縣|