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

## prometheus-nginxlog-exporter와 설치

## 메트릭 수집을 위한 각 컴포넌트 설정

## 서비스 메트릭 모니터링을 위한 Grafana 대시보드 구축