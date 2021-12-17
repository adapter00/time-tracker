
time tracking command for slack

## prepare

setup slack bot [see](https://slack.com/intl/ja-jp/help/articles/115005265703-%E3%83%AF%E3%83%BC%E3%82%AF%E3%82%B9%E3%83%9A%E3%83%BC%E3%82%B9%E3%81%A7%E5%88%A9%E7%94%A8%E3%81%99%E3%82%8B%E3%83%9C%E3%83%83%E3%83%88%E3%81%AE%E4%BD%9C%E6%88%90)

get app token and bot token

```
cp .env.sample .env
```
set `APP_TOKEN` and `BOT_TOKEN` to .env

## start

```
docker compose up -d
```

 start postgresql and slack cmd

## command

```
・start work
/tm start

・finish work
/tm stop [default now or yyyymm]

・start rest
/tm rstart 

・finish rest
/tm rstop 

・show month total duration
/tm show 


・show detail in month
/tm detail

・help
/tm help
```

