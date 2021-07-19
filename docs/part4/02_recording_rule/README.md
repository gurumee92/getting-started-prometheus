# 기록 규칙(Recording Rule)이란 무엇인가

![logo](../../logo.png)

## 개요

이번 장에서는 `PromQL`의 쿼리 성능을 보다 높여줄 수 있는 `Recording Rule`에 대해 알아볼 것이다. 이 문서에서는 편의성을 위해서 `Docker` 환경에서 진행할 것이나, 실제 서버 환경에서 어떻게 작업해야 하는지까지 최대한 다루도록 하겠다. 관련 코드는 다음 링크를 참고하길 바란다.

* 이번 장 코드 : [https://github.com/gurumee92/gurumee-prometheus-code/tree/master/part4/ch02](https://github.com/gurumee92/gurumee-prometheus-code/tree/master/part4/ch02)

## Recording Rule은 무엇인가?

`Prometheus` 공식 문서에 따르면, `Recording Rule`의 정의는 다음과 같다.

> 기록 규칙은 일관된 이름 체계를 사용함으로써, 한 눈에 규칙을 쉽게 파악할 수 있다. 이것은 또한 부정확하거나 무의미한 계산을 두드러지게 함으로써 실수를 피한다.

개인적으로는 다소 뜬금 없는 표현이라고 생각되는데 쉽게 생각해서, `PromQL`을 통해서 집계한 시계열 데이터의 이름을 붙인다고 생각하면 된다. 예를 들면 지난 장 "[서비스 메트릭 모니터링하기 (1) prometheus-nginxlog-exporter]()"에서 `Nginx`의 1분간 평균 응답 시간에 대한 패널을 어떻게 구축했는가?

![grafana > response time](./01.png)

해당 패널은 다음 쿼리로 구성된다.

```
sum by (method, request_uri)(rate(nginx_http_response_time_seconds_hist_sum[1m])) 
/ 
sum by (method, request_uri)(rate(nginx_http_response_time_seconds_hist_count[1m]))
```

위의 쿼리는 단순한 수준에 속하지만 `PromQL`을 잘 모르는 경우에는 무슨 쿼리인지 도통 알아보기가 힘들다. 쿼리가 복잡해지면 복잡해질수록, 해석하기가 점점 히들어질 것이다. 이를 위해서 `Recording Rule`을 사용하면, 다음과 같이 단순하게 표현 가능하다.

```
job:nginx_http_response_time_seconds:avg_1m
```

다음과 같이 쿼리하게 되면 적어도, 패널이 쿼리하는 데이터가 1분간 집계된 평균 응답 시간이라는 것을 쉽게 알 수 있다. 사실은 이런 효과보다도, `InfluxDB`의 `Continuous Query` 같이, 복잡한 집계 쿼리를 배치를 돌듯 쿼리 결과를 하나의 시계열 데이터로 저장하여 쿼리 성능을 높이는 데 더 큰 의미가 있다고 생각한다.

## Recording Rule 사용법 (1) 작성 방법과 관례

`Recording Rule`을 만들기 위해서는 다음과 같은 절차가 필요하다.

1. Recording Rule 작성 (보통 <job_name>_rules.yml 파일에 작성한다.)
2. prometheus.yml에서 `rule_files` 설정 (<job_name>_rules.yml 상대 경로를 지정)

먼저 `Recording Rule`을 만들어보자. `prometheus.yml`이 저장된 디렉토리 경로에서 `rules`라는 디렉토리를 만들고 `prometheus_nginxlog_exporter_rules.yml`을 다음과 같이 작성한다.

[part4/ch02/prometheus/rules/prometheus_nginxlog_exporter_rules.yml](https://github.com/gurumee92/gurumee-prometheus-code/blob/master/part4/ch02/prometheus/rules/prometheus_nginx_log_exporter_rules.yml)
```yml
groups:
- name: prometheus_nginxlog_exporter
  rules:
  - record: job:nginx_http_response_time_seconds:avg_1m
    expr: |
      sum by (method, request_uri)(rate(nginx_http_response_time_seconds_hist_sum[1m])) 
      / 
      sum by (method, request_uri)(rate(nginx_http_response_time_seconds_hist_count[1m]))

  # ...
```

공식 문서에 따르면 다음의 네이밍 규칙을 권장하고 있다.

```
<level>:<metric>:<operations>
```

개인적으로는 "level"이 제일 이해가 안갔었는데, 보통 집계 수준을 정의한다. 쉽게 설명하자면, 현재처럼 method, request_uri 별이 아닌 method 별 응답 시간이었다면 `Recording Rule`은 다음과 같이 작성할 수 있다.

```yml
# ...
  rules:
  - record: method:nginx_http_response_time_seconds:avg_1m
    expr: |
      sum by (method)(rate(nginx_http_response_time_seconds_hist_sum[1m])) 
      / 
      sum by (method)(rate(nginx_http_response_time_seconds_hist_count[1m]))
```

어떤 기준으로 집계하는지에 따라, level을 달리 작성할 수 있을 것이다. 집계 수준이 2개 이상인 경우는 "job"으로 표현한다. "metric"은 실제 메트릭의 이름을 지정하면 된다. `Histogram`, `Summary` 메트릭 타입은 `nginx_http_response_time_seconds_hist_sum`처럼 메트릭 이름 이후에 `postfix`가 붙는다. 이때는 `postfix`를 제거한다. 그래서 `nginx_http_response_time_seconds`으로 지정한 것이다. "operations"는 집계 연산을 나타내면 된다. 위 `Recording Rule`은 1분간 평균이니까 `avg_1m`으로 지정했다.

`<job_name>_rules.yml`을 작성했으면, `prometheus.yml`에서 다음과 같이 설정하면 된다.

[part4/ch02/prometheus/prometheus.yml](https://github.com/gurumee92/gurumee-prometheus-code/blob/master/part4/ch02/prometheus/prometheus.yml)
```yml
# my global config
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # By default, scrape targets every 15 seconds.

rule_files:
  - 'rules/prometheus_nginx_log_exporter_rules.yml' # prometheus.yml 상대 경로를 지정한다.

# ...
```

여기서 중요한 점은 `rule_files`에서 `<job_name>_rules.yml`의 경로를 지정해야 하는데, 이는 `prometheus.yml`에서 상대 경로로 지정된다. 보통 다음과 같은 위치를 가진다.

```
|- /home/ec2-user/apps/prometheus (prometheus 디렉토리)
    |- prometheus.yml
    |- rules
        |- prometheus_nginxlog_exporter_rules.yml
        |- ...
```

`prometheus.yml`을 수정한 후 `Prometheus`를 재시작 혹은 설정 리로드를 해주면 된다. 재시작은 다음 명령어로 실행할 수 있다.

```bash
$ sudo systemctl restart prometheus
```

설정 리로드는 "HUP" 시그널을 통해서 할 수 있다. 먼저 `Prometheus`의 "process id"를 알아낸 후 "kill -HUP" 명령어를 주면 된다.

```bash
# prometheus PID 알기
$ ps -ef | grep prometheus
# ...
root 17625 1 2 6월03 ? 09:53:10 /home/ec2-user/apps/prometheus/prometheus --config.file=/home/ec2-user/apps /prometheus/prometheus.yml --storage.tsdb.path=/home/ec2-user/apps/prometheus/data

# 설정 리로드
# sudo kill -HUP <pid> 
$ sudo kill -HUP 17625
```

이렇게 설정이 완료되면 `Prometheus` 웹 UI 상단 메뉴 "Status"의 "Rules"에서 작성된 `Recording Rule`을 확인할 수 있다.

![Prometheus Web UI > Status > Rules](./02.png)

위 화면에서 알 수 있듯이, `Recording Rule`의 이름, 쓰였던 집계 쿼리는 무엇인지는 물론 사용 가능한 상태인지, 언제 집계 데이터를 저장했는지, 저장할 때 얼마나 시간이 걸렸는지까지 확인할 수 있다.

## Recording Rule 사용법 (2) 권장되는 상황

`Recording Rule`이 권장되는 상황은 크게 다음과 같다.

1. PromQL 연산이 복잡한 경우 (이전 절 "Recording Rule 사용법 (1) 작성 방법과 관례"에서 다루었다.)
2. PromQL 집계 결과의 성능을 올리고 싶을 때 (=카디널리티를 줄이고 싶을 때)
3. 함수 입력으로 Range Vector를 넣고 싶을 때 (거의 사용되진 않음.)

`Recording Rule`의 주된 사용처는 2번이다. `PromQL`을 통해서 집계된 여러 시계열의 개수를 줄여, 쿼리 성능을 높이고 싶을 때 이것을 사용하면 매우 좋다. 예를 들어서 다음 패널을 보자.

![grafana > status 4xx](./03.png)

위 패널은 `Nginx`에서 발생하는 상태 코드에 대한 정보를 나타낸다. 이 중 클라이언트 에러를 나타내는 4xx 코드를 집계하는 쿼리는 다음과 같다.

```
sum(rate(nginx_http_response_count_total{status=~"4.."}[1m]))
```

이 때 만약 여태까지 발생했던 상태 코드가 400, 401, 402, 403, 404가 있었다면, 위 쿼리는 총 5개의 시계열을 집계하고 있는 상황이다. 

![4xx time sereies 1](./04.png)

이를 `Recording Rule`을 이용해서, 쿼리 성능을 올려보자. 다음과 같이 작성할 수 있다.

[part4/ch02/prometheus/rules/prometheus_nginxlog_exporter_rules.yml](https://github.com/gurumee92/gurumee-prometheus-code/blob/master/part4/ch02/prometheus/rules/prometheus_nginx_log_exporter_rules.yml)
```yml
groups:
- name: prometheus_nginxlog_exporter
  rules:
  # ...
  - record: status:nginx_http_response_count_total:4xx_1m
    expr: sum(rate(nginx_http_response_count_total{status=~"4.."}[1m]))
  # ...
```

이 경우 다음과 같이 쿼리할 수 있다.

```
status:nginx_http_response_count_total:4xx_1m
```

이 둘의 쿼리 결과는 같다. 하지만, `Recording Rule`로 쿼리한 결과는 하나의 시계열을 가져온 것으로, `PromQL`로 집계한 결과보다 시계열이 총 4개 감소한다. 

![4xx time sereies 2](./05.png)

현재는 작은 규모의 `Prometheus`를 운영하는 중이니 큰 성능 향상은 없으나, 카디널리티가 늘어날 수록 모니터링 요소가 많아질 수록 쿼리 성능 최적화가 필요하게 된다. 이 때 `Recording Rule`이 좋은 해결책이 된다. 보통, `PromQL`에서 정규 표현식으로 결과를 집계하는 경우가 `Recording Rule`을 사용하라는 신호이니, 꼭 사용할 것을 검토해보길 바란다.

굉장히 드물게, `Range Vector`를 `max_over_time` 같은 함수의 입력으로 쓰고 싶을 수가 있다. 

```
max_over_time(sum(rate(nginx_http_response_count_total{status=~"4.."}[1m]))[1h])
```

아쉽게도 `Prometheus`의 `max_over_time` 등의 함수들은 `Instant Vector`를 입력으로 받는 상황이기 때문에 위의 `PromQL`은 실패한다. 

![max_over_time 실패](./06.png)

하지만, 이 때 입력을 `Recording Rule`을 통해서 `Instant Vector`로 변환할 수 있다. 따라서 다음 쿼리는 유효하다.

```
max_over_time(status:nginx_http_response_count_total:4xx_1m[1h])
```

위 쿼리를 `Prometheus` 웹 UI로 쿼리했을 때 결과는 다음과 같다. 유효한 쿼리이기 때문에 결과가 나온다.

![max_over_time 성공](./07.png)
 
## Recording Rule 사용법 (3) 권장되지 않는 상황

쿼리 결과의 카디널리티를 줄여서 성능을 올릴 때 `Recording Rule`은 아주 좋은 해결책이 된다. 하지만 만능 해결책은 아니다. 다음과 같은 상황에서 `Recording Rule` 사용을 권장하지 않는다.

1. 현재 수집된 모든 메트릭을 사전에 집계하는 경우
2. metric 이름과 label의 키를 수정해야 하는 경우

보통 메트릭 모니터링 시스템의 경우 90% 이상의 메트릭들이 사용되지 않기 때문에 모든 메트릭은 절대 사전에 집계해서는 안된다. 또한 metric 이름이나 lable의 키를 변경해야 하는 경우에도 `Recording Rule`은 좋은 해결책이 되지 않는다. `Recording Rule` 역시 비용이 들어가기 때문에 과하게 사용하는 것은 금물이다. 이런 상황의 경우, `relabel_configs`를 이용하는 것이 좋다.

그리고 보통 잘못되게 사용하는 경우가 있는데, 바로 label의 이점을 제거하는 경우이다. 일반적으로 `Recording Rule`을 작성할 때 metric의 특정 label에 대한 데이터를 뽑는 것은 지양해야 한다. 예를 들면 다음과 같다.

> 참고! 아래 예는 눈으로만 확인하세요
> 아래는 node-exporter가 수집하는 메트릭 기반으로 작성된 Recording Rule입니다. 눈으로만 왜 잘못된 것인지만 짚고 넘어가는 것을 추천드립니다. 직접 Recording Rule을 작성하고 그 결과를 확인하고 싶다면, 이전 장 "[시스템 메트릭 모니터링하기 (1) node-exporter]()"를 확인하세요.

```yml
# ...
  - record: instance:node_disk_read_bytes_total_nbd1:rate5m
    expr: sum without(instance)(rate(node_disk_read_bytes_total{device="nbd1"}[5m])) 
  - record: instance:node_disk_read_bytes_total_nbd2:rate5m
    expr: sum without(instance)(rate(node_disk_read_bytes_total{device="nbd2"}[5m]))
# ...
```

위의 작성된 `Recording rule`은 device 별로 `node_disk_read_bytes_total`를 집계하게 된다. 그러나 이 경우에는 실제로 별로 좋지 않다. 왜냐하면, `Recording Rule` 생성 비용이, 쿼리 비용보다 앞서기 때문이다. 무슨 말일까? 먼저 아래처럼 `Recording Rule`을 수정한다.

```yml
# ...
  - record: instance:node_disk_read_bytes_total:rate5m
    expr: sum without(instance)(rate(node_disk_read_bytes_total[5m])) 
# ...
```

그 후 쿼리 시에는 다음과 같이 쿼리하는 것이다.

```
instance:node_disk_read_bytes_total:rate5m{device="nbd1"}
```

전자의 경우 후자의 경우보다 2가지 관점에서 더 좋지 않다. 첫째로는 device는 실제로 훨씬 더 많을 수 있다. 근데 device의 값을 지정해서 `Recording Rule`로 만들게 되면 모든 device에 대한 `Recording Rule`을 작성하지 않는 한, 더 이상 다른 device들은 쿼리할 수가 없다. 

둘째로는 카디널리티를 줄이는 효과가 사라진다는 것이다. 첫번째 문제를 해결하기 위해서 특정 label 값마다 `Recording Rule`을 작성하게 되면, 시계열 수만큼, 즉 위 예제에서는 device 개수만큼 선형적으로 그 비용이 증가하게 된다. 이는 카디널리티를 줄이기 위해 사용되는 `Recording Rule`의 특성과 정 반대로 사용하는 것이다. 따라서 label 값마다 뽑아서 쓰는 `Recording Rule`은 절대 만들지 않는 것이 좋다. 

개인적으로는 다음과 같은 상황에 사용하는 것을 추천한다.

1. 집계 연산이 매우 복잡할 때
2. 정규 표현식을 매칭해서 쿼리할 때