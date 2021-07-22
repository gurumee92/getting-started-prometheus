# 서비스 메트릭 모니터링하기 (1) prometheus-nginxlog-exporter

![logo](../../logo.png)

## 개요

이 문서에서는 `prometheus-nginxlog-exporter`를 이용해서 `Nginx`의 액세스 로그를 분석하여 RPS, 상태 코드 개수 등의 서비스 메트릭을 수집한다. 그 후 `Grafana`, `Prometheus`를 이용해서 시스템 메트릭을 모니터링할 수 있는 대시보드를 구축하는 것에 대하여 다룬다. 자세한 내용은 다음과 같다.

* Nginx와 설치
* prometheus-nginxlog-exporter와 설치
* 메트릭 수집을 위한 각 컴포넌트 설정
* 서비스 메트릭 모니터링을 위한 Grafana 대시보드 구축

이 문서에서 진행되는 실습 코드는 편의성을 위해 로컬 `Docker` 환경에서 진행되나, 실세 서버 환경에서도 거의 동일하게 적용할 수 있도록 작성되었다. 이번 장의 코드는 다음 링크에서 확인할 수 있다.

* 이번 장 코드 : [https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04)

이 문서에서 구성하는 인프라스트럭처는 다음과 같다.

![01](./01.png)

## Nginx와 설치

`Nginx`는 대표적인 웹 서버 중 하나로, 가볍고 높은 성능으로 많은 엔지니어들의 사랑(?)을 받고 있다. 상용 솔루션 뿐 아니라 오픈 소스조차 굉장히 성능이 우수하고, 필요 기능은 공개된 모듈을 통해서 쉽게 커스텀이 가능하기 떄문에 업계 표준으로 자리잡았다.

로컬 환경에서는 다음과 같이 `Docker`로 간단하게 설치 및 구동 가능하다.

```bash
$ docker run --rm -p 8080:80 nginx
```

역시 이 장의 코드를 다운 받았다면, 다음과 같이 `docker-compose`로 간단하게 설치 및 구동할 수 있다.

```bash
$ pwd
/Users/gurumee/Workspace/gurumee-book-prometheus/src/part2/ch04

$ docker compose up -d nginx
[+] Running 2/2
 ⠿ Network ch04_default  Created                                                                                                                                                                                                   0.3s
 ⠿ Container nginx       Started  
```

서버 환경에서는 다음 명령어로 설치 및 구동이 가능하다.

```bash
# 필요 패키지 설치
$ sudo yum install -y yum-utils

# nginx 패키지 레포지토리 추가
$ sudo tee /etc/yum.repos.d/nginx.repo << EOF
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/\$releasever/\$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true

[nginx-mainline]
name=nginx mainline repo
baseurl=http://nginx.org/packages/mainline/centos/\$releasever/\$basearch/
gpgcheck=1
enabled=0
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
EOF

# nginx 레포지토리 선택
$ sudo yum-config-manager --enable nginx-stable

# nginx 설치
$ sudo yum install -y nginx 

# nginx 구동
$ sudo systemctl restart nginx 

# nginx 구동 상태 확인
$ sudo systemctl status nginx
● nginx.service - nginx - high performance web server
   Loaded: loaded (/usr/lib/systemd/system/nginx.service; disabled; vendor preset: dis>
   Active: active (running) since Thu 2021-07-22 02:20:48 UTC; 4s ago
     Docs: http://nginx.org/en/docs/
  Process: 2037 ExecStart=/usr/sbin/nginx -c /etc/nginx/nginx.conf (code=exited, statu>
...
```

그 후 터미널에 다음을 입력하면 다음 결과를 얻을 수 있다.

```bash
# 로컬의 경우
$ curl localhost:8080

# 서버의 경우
$ curl localhost

# 결과 출력
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

## prometheus-nginxlog-exporter와 설치

RPS, 상태 코드 개수, Fail Rate, Response Time 등의 서비스 메트릭은 대부분 액세스 로그의 그 정보가 들어 있다. 아쉽게도 `Prometheus`에서 공식적으로 지원하는 `Nginx` 액세스 로그르 파싱해서 메트릭을 수집하는 `Exporter`는 없다. 

하지만, 제 3자에 의해서 `Nginx` 로그를 파싱하는 `Exporter`를 만들어서 공개했는데 그게 바로 `prometheus-nginxlog-exporter`이다. 터미널에 다음을 입력해서 설치 및 구동이 가능하다. 

```bash
$ docker run --rm quay.io/martinhelmich/prometheus-nginxlog-exporter
```  

역시 이 장의 코드를 다운 받았다면, 다음과 같이 `docker-compose`로 간단하게 설치 및 구동할 수 있다.
   
```bash
$ pwd
/Users/gurumee/Workspace/gurumee-book-prometheus/src/part2/ch04

$ docker compose up -d prometheus-nginxlog-exporter
[+] Running 2/2
⠿ Network ch04_default  Created                                                                                                                                                                                                   0.3s
⠿ Container nginx       Started  
```

서버 환경에서는 다음과 같이 설치 및 구동할 수 있다.

```bash
$ pwd
/home/sidelineowl

$ mkdir -p ~/apps/prometheus-nginxlog-exporter

# 압축 파일 다운로드
$ wget https://github.com/martin-helmich/prometheus-nginxlog-exporter/releases/download/v1.8.0/prometheus-nginxlog-exporter_1.8.0_linux_amd64.tar.gz

# 압축 파일 해제
$ tar -xvf prometheus-nginxlog-exporter_1.8.0_linux_amd64.tar.gz -C ~/apps/prometheus-nginxlog-exporter

# 압축 파일 삭제
$ rm prometheus-nginxlog-exporter_1.8.0_linux_amd64.tar.gz 

# prometheus-nginxlog-exporter 설치된 디렉토리로 이동
$ cd apps/prometheus-nginxlog-exporter/

# prometheus-nginxlog-exporter 실행
$ ./prometheus-nginxlog-exporter 
...
```

이제 손쉽게 구동하기 위해서 서비스로 등록해보자.

```bash
$ pwd
/home/sidelineowl/apps/prometheus-nginxlog-exporter

# user 추가
$ sudo useradd -M -r -s /bin/false prometheus_nginxlog_exporter

# 실행 파일 /usr/local/bin/으로 경로 이동
$ sudo cp ./prometheus-nginxlog-exporter /usr/local/bin

# /usr/local/bin/prometheus-nginxlog-exporter prometheus_nginxlog_exporter 유저, 그룹 권한 주기
$ sudo chown prometheus_nginxlog_exporter:prometheus_nginxlog_exporter /usr/local/bin/prometheus-nginxlog-exporter

# 서비스 파일 등록
$ sudo tee /etc/systemd/system/prometheus_nginxlog_exporter.service << EOF
[Unit]
Description=Prometheus Nginxlog Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus_nginxlog_exporter
Group=prometheus_nginxlog_exporter
Type=simple
ExecStart=/usr/local/bin/prometheus-nginxlog-exporter

[Install]
WantedBy=multi-user.target
EOF

# 데몬 리로드
# sudo systemctl daemon-reload
```

그 후 터미널에 다음을 입력해서 서비스를 구동시킨다.

```bash
# 서비스 가동
$ sudo systemctl restart prometheus_nginxlog_exporter

# 서비스 상태 확인
$ sudo systemctl status prometheus_nginxlog_exporter

● prometheus_nginxlog_exporter.service - Prometheus Nginxlog Exporter
   Loaded: loaded (/etc/systemd/system/prometheus_nginxlog_exporter.service; disabled;>
   Active: active (running) since Thu 2021-07-22 02:42:43 UTC; 5s ago
 Main PID: 2160 (prometheus-ngin)
    Tasks: 6 (limit: 23679)
   Memory: 2.4M
   CGroup: /system.slice/prometheus_nginxlog_exporter.service
           └─2160 /usr/local/bin/prometheus-nginxlog-exporter
...
```

그 후 터미널에 다음을 입력하면 다음 결과를 얻을 수 있다.

```bash
$ curl localhost:4040/metrics
# HELP nginx_parse_errors_total Total number of log file lines that could not be parsed
# TYPE nginx_parse_errors_total counter
nginx_parse_errors_total 0
```

## 메트릭 수집을 위한 각 컴포넌트 설정

이제 각 컴포넌트를 설정해서 메트릭을 수집해보자. 로컬의 경우엔, 이미 다 설정되어 있으므로 어떻게 설정하는지만 확인하면 된다. 먼저 `Nginx`부터 설정해보자. 기본적으로 `Nginx` 액세스 로그 포맷은 `Response Time`이 들어있지 않다. 관련 메트릭 수집을 위해서 액세스 로그를 다음과 같이 설정한다. 서버 환경에서는 `/etc/nginx/nginx.conf`에서 작업하면 된다.

[src/part2/ch04/nginx/nginx.conf](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/nginx/nginx.conf)
```conf
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    
    # 이 부분을 수정
    log_format  main  '$remote_addr - $remote_user [$time_local] '
                      '"$request" $status $body_bytes_sent '
                      '"$http_referer" "$http_user_agent" "$request_time"';


    access_log  /var/log/nginx/access.log  main;
    sendfile        on;
    keepalive_timeout  65;
    include /etc/nginx/conf.d/*.conf;
}
```

터미널에 다음을 입력해서 `Nginx`를 재구동한다.

```bash
$ sudo systemctl restart nginx
```

이제 `prometheus-nginxlog-exporter`가 `Nginx`의 액세스 로그를 파싱하기 위해서 해당 파일에 모든 접근을 허용하는 권한을 준다.

```bash
$ sudo chmod 777 /var/log/nginx/access.log
```

그 후 `prometheus-nginxlog-exporter` 설정을 진행한다. `/etc/` 경로에 `prometheus-nginxlog-exporter.yml`을 다음과 같이 설정한다.

[src/part2/ch04/prometheus-nginxlog-exporter/prometheus-nginxlog-exporter.yml](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/prometheus-nginxlog-exporter/prometheus-nginxlog-exporter.yml)
```yml
listen:
  port: 4040
  metrics_endpoint: "/metrics"

consul:
  enable: false

namespaces:
  - name: nginx
    format: "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$request_time\""
    source:
      files:
        - /var/log/nginx/access.log
    only_count: true
    relabel_configs:
    - target_label: request_uri
      from: request
      split: 2
      separator: ' '  
```

`format`은 `Nginx` 설정에서의 `log_format`을 따른다. 그 후 파일 경로를 주면 된다. `relabel_configs`는 메트릭 수집할 때 기본 `Label`이 있는데, 이를 조정하기 위해서 썼다. 자세한 내용은 다음을 참조하면 된다.

* [prometheus-nginxlog-exporter README](https://github.com/martin-helmich/prometheus-nginxlog-exporter)

이제 `prometheus-nginxlog-exporter` 서비스를 다음과 같이 수정한다.

/etc/systemd/system/prometheus_nginxlog_exporter.service
``` 
[Unit]
Description=Prometheus Nginxlog Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus_nginxlog_exporter
Group=prometheus_nginxlog_exporter
Type=simple

# -config-file 수정
ExecStart=/usr/local/bin/prometheus-nginxlog-exporter -config-file /etc/prometheus-nginxlog-exporter.yml

[Install]
WantedBy=multi-user.target
```

그 후 `prometheus-nginxlog-exporter` 서비스를 재구동한다.

```bash
# 데몬 리로드
$ sudo systemctl daemon-reload

# 서비스 재구동
$ sudo systemctl restart prometheus_nginxlog_exporter
```

그 다음 `Nginx`에 여러 번 `curl`을 날려본다.

```bash
# curl 5번 요청
$ for i in {1..5}
do
curl localhost
done
```

그 후 터미널에 다음을 입력한다.

```bash
$ curl localhost:4040/metrics
# HELP nginx_http_response_count_total Amount of processed HTTP requests
# TYPE nginx_http_response_count_total counter
nginx_http_response_count_total{method="GET",request_uri="/",status="200"} 5
# HELP nginx_http_response_size_bytes Total amount of transferred bytes
# TYPE nginx_http_response_size_bytes counter
nginx_http_response_size_bytes{method="GET",request_uri="/",status="200"} 3060
# HELP nginx_http_response_time_seconds Time needed by NGINX to handle requests
# TYPE nginx_http_response_time_seconds summary
nginx_http_response_time_seconds{method="GET",request_uri="/",status="200",quantile="0.
5"} 0
nginx_http_response_time_seconds{method="GET",request_uri="/",status="200",quantile="0.
9"} 0
nginx_http_response_time_seconds{method="GET",request_uri="/",status="200",quantile="0.
99"} 0
nginx_http_response_time_seconds_sum{method="GET",request_uri="/",status="200"} 0
nginx_http_response_time_seconds_count{method="GET",request_uri="/",status="200"} 5
# HELP nginx_http_response_time_seconds_hist Time needed by NGINX to handle requests
# TYPE nginx_http_response_time_seconds_hist histogram
nginx_http_response_time_seconds_hist_bucket{method="GET",request_uri="/",status="200",
le="0.005"} 5
...
```

`nginx_*`으로 시작하는 메트릭들이 수집된다면 성공이다. 이제 마지막으로 `Prometheus`를 다음과 같이 설정한다. (`/etc/prometheus/prometheus.yml`)

[src/part2/ch04/prometheus/prometheus.yml](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/prometheus/prometheus.yml)
```yml
# my global config
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # By default, scrape targets every 15 seconds.
  
  external_labels:
    monitor: 'my-project'

rule_files:

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'prometheus-nginxlog-exporter'
    scrape_interval: 5s

    static_configs:
      # Nginx와 prometheus-nginxlog-exporter 가 설치된 IP:4040
      - targets: ['prometheus-nginxlog-exporter:4040']
```

이제 `Prometheus`를 재구동한다.

```bash
$ sudo systemctl restart prometheus
```

그 후, `Prometheus UI`에서 다음을 쿼리해보자.

```
nginx_http_response_count_total
```

![02](./02.png)

## 서비스 메트릭 모니터링을 위한 Grafana 대시보드 구축
## 개요

이 문서에서는 `prometheus-nginxlog-exporter`를 이용해서 `Nginx`의 액세스 로그를 분석하여 RPS, 상태 코드 개수 등의 서비스 메트릭을 수집한다. 그 후 `Grafana`, `Prometheus`를 이용해서 시스템 메트릭을 모니터링할 수 있는 대시보드를 구축하는 것에 대하여 다룬다. 자세한 내용은 다음과 같다.

* Nginx와 설치
* prometheus-nginxlog-exporter와 설치
* 메트릭 수집을 위한 각 컴포넌트 설정
* 서비스 메트릭 모니터링을 위한 Grafana 대시보드 구축

이 문서에서 진행되는 실습 코드는 편의성을 위해 로컬 `Docker` 환경에서 진행되나, 실세 서버 환경에서도 거의 동일하게 적용할 수 있도록 작성되었다. 이번 장의 코드는 다음 링크에서 확인할 수 있다.

* 이번 장 코드 : [https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04)

이 문서에서 구성하는 인프라스트럭처는 다음과 같다.

![01](./01.png)

## Nginx와 설치

`Nginx`는 대표적인 웹 서버 중 하나로, 가볍고 높은 성능으로 많은 엔지니어들의 사랑(?)을 받고 있다. 상용 솔루션 뿐 아니라 오픈 소스조차 굉장히 성능이 우수하고, 필요 기능은 공개된 모듈을 통해서 쉽게 커스텀이 가능하기 떄문에 업계 표준으로 자리잡았다.

로컬 환경에서는 다음과 같이 `Docker`로 간단하게 설치 및 구동 가능하다.

```bash
$ docker run --rm -p 8080:80 nginx
```

역시 이 장의 코드를 다운 받았다면, 다음과 같이 `docker-compose`로 간단하게 설치 및 구동할 수 있다.

```bash
$ pwd
/Users/gurumee/Workspace/gurumee-book-prometheus/src/part2/ch04

$ docker compose up -d nginx
[+] Running 2/2
 ⠿ Network ch04_default  Created                                                                                                                                                                                                   0.3s
 ⠿ Container nginx       Started  
```

서버 환경에서는 다음 명령어로 설치 및 구동이 가능하다.

```bash
# 필요 패키지 설치
$ sudo yum install -y yum-utils

# nginx 패키지 레포지토리 추가
$ sudo tee /etc/yum.repos.d/nginx.repo << EOF
[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/\$releasever/\$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true

[nginx-mainline]
name=nginx mainline repo
baseurl=http://nginx.org/packages/mainline/centos/\$releasever/\$basearch/
gpgcheck=1
enabled=0
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true
EOF

# nginx 레포지토리 선택
$ sudo yum-config-manager --enable nginx-stable

# nginx 설치
$ sudo yum install -y nginx 

# nginx 구동
$ sudo systemctl restart nginx 

# nginx 구동 상태 확인
$ sudo systemctl status nginx
● nginx.service - nginx - high performance web server
   Loaded: loaded (/usr/lib/systemd/system/nginx.service; disabled; vendor preset: dis>
   Active: active (running) since Thu 2021-07-22 02:20:48 UTC; 4s ago
     Docs: http://nginx.org/en/docs/
  Process: 2037 ExecStart=/usr/sbin/nginx -c /etc/nginx/nginx.conf (code=exited, statu>
...
```

그 후 터미널에 다음을 입력하면 다음 결과를 얻을 수 있다.

```bash
# 로컬의 경우
$ curl localhost:8080

# 서버의 경우
$ curl localhost

# 결과 출력
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

## prometheus-nginxlog-exporter와 설치

RPS, 상태 코드 개수, Fail Rate, Response Time 등의 서비스 메트릭은 대부분 액세스 로그의 그 정보가 들어 있다. 아쉽게도 `Prometheus`에서 공식적으로 지원하는 `Nginx` 액세스 로그르 파싱해서 메트릭을 수집하는 `Exporter`는 없다. 

하지만, 제 3자에 의해서 `Nginx` 로그를 파싱하는 `Exporter`를 만들어서 공개했는데 그게 바로 `prometheus-nginxlog-exporter`이다. 터미널에 다음을 입력해서 설치 및 구동이 가능하다. 

```bash
$ docker run --rm quay.io/martinhelmich/prometheus-nginxlog-exporter
```  

역시 이 장의 코드를 다운 받았다면, 다음과 같이 `docker-compose`로 간단하게 설치 및 구동할 수 있다.
   
```bash
$ pwd
/Users/gurumee/Workspace/gurumee-book-prometheus/src/part2/ch04

$ docker compose up -d prometheus-nginxlog-exporter
[+] Running 2/2
⠿ Network ch04_default  Created                                                                                                                                                                                                   0.3s
⠿ Container nginx       Started  
```

서버 환경에서는 다음과 같이 설치 및 구동할 수 있다.

```bash
$ pwd
/home/sidelineowl

$ mkdir -p ~/apps/prometheus-nginxlog-exporter

# 압축 파일 다운로드
$ wget https://github.com/martin-helmich/prometheus-nginxlog-exporter/releases/download/v1.8.0/prometheus-nginxlog-exporter_1.8.0_linux_amd64.tar.gz

# 압축 파일 해제
$ tar -xvf prometheus-nginxlog-exporter_1.8.0_linux_amd64.tar.gz -C ~/apps/prometheus-nginxlog-exporter

# 압축 파일 삭제
$ rm prometheus-nginxlog-exporter_1.8.0_linux_amd64.tar.gz 

# prometheus-nginxlog-exporter 설치된 디렉토리로 이동
$ cd apps/prometheus-nginxlog-exporter/

# prometheus-nginxlog-exporter 실행
$ ./prometheus-nginxlog-exporter 
...
```

이제 손쉽게 구동하기 위해서 서비스로 등록해보자.

```bash
$ pwd
/home/sidelineowl/apps/prometheus-nginxlog-exporter

# user 추가
$ sudo useradd -M -r -s /bin/false prometheus_nginxlog_exporter

# 실행 파일 /usr/local/bin/으로 경로 이동
$ sudo cp ./prometheus-nginxlog-exporter /usr/local/bin

# /usr/local/bin/prometheus-nginxlog-exporter prometheus_nginxlog_exporter 유저, 그룹 권한 주기
$ sudo chown prometheus_nginxlog_exporter:prometheus_nginxlog_exporter /usr/local/bin/prometheus-nginxlog-exporter

# 서비스 파일 등록
$ sudo tee /etc/systemd/system/prometheus_nginxlog_exporter.service << EOF
[Unit]
Description=Prometheus Nginxlog Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus_nginxlog_exporter
Group=prometheus_nginxlog_exporter
Type=simple
ExecStart=/usr/local/bin/prometheus-nginxlog-exporter

[Install]
WantedBy=multi-user.target
EOF

# 데몬 리로드
# sudo systemctl daemon-reload
```

```bash
# 서비스 가동
$ sudo systemctl restart prometheus_nginxlog_exporter

# 서비스 상태 확인
$ sudo systemctl status prometheus_nginxlog_exporter

● prometheus_nginxlog_exporter.service - Prometheus Nginxlog Exporter
   Loaded: loaded (/etc/systemd/system/prometheus_nginxlog_exporter.service; disabled;>
   Active: active (running) since Thu 2021-07-22 02:42:43 UTC; 5s ago
 Main PID: 2160 (prometheus-ngin)
    Tasks: 6 (limit: 23679)
   Memory: 2.4M
   CGroup: /system.slice/prometheus_nginxlog_exporter.service
           └─2160 /usr/local/bin/prometheus-nginxlog-exporter
...
```

그 후 터미널에 다음을 입력하면 다음 결과를 얻을 수 있다.

```bash
$ curl localhost:4040/metrics
# HELP nginx_parse_errors_total Total number of log file lines that could not be parsed
# TYPE nginx_parse_errors_total counter
nginx_parse_errors_total 0
```

## 메트릭 수집을 위한 각 컴포넌트 설정

이제 각 컴포넌트를 설정해서 메트릭을 수집해보자. 로컬의 경우엔, 이미 다 설정되어 있으므로 어떻게 설정하는지만 확인하면 된다. 먼저 `Nginx`부터 설정해보자. 기본적으로 `Nginx` 액세스 로그 포맷은 `Response Time`이 들어있지 않다. 관련 메트릭 수집을 위해서 액세스 로그를 다음과 같이 설정한다. 서버 환경에서는 `/etc/nginx/nginx.conf`에서 작업하면 된다.

[src/part2/ch04/nginx/nginx.conf](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/nginx/nginx.conf)
```conf
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    
    # 이 부분을 수정
    log_format  main  '$remote_addr - $remote_user [$time_local] '
                      '"$request" $status $body_bytes_sent '
                      '"$http_referer" "$http_user_agent" "$request_time"';


    access_log  /var/log/nginx/access.log  main;
    sendfile        on;
    keepalive_timeout  65;
    include /etc/nginx/conf.d/*.conf;
}
```

터미널에 다음을 입력해서 `Nginx`를 재구동한다.

```bash
$ sudo systemctl restart nginx
```

이제 `prometheus-nginxlog-exporter`가 `Nginx`의 액세스 로그를 파싱하기 위해서 해당 파일에 모든 접근을 허용하는 권한을 준다.

```bash
$ sudo chmod 777 /var/log/nginx/access.log
```

그 후 `prometheus-nginxlog-exporter` 설정을 진행한다. `/etc/` 경로에 `prometheus-nginxlog-exporter.yml`을 다음과 같이 설정한다.

[src/part2/ch04/prometheus-nginxlog-exporter/prometheus-nginxlog-exporter.yml](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/prometheus-nginxlog-exporter/prometheus-nginxlog-exporter.yml)
```yml
listen:
  port: 4040
  metrics_endpoint: "/metrics"

consul:
  enable: false

namespaces:
  - name: nginx
    format: "$remote_addr - $remote_user [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\" \"$request_time\""
    source:
      files:
        - /var/log/nginx/access.log
    only_count: true
    relabel_configs:
    - target_label: request_uri
      from: request
      split: 2
      separator: ' '  
```

`format`은 `Nginx` 설정에서의 `log_format`을 따른다. 그 후 파일 경로를 주면 된다. `relabel_configs`는 메트릭 수집할 때 기본 `Label`이 있는데, 이를 조정하기 위해서 썼다. 자세한 내용은 다음을 참조하면 된다.

* [prometheus-nginxlog-exporter README](https://github.com/martin-helmich/prometheus-nginxlog-exporter)

이제 `prometheus-nginxlog-exporter` 서비스를 다음과 같이 수정한다.

/etc/systemd/system/prometheus_nginxlog_exporter.service
``` 
[Unit]
Description=Prometheus Nginxlog Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus_nginxlog_exporter
Group=prometheus_nginxlog_exporter
Type=simple

# -config-file 수정
ExecStart=/usr/local/bin/prometheus-nginxlog-exporter -config-file /etc/prometheus-nginxlog-exporter.yml

[Install]
WantedBy=multi-user.target
```

그 후 `prometheus-nginxlog-exporter` 서비스를 재구동한다.

```bash
# 데몬 리로드
$ sudo systemctl daemon-reload

# 서비스 재구동
$ sudo systemctl restart prometheus_nginxlog_exporter
```

그 다음 `Nginx`에 여러 번 `curl`을 날려본다.

```bash
# curl 5번 요청
$ for i in {1..5}
do
curl localhost
done
```

그 후 터미널에 다음을 입력한다.

```bash
$ curl localhost:4040/metrics
# HELP nginx_http_response_count_total Amount of processed HTTP requests
# TYPE nginx_http_response_count_total counter
nginx_http_response_count_total{method="GET",request_uri="/",status="200"} 5
# HELP nginx_http_response_size_bytes Total amount of transferred bytes
# TYPE nginx_http_response_size_bytes counter
nginx_http_response_size_bytes{method="GET",request_uri="/",status="200"} 3060
# HELP nginx_http_response_time_seconds Time needed by NGINX to handle requests
# TYPE nginx_http_response_time_seconds summary
nginx_http_response_time_seconds{method="GET",request_uri="/",status="200",quantile="0.
5"} 0
nginx_http_response_time_seconds{method="GET",request_uri="/",status="200",quantile="0.
9"} 0
nginx_http_response_time_seconds{method="GET",request_uri="/",status="200",quantile="0.
99"} 0
nginx_http_response_time_seconds_sum{method="GET",request_uri="/",status="200"} 0
nginx_http_response_time_seconds_count{method="GET",request_uri="/",status="200"} 5
# HELP nginx_http_response_time_seconds_hist Time needed by NGINX to handle requests
# TYPE nginx_http_response_time_seconds_hist histogram
nginx_http_response_time_seconds_hist_bucket{method="GET",request_uri="/",status="200",
le="0.005"} 5
...
```

`nginx_*`으로 시작하는 메트릭들이 수집된다면 성공이다. 이제 마지막으로 `Prometheus`를 다음과 같이 설정한다. (`/etc/prometheus/prometheus.yml`)

[src/part2/ch04/prometheus/prometheus.yml](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/prometheus/prometheus.yml)
```yml
# my global config
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # By default, scrape targets every 15 seconds.
  
  external_labels:
    monitor: 'my-project'

rule_files:

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'prometheus-nginxlog-exporter'
    scrape_interval: 5s

    static_configs:
      # Nginx와 prometheus-nginxlog-exporter 가 설치된 IP:4040
      - targets: ['prometheus-nginxlog-exporter:4040']
```

이제 `Prometheus`를 재구동한다.

```bash
$ sudo systemctl restart prometheus
```

그 후, `Prometheus UI`에서 다음을 쿼리해보자.

```
nginx_http_response_count_total
```

![02](./02.png)

## 서비스 메트릭 모니터링을 위한 Grafana 대시보드 구축

이제 대시보드를 구축한다. 다음 JSON 파일을 복사해서 대시보드를 임포트한다. (로컬 환경에는 이미 대시보드가 로드되어 있다.)

* [src/part2/ch04/grafana/dashboard.json](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/grafana/dashboard.json)

그럼 다음과 같은 대시보드를 확인할 수 있다.

![03](./03.png)
이제 대시보드를 구축한다. 다음 JSON 파일을 복사해서 대시보드를 임포트한다. (로컬 환경에는 이미 대시보드가 로드되어 있다.)

* [src/part2/ch04/grafana/dashboard.json](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch04/grafana/dashboard.json)

그럼 다음과 같은 대시보드를 확인할 수 있다.

![03](./03.png)