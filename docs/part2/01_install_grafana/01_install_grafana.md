# Grafana 설치

이 장에서는 `Grafana` 설치를 다룬다. 역시 로컬 환경, `AWS` 클라우드 환경에서의 설치를 다룬다.

## Grafana란 무엇인가

먼저 설치 전에 왜 우리가 `Grafana`를 설치해야 하는지, `Grafana`가 무엇인지 알아두면 좋을 것 같다. `Grafana`란, `Grafana Labs`에서 관리하고 있는 오픈 소스 시각화 및 분석 도구이다. `Prometheus` 물론 `InfluxDB`, `Elasticsearch` 등 여러 데이터 소스와 통합이 가능하다.

물론 이전 장에서 잠깐 봤듯이 `Prometheus` 역시 자체적으로 UI를 제공하고 있다. 쿼리도 가능하고 심지어 여러 패널을 만들어서 대시보드 구성도 가능하다. 하지만, 그 기능이 너무나도 빈약해서 혹은 불편해서 보통 상용 환경에서는 `Grafana`와 함께 연동해서 사용하는 것이 일반적이다. `Grafana`는 여러 데이터 소스에 대한 대시보드 템플릿을 제공하기 때문에, `Prometheus` 등의 데이터 소스의 쿼리 방법을 잘 모른다 하더라도 기본적인 대시보드 구성이 가능하다.

기본적으로 `Prometheus`와 `Grafana`는 모두 `Grafana Labs`에서 관리하고 있기 때문에, 궁합이 어떤 데이터 소스와 비교하더라도 매우 좋은 편이다. 이것만으로 우리가 `Grafana`를 설치하는 이유는 충분하다.

## 로컬에서 Grafana 설치

먼저 로컬 환경에서 `Grafana`를 설치한다. 역시 `Docker` 기반으로 설치를 할 것이다. 터미널에 다음을 입력하면 바로 설치를 할 수 있다.

```bash
$ docker run --name=grafana -p 3000:3000 grafana/grafana
```

끝이다. 이후 이어지는 절에서, `Prometheus` 연동하는 방법을 알아보기 때문에, 로컬에서 진행할 사람들은 다음과 같이 `docker-compose.yml`을 만들어둔다.

part1/ch03/docker-compose.yml
```yml
version: "3"

services:
  prometheus:
    container_name: prometheus
    image: prom/prometheus:latest
    ports:
      - 9090:9090

  # dashboard & alert
  grafana:
    container_name: grafana
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
```

위 코드는 다음 링크에서 확인할 수 있다.

* [https://github.com/gurumee92/gurumee-prometheus-code/blob/master/part1/ch03/docker-compose.yml](https://github.com/gurumee92/gurumee-prometheus-code/blob/master/part1/ch03/docker-compose.yml)

로컬에서 `docker-compose`로 `Prometheus`와 `Grafana`를 조작하려면 터미널에 다음을 입력하면 된다.

```bash
# 현재 위치
$ pwd
# docker-compose.yml이 있는 위치
/Users/a1101320/Studies/gitbooks/gurumee-prometheus/code/part1/ch03

# 컴포넌트 실행
$ docker-compose up -d
Creating network "ch03_default" with the default driver
Creating prometheus ... done
Creating grafana    ... done

# 컴포넌트 상태 확인
$ docker ps
CONTAINER ID   IMAGE                    COMMAND                  CREATED          STATUS          PORTS                    NAMES
c2203aaa4464   prom/prometheus:latest   "/bin/prometheus --c…"   14 seconds ago   Up 13 seconds   0.0.0.0:9090->9090/tcp   prometheus
7fa2f3ecc913   grafana/grafana:latest   "/run.sh"                14 seconds ago   Up 13 seconds   0.0.0.0:3000->3000/tcp   grafana

# 컴포넌트 중지
$ docker-compose down -v
Stopping prometheus ... done
Stopping grafana    ... done
Removing prometheus ... done
Removing grafana    ... done
Removing network ch03_default
```

## AWS에서 Grafana 설치

지난 장에서 `Prometheus`처럼 직접 바이너리 설치하는 것이 아닌 rpm 파일을 다운 받아 로컬 패키지 매니저에 설치 명령을 내리는 것으로 `Grafana`를 설치해보겠다. 쉽게 말해서, 조금 더 쉽게 `Grafana`를 설치하고 서비스로 등록한다. **프로메테우스를 설치한 AWS 인스턴스의 터미널에 접속한 후** `Grafana` 설치를 위해 다음을 입력한다.

```bash
# .rpm install
$ wget https://dl.grafana.com/oss/release/grafana-7.3.6-1.x86_64.rpm
 
# rpm localinstall
$ sudo yum localinstall grafana-7.3.7-1.x86_64.rpm
```

그럼 자동으로 `Grafana`가 설치되고 `grafana-server`라는 이름으로 서비스가 등록된다. 이제 터미널에 다음을 입력해서 `Grafana`를 실행하면 된다.

```bash
$ sudo systemctl start grafana-server
 
# 그라파나 상태 확인
$ sudo systemctl status grafana-server
● grafana-server.service - Grafana instance
Loaded: loaded (/usr/lib/systemd/system/grafana-server.service; enabled; vendor preset: disabled)
Active: active (running) since 목 2021-01-14 07:04:46 UTC; 3 days ago
Docs: http://docs.grafana.org
....
```

> 참고!
> 
> 물론 Prometheus와 Grafana 각기 다른 서버에 설치해서 연동해도 됩니다. 다만 그에 따른 AWS 서버 설정이 필요하기 때문에 프로메테우스를 다루는 범위를 넘어간다고 생각합니다. 따라서 이에 대한 내용은 설명하지 않습니다. 이후에 진행되는 실습도 두 컴포넌트가 같은 서버에 설치되었다고 가정하고 진행할 것입니다. 실습을 따라하는 사람들은 최대한 같은 환경을 맞춰주시길 바랍니다.

## Prometheus와의 연동

이제 `Prometheus`와 `Grafana`를 연동해 볼 것이다. 로컬 기준으로 자세하게 알아볼 것이다. 서버 작업은 모두 동일하고 URL만 잘 설정해주면 된다. 

### 로컬에서 연동 

먼저 로컬 환경이다. `docker-compose`로 `Grafana`와 `Prometheus`를 실행한다.

```bash
# 현재 위치
$ pwd
# docker-compose.yml이 있는 위치
/Users/gurumee/gurumee-prometheus/code/part1/ch03

# 컴포넌트 실행
$ docker-compose up -d
```

그 후 브라우저에서 "localhost:3000"을 접속한다. 그럼 다음 화면이 뜨느데, Email과 Password 입력란에 "admin"을 입력한다.

* Email : admin
* Password : admin
  
![01](./01.png)

입력하고 나면 다음 화면이 뜨는데 좌측 하단에 "Skip"을 클릭한다. 물론 원하는 Email/Password가 있다면 입력하고 넘어가도 무방하다. 

> 참고!
> 
> 실제 서버 환경에서는 Email/Password를 보안을 위해서 재 설정하는 것이 좋습니다.

![02](./02.png)

그럼 다음 메인 UI로 이동하게 된다. 좌측 메뉴바에 6번째 메뉴 "톱니바퀴" 메뉴를 클릭하면, "Datasources"라는 메뉴가 보인다. 이를 클릭한다.

![03](./03.png)

그럼 다음 화면이 보이는데 "Add data sources"를 클릭한다.

![04](./04.png)

그럼 다음 화면이 보인다. 맨 첫 번째 보이는 `Prometheus`를 클릭한다. 오른쪽에 "Select"를 누르면 된다.

![05](./05.png)

그럼 다음 화면에서 `Prometheus` URL을 적어주면 된다. 현재는 도커 컨테이너 기반이니까 "컨테이너 이름:포트"로 접속이 가능하다. 따라서 현재는 "prometheus:9090"으로 접속하면 된다.

![06](./06.png)

그 후 아래에 "Save & Test"를 누르면 된다.

![07](./07.png)

성공적으로 연결되면 아래 화면처럼 성공했다는 문구가 보인다. 

![08](./08.png)

이제 다시 메인 메뉴로 돌아간다. 쿼리가 되는지 확인해보자. 좌측 탭의 4번째 "Explore"를 선택한다.

![09](./09.png)

쿼리 입력란에 "up"이라는 `Prometheus` 쿼리를 입력한다. 그 후 "Run Query"를 누른다.

> 참고!
> 
> up 쿼리는 프로메테우스가 수집하는 인스턴스 상태를 보여주는 쿼리입니다. 자세한 내용은 추후에 더 깊이 다뤄보도록 하겠습니다.

![10](./10.png)

그럼 아래 화면처럼 그래프와 테이블을 확인할 수 있다.

![11](./11.png)

### 서버에서 연동

만약 서버에서 직접 연동하는 경우는 `Grafana`와 `Prometheus` 서비스가 잘 실행되는지 확인하면 된다.

```bash
# 프로메테우스 상태 확인
$ sudo systemctl status prometheus

# 그라파나 상태 확인
$ sudo systemctl status grafana-server
● grafana-server.service - Grafana instance
Loaded: loaded (/usr/lib/systemd/system/grafana-server.service; enabled; vendor preset: disabled)
Active: active (running) since 목 2021-01-14 07:04:46 UTC; 3 days ago
Docs: http://docs.grafana.org
```

이후 진행 작업은 로컬 환경과 같다. 다만 `Prometheus` URL은 이번 실습 그대로 했다면 "localhost:9090"으로 하면 된다. 만약 다른 서버에서 작업했다면, "해당 서버 주소:9090"으로 설정하면 된다. 같은 VPC 내라면, 별 다른 작업 없이 private_ip, port로 작업이 가능하겠지만, 다른 VPC라면, public_ip 할당과, port에 대한 방화벽 설정이 필요하다.