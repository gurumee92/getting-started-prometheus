# 시스템 메트릭 모니터링하기 (1) node-exporter

## 개요

이 문서에서는 `Grafana`, `Prometheus`, `node-exporter`를 이용해서 시스템 메트릭을 모니터링할 수 있는 대시보드를 구축하는 것에 대하여 다룬다. 자세한 내용은 다음과 같다.

* node-exporter란 무엇인가
* node-exporter 설치
* node-exporter, Prometheus 연동
* 시스템 메트릭 모니터링을 위한 Grafana 대시보드 구축

이 문서에서 진행되는 실습 코드는 편의성을 위해 로컬 `Docker` 환경에서 진행되나, 실세 서버 환경에서도 거의 동일하게 적용할 수 있도록 작성되었다. 이번 장의 코드는 다음 링크에서 확인할 수 있다.

* 이번 장 코드 : [https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch02](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch02)

## node-exporter란 무엇인가

`node-exporter`란 UNIX 계열 커널을 가진 하드웨어와 OS릭 등 "시스템 메트릭"을 수집하는 `Exporter`이다. `Prometheus` 재단이 공식적으로 지원하고 있는 `Exporter` 중 하나이며 `Prometheus`로 모니터링 시스템을 구축 시 시스템 메트릭 수집을 위해 가장 우선적으로 고려되는 `Exporter`이기도 하다. 

> 참고! Exporter가 무엇인가요?
> 
> Exporter란 특정 메트릭을 수집해서 엔드포인트에 노출시키는 소프트웨어 혹은 에이전트라고 보시면 됩니다. node-exporter가 UNIX 계열 서버의 cpu, memory 등의 메트릭을 수집할 수 있는 것처럼, DB, 하드웨어, 메세지 시스템, 저장소 등 여러 시스템에 대한 익스포터가 존재하며, CollectD 등 기존의 서버 모니터링에 사용되는 에이전트들과 통합할 수 있는 익스포터도 존재합니다.

## node-exporter 설치 (로컬)

로컬 환경에서는 `Docker`로 설치를 진행한다. 터미널에 다음을 입력하여 설치할 수 있다.

```bash
$ docker run --rm -p 9100:9100 prom/node-exporter 
```

그 후 새로운 터미널을 열어 다음 명령어를 입력해서 수집되는 데이터가 있는지 확인한다.

```bash
$ curl localhost:9100/metrics
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
...
```

## node-exporter 설치 (서버)

## Prometheus 연동

## Grafana 대시보드 구축 