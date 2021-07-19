# Cortex란 무엇인가

![logo](../../logo.png)

## 개요

이번 장에서는 `Cortex`가 무엇인지 공부한다. "갑자기 생뚱맞게 `Prometheus` 공부하다가 무슨 `Cortex`야?" 라는 의문이 생길 수도 있다. `Prometheus`는 뛰어난 성능과, 쉬운 접근성을 바탕으로 최근 나온 모니터링 기술 중 `InfluxDB`와 함께 업계 표준으로 자리 잡았다. 하지만 다음과 같은 문제점이 존재한다.

1. Prometheus는 scale-out 구조를 고려하지 않고 설계되었다.
2. Prometheus에서 긴 시간 범위를 지닌 데이터를 쿼리할 때 성능이 매우 저하되었다.
3. Prometheus가 저장하는 데이터 특성 상 오랜 시간 저장하는 것이 쉽지 않다.

위 문제점들을 해결하기 위한 대표적인 솔루션이 바로 `Cortex`이다. 단점이라면 백엔드 구성을 S3, GCS 등으로 구성해야 하기 때문에 오픈 소스 솔루션 치고 비싸다라는 점이 있긴 하지만 대규모 시스템을 모니터링 하기 위해서 `Prometheus`를 운영하고 있는 곳이라면, 한 번쯤 고려해볼만 하다.   

> 참고! Loki에 대해서
>
> 현재 업계에서 뜨고 있는 기술 중 하나인 "Loki"는 "Grafana Labs"가 만든 오픈 소스 로그 모니터링 시스템입니다. 이는 "Cortex"를 래핑하여 로그 메세지를 수집 및 쿼리할 수 있게 만들었기 떄문에, 서로의 구조가 매우 유사합니다. 그래서 둘 중 하나를 배워두면 다른 하나를 배우는데 큰 도움이 될 것입니다.

> 참고! Thanos에 대해서
>
> "Cortex"와 더불어서 "Prometheus" HA 솔루션으로 거론되는 것이 바로 "Thanos"입니다. "Thanos"는 "Cortex"에 비해 저렴하게 운영할 수 있으나, 쿼리 성능은 다소 낮습니다. (작게는 2배 크게는 10배 정도) 둘 다 훌륭한 솔루션이며 서로 보완적인 관계에 있으므로 같이 알아두면 큰 도움이 될 것입니다. 

## Cortex란 무엇인가?

`Cortex`는 `CNCF`의 인큐베이션 프로젝트 중 하나로, `Weave Cloud`와 `Grafana Labs`가 관리하는 `Prometheus` 기반 오픈 소스 솔루션이다. `Cortex` 공식 문서에 따르면, 자신을 이렇게 소개하고 있다.

>  "Cortex provides horizontally scalable, highly available, multi-tenant, long term storage for Prometheus."

번역하면 `Cortex`는 `Prometheus`를 위한, 수평적으로 확장 가능한, 고 가용성의, 멀티 테넌트, 롱 텀 스토리지이다. 정의만 보더라도 위에서 지적한 `Prometheus`의 문제점과 매우 밀접한 관계를 맺고 있다는 것을 볼 수 있다. 

문서에 따르면, `Cortex`는 다음의 특징을 갖는다.

* Horizontally Scalable
* Highly Available
* Multi-Tenant
* Long Term Storage

`Cortex`는 "Horizontally scalable" 즉 수평적으로 확장이 가능하다. 이 말은 여러 머신을 클러스터로 구성할 수 있다는 뜻이다. 여기서 중요한 점이 바로 "Globally Aggregated"를 지원한다는 점이다. 가령, 시스템 메트릭 데이터만 수집하는 `Prometheus A`와 액세스 로그 데이터만 수집하는 `Prometheus B`가 있다고 해보자. 하지만 웹서버가 이상 동작을 일으킬 때, 우리는 `Prometheus A`와 `Prometheus B`를 모두 살펴보아야 한다. 하지만 각각의 지표를 봤을 때 해결점을 찾기란 쉽지 않다.

여기서는 메트릭을 기준으로 잡았지만 지역별로 쪼갠다든가 다른 방식으로 `Prometheus`를 나눈다 하더라도 결국 전역적으로 집계할 필요가 생긴다. 이 때 `Cortex`를 이용하면 쉽게 해결할 수 있다. `Prometheus A`와 `Prometheus B` 모두 수집한 데이터를 `Cortex`가 구성된 클러스터에게 넘겨주기만 하면 모두 쿼리 및 집계가 가능하다.

`Cortex`는 "High Available"을 지원한다. 이 말은 클러스터 내 `Cortex`끼리는 데이터를 복제가 가능하다. 즉 클러스터 내에 특정 한 개의 `Cortex`가 망가지더라도 구성된 다른 `Cortex`들이 고장난 녀석의 몫까지 커버할 수 있다는 뜻이다.

또한, "Long Term Storage"를 지원한다. `Cortex`는 AWS의 DynamoDB, S3 혹은 GCP의 BigTable, GCS 등을 데이터 저장소로써 쓸 수 있다. 이 때 Index/Chunk 단위로 데이터를 나누고 Chunk는 압축해서, Index는 DynamoDB/BigTable에 저장되고 Chunk는 S3/GCS 등에 저장된다. 이렇게 저장된 데이터는 일반 `Prometheus`에 저장된 데이터보다 훨씬 긴 라이프사이클을 가지게 된다. 중요한 점은 긴 시간 범위에 해당하는 데이터를 조회하더라도 매우 빠르게 조회가 가능하도록 설계가 되었다는 점이다.

마지막으로 "Multi-Tenant"를 지원한다. 이 부분은 정확하게 뜻이 이해되지 않는데 단일 클러스터에 있는 `Cortex`는 여러 독립적인 `Prometheus` 소스의 데이터와 쿼리를 격리하여 신뢰할 수없는 당사자가 동일한 클러스터를 공유 할 수 있도록 만들어준다고 한다. 이 부분은 공부하면서 나중에 이해 되면 그 때 채워넣기로 하자.

## Cortex 아키텍처

다음은 공식 문서를 기반으로 작성한 `Cortex` 아키텍처의 구성도이다.

![01](./01.png)

위의 관계도를 보면 `Cortex`에서 `Prometheus`는 일종의 중개자 역할이 된다. 각 `Exporter` 및 `Push Gateway`에서 수집된 메트릭들을 `Cortex`로 넘기는 역할을 수행한다. 이 때 `Cortex`는 클러스터 내에서, 역할에 따라 다음과 같이 분류될 수 있다.

* Distributor
* Ingester
* Store Gateway
* Compactor
* Querier
* Query Frontend
* Ruler
* Alertmanager

### Distributor

`Prometheus`가 전송한 데이터를 처리하는 역할을 한다. 전송된 데이터의 유효성을 체크한 후, 문제가 없으면 `Ingester`로 병렬적으로 데이터를 전송한다. 이 때 데이터를 전송할 `Ingester`를 선택할 때 구성된 "Hash Ring"을 통해서 해싱하여 선택한다. 이 때 해싱할 때 데이터의 내부 속성을 이용해서 해싱하게 된다. 이용할 수 있는 내부 속성의 조합은 다음과 같다.

* 메트릭 이름 + 테넌트 ID (기본)
* 메트릭 이름 + 라벨 + 테넌트 ID

또한, 데이터 유효성 체크는 다음을 포함한다.

* 메트릭의 라벨(=태그)들의 이름이 형식적으로 정확한가
* 각 메트릭마다 최대/최소 라벨 개수를 준수하는가 (=최소 ~ 최대 개수 사이인가)
* 수집된 데이터의 timestamp가 최소/최대 시간 범위 안에 있는가

`Distributor`는 "stateless"하며, 필요에 따라 스케일 업/다운이 가능하다. 

### Ingester

`Ingester`는 `Distributor`로 전달 받은 데이터들을 `Storage`에 넘겨주는 역할을 한다. `Cortex`가 지원하고 있는 데이터 저장 방법은 크게 2가지이다.

* Chunks Storage (Deprecated)
* Blocks Storage

`Chunks Storage` 방식이 기본이며, 수신된 데이터를 Index/Chunk를 나누어서 저장한다. 가능한 데이터 저장소로는 `AWS DynamoDB/S3`, `GCP BigTable/GCS`, `Cassandra/Cassandra` 가 있다. 하지만 v1.9 이후부터는 deprecated 되었다. 

`Blocks Storage` 방식은 `Prometheus TSDB`와 유사한 방식으로 데이터를 저장하는 방법이며, `AWS S3`, `GCP GCS` 등을 저장소로 사용할 수 있다.

또한, 앞서 언급했던 것처럼 `Ingester` 들은 각각 토큰과 함께 "Hash Ring"에 구성되는데 "KV Store"가 필요하다. 이를 지원하는 KV Store는 다음과 같다.

* Consul
* Etcd
* Gossip Memberlist

### Store Gateway

`Store Gateway`는 "Blocks"으로부터 시계열을 쿼리하는 서비스이다. `Store Gateway`는 "Blocks Storage" 저장 방식으로 `Cortex` 클러스터를 구성할 때 필수적인 컴포넌트이다. `Store Gateway`는 "semi-stateful"하다. 

### Compactor

`Compactor`는 다음 작업을 수행한다.

* 지정된 테넌트의 여러 "Block"을 최적화된 하나의 큰 "Block"으로 압축하는 역할을 한다. 이를 통해서 스토리지 비용을 절감하고 쿼리 속도를 높인다.
* 테넌트 별 버킷 인덱스를 업데이트 상태로 유지한다. 버킷 인덱스는 `Querier`, `Store Gateway`, `Ruler`가 새 "Block"을 검색하는데 사용된다.

`Compactor`는 "stateless"하다.

### Querier

`Querier`는 이름 그대로 `Cortex` 클러스터의 `Ingester`, `Store Gateway` 그리고 캐시 저장소에 저장된 데이터를 `PromQL`로 쿼리하는 역할을 한다. 복제본으로 인해 중복되는 데이터가 있다면, 내부적으로 이를 제거해서 쿼리 결과로 보여준다.

### Query Frontend

`Query Frontend`는 선택적인 역할로, `Querier`의 쿼리 속도를 높이는데 사용된다. 내부적으로 쿼리 조정을 수행하고 내부 큐의 쿼리들을 보관한다. 만약 `Query Frontend` 역할을 하는 `Cortex`가 있다면, `Querier`들은 이 내부 큐의 작업들을 가져와서 실행한 후 반환하는 작업자들이 된다. 또한 `Query Scheduler`를 따로 구성하면 내부 큐를 외부로 뺄 수 있다.

### Ruler

`Ruler`는 선택적인 역할로써, 규칙과 알림을 기록하기 위해서 `PromQL` 쿼리들을 실행한다. `Ingester`와 유사하게 "Hash Ring"으로 클러스터를 구성하며, `S3`, `GCS` 등의 장기 저장소에 `rule`들을 저장한다. `Ruler`는 "semi-stateful"하며, 병렬적으로 스케일할 수 있다. 

### Alertmanager

`Alertmanager`는 선택적인 역할로,  `Prometheus AlertManager` 기반으로 작성되었다. `Ingester`와 유사하게 "Hash Ring"으로 클러스터를 구성하며 `Cortex Ruler`를 통해서 생성된 알람들을 `Slack` 등의 외부 엔드포인트로 전송하는 역할을 한다. `Alertmanager`는 "semi-stateful"하며, 병렬적으로 스케일할 수 있다.