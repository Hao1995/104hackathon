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
### Download the API service
```
go get -u -v github.com/Hao1995/104hackathon
```
### Enter your DB info
Copy an app.conf from example.conf
### Create DataBase
* Import data that I have been completed.   
    [[Data Link]](https://drive.google.com/drive/folders/1NGbvKBhIuSH1Dm_krbTeLBA5_97JJNqy?usp=sharing)
* Create a schema and name it '104hackathon-welfare'. Then produce the data by yourself.    
    * Create tables from floder "/sql/create_schema"
        1. users
        2. categories
        3. companies
        4. jobs
        5. welfares
        6. job_welfares
        7. welfare_user_score
        8. job_user_score
        9. train_action
        10. train_click
    * Synchronize 104hackathon open data.
        1. /api/sync/categories
        2. /api/sync/companies
        3. /api/sync/jobs
        4. /api/sync/train_click
        5. /api/sync/train_click/key
        6. /api/sync/train_action
    * Adding the foreign key.
        * /sql/create_schema/build_fk.sql
    * Add other info by manipulated.
        1. /api/welfare
        2. /api/user (Will add data to welfare_user_score)
        3. /api/user/welfare/score (Adjust the score of each welfare of each user)
        4. /api/job/welfare
        5. /api/user/job/score

## Manual
*-More detail-*

### country_id
Enter the city id to searching data.     

台北市          ：6001001   
新北市          ：6001002   
宜蘭縣          ：6001003   
基隆市          ：6001004   
桃園市          ：6001005   
新竹縣市        ：6001006   
苗栗縣          ：6001007   
台中市          ：6001008   
台中市(原台中縣) ：6001009   
彰化縣          ：6001010   
南投縣          ：6001011   
雲林縣          ：6001012   
嘉義縣市        ：6001013   
台南市          ：6001014   
台南市(原台南縣) ：6001015   
高雄市          ：6001016   
高雄市(原高雄縣) ：6001017   
屏東縣          ：6001018   
台東縣          ：6001019   
花蓮縣          ：6001020   
澎湖縣          ：6001021   
金門縣          ：6001022   
連江縣          ：6001023   