# 웹 서버 모니링하기 (2) apache-exporter

![logo](../../logo.png)

## 개요

이 문서에서는 `apache-exporter`를 이용해서 `Apache`의 커넥션 정보에 대한 메트릭을 수집한다. 그 후 `Grafana`, `Prometheus`를 이용해서 `Apache` 웹 서버를 모니터링할 수 있는 대시보드를 구축하는 것에 대하여 다룬다. 자세한 내용은 다음과 같다.

* Apache 서버와 설치
* apache-exporter와 설치
* 메트릭 수집을 위한 각 컴포넌트 설정
* Apache 서버 모니터링을 위한 Grafana 대시보드 구축

이 문서에서 진행되는 실습 코드는 편의성을 위해 로컬 `Docker` 환경에서 진행되나, 실세 서버 환경에서도 거의 동일하게 적용할 수 있도록 작성되었다. 이번 장의 코드는 다음 링크에서 확인할 수 있다.

* 이번 장 코드 : [https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch07](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch07)

이 문서에서 구성하는 인프라스트럭처는 다음과 같다.

![01](./01.png)

## Apache 서버와 설치
## apache-exporter와 설치
## 메트릭 수집을 위한 각 컴포넌트 설정
## Apache 서버 모니터링을 위한 Grafana 대시보드 구축

