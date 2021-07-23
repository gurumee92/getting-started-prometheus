# 자바 어플리케이션 모니터링하기 (1) spring-boot

## 개요

이 문서에서는 `Prometheus`로 `spring-boot`기반의 자바 애플리케이션의 메트릭을 수집한 후 `Grafana` 대시보드를 구축하는 것에 대하여 다룬다. 자세한 내용은 다음과 같다.

* 자바, 프로젝트 설치
* Spring Boot Application 설정 살펴보기
* Prometheus 설정
* Spring Boot Application 서버 모니터링을 위한 Grafana 대시보드 구축

이 문서에서 진행되는 실습 코드는 편의성을 위해 로컬 `Docker` 환경에서 진행되나, 실세 서버 환경에서도 거의 동일하게 적용할 수 있도록 작성되었다. 이번 장의 코드는 다음 링크에서 확인할 수 있다.

* 이번 장 코드 : [https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch08](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch08)

이 문서에서 구성하는 인프라스트럭처는 다음과 같다.

![01](./01.png)

## 자바, 프로젝트 설치 (서버 환경)

먼저 자바11을 설치한다.

```bash
# 패키지 인스톨
$ sudo yum install java-11-openjdk-devel

# 자바 버전 확인
$ java -version
openjdk version "11.0.11" 2021-04-20 LTS
OpenJDK Runtime Environment 18.9 (build 11.0.11+9-LTS)
OpenJDK 64-Bit Server VM 18.9 (build 11.0.11+9-LTS, mixed mode, sharing)
```

만약 11버전이 아니라 8(1.8)이 나온다면 서버를 재부팅한다. 그럼 11버전이 나올 것이다. 이제 프로젝트를 서버에 설치해보자.

```bash
$ pwd
/home/sidelineowl

# git clone
$ git clone https://github.com/gurumee92/gurumee-book-prometheus.git

# java app 소스 코드 복제
$ cp -R ./gurumee-book-prometheus/src/part2/ch08/app ~/apps/java-app

# java app 디렉토리로 이동
$ cd ~/apps/java-app

# gradlew 실행 권한 주기
$ chmod +x gradlew

# jar 빌드
$ ./gradlew bootJar

# 애플리케이션 실행
$ java -jar build/libs/prometheus-example-0.0.1-SNAPSHOT.jar 
  .   ____          _            __ _ _
 /\\ / ___'_ __ _ _(_)_ __  __ _ \ \ \ \
( ( )\___ | '_ | '_| | '_ \/ _` | \ \ \ \
 \\/  ___)| |_)| | | | | || (_| |  ) ) ) )
  '  |____| .__|_| |_|_| |_\__, | / / / /
 =========|_|==============|___/=/_/_/_/
 :: Spring Boot ::                (v2.4.3)
2021-07-23 08:06:34.501  INFO 3465 --- [           main] c.gurumee.prometheusexample.Application  : Starting Application using Java 11.0.11 on gbp-02 with PID 3465 (
/home/sidelineowl/apps/java-app/build/libs/prometheus-example-0.0.1-SNAPSHOT.jar started by sidelineowl in /home/sidelineowl/apps/java-app)
2021-07-23 08:06:34.507  INFO 3465 --- [           main] c.gurumee.prometheusexample.Application  : No active profile set, falling back to default profiles: default
2021-07-23 08:06:36.607  INFO 3465 --- [           main] .s.d.r.c.RepositoryConfigurationDelegate : Bootstrapping Spring Data JPA repositories in DEFAULT mode.
2021-07-23 08:06:36.665  INFO 3465 --- [           main] .s.d.r.c.RepositoryConfigurationDelegate : Finished Spring Data repository scanning in 10 ms. Found 0 JPA re
pository interfaces.
2021-07-23 08:06:37.659  INFO 3465 --- [           main] o.s.b.w.embedded.tomcat.TomcatWebServer  : Tomcat initialized with port(s): 8080 (http)
2021-07-23 08:06:37.682  INFO 3465 --- [           main] o.apache.catalina.core.StandardService   : Starting service [Tomcat]
2021-07-23 08:06:37.682  INFO 3465 --- [           main] org.apache.catalina.core.StandardEngine  : Starting Servlet engine: [Apache Tomcat/9.0.43]
2021-07-23 08:06:37.798  INFO 3465 --- [           main] o.a.c.c.C.[Tomcat].[localhost].[/]       : Initializing Spring embedded WebApplicationContext
2021-07-23 08:06:37.799  INFO 3465 --- [           main] w.s.c.ServletWebServerApplicationContext : Root WebApplicationContext: initialization completed in 3131 ms
2021-07-23 08:06:38.535  INFO 3465 --- [           main] com.zaxxer.hikari.HikariDataSource       : HikariPool-1 - Starting...
2021-07-23 08:06:38.842  INFO 3465 --- [           main] com.zaxxer.hikari.HikariDataSource       : HikariPool-1 - Start completed.
2021-07-23 08:06:39.019  INFO 3465 --- [           main] o.hibernate.jpa.internal.util.LogHelper  : HHH000204: Processing PersistenceUnitInfo [name: default]
2021-07-23 08:06:39.113  INFO 3465 --- [           main] org.hibernate.Version                    : HHH000412: Hibernate ORM core version 5.4.28.Final
2021-07-23 08:06:39.376  INFO 3465 --- [           main] o.hibernate.annotations.common.Version   : HCANN000001: Hibernate Commons Annotations {5.1.2.Final}
2021-07-23 08:06:39.615  INFO 3465 --- [           main] org.hibernate.dialect.Dialect            : HHH000400: Using dialect: org.hibernate.dialect.H2Dialect
2021-07-23 08:06:40.027  INFO 3465 --- [           main] o.h.e.t.j.p.i.JtaPlatformInitiator       : HHH000490: Using JtaPlatform implementation: [org.hibernate.engin
e.transaction.jta.platform.internal.NoJtaPlatform]
2021-07-23 08:06:40.045  INFO 3465 --- [           main] j.LocalContainerEntityManagerFactoryBean : Initialized JPA EntityManagerFactory for persistence unit 'defaul
t'
2021-07-23 08:06:40.172  WARN 3465 --- [           main] JpaBaseConfiguration$JpaWebConfiguration : spring.jpa.open-in-view is enabled by default. Therefore, databas
e queries may be performed during view rendering. Explicitly configure spring.jpa.open-in-view to disable this warning
2021-07-23 08:06:40.555  INFO 3465 --- [           main] o.s.s.concurrent.ThreadPoolTaskExecutor  : Initializing ExecutorService 'applicationTaskExecutor'
2021-07-23 08:06:41.044  INFO 3465 --- [           main] o.s.b.a.e.web.EndpointLinksResolver      : Exposing 1 endpoint(s) beneath base path '/actuator'
2021-07-23 08:06:41.149  INFO 3465 --- [           main] o.s.b.w.embedded.tomcat.TomcatWebServer  : Tomcat started on port(s): 8080 (http) with context path ''
2021-07-23 08:06:41.183  INFO 3465 --- [           main] c.gurumee.prometheusexample.Application  : Started Application in 7.858 seconds (JVM running for 8.714)
...
```

이번엔 애플리케이션을 리눅스 서비스로 말고 백그라운드 프로세스로 애플리케이션을 실행해보자.

```bash
$ pwd
/home/sidelineowl/apps/java-app

# 백그라운드 실행
$ nohup java -jar build/libs/prometheus-example-0.0.1-SNAPSHOT.jar 1>/dev/null 2>&1 &

# 메트릭 수집 엔드포인트 호출
$ curl localhost:8080/actuator/prometheus
# HELP tomcat_sessions_rejected_sessions_total  
# TYPE tomcat_sessions_rejected_sessions_total counter
tomcat_sessions_rejected_sessions_total{application="example",} 0.0
# HELP jvm_buffer_count_buffers An estimate of the number of buffers in the pool
# TYPE jvm_buffer_count_buffers gauge
jvm_buffer_count_buffers{application="example",id="mapped",} 0.0
jvm_buffer_count_buffers{application="example",id="direct",} 1.0
...
```

이렇게 나오면 성공이다.

## Spring Boot Application 설정 살펴보기

먼저 스프링 부트 기반의 자바 애플리케이션이라면, `spring-boot-starter-actuator`와 `micrometer-registry-prometheus` 의존성이 필요하다. 보통 `gradle` 혹은 `maven`이라는 빌드 툴로 관리하는데, 각각의 도구에서 다음 코드처럼 의존성을 명시하면 된다.

> 참고! 빌드 도구에 따른 의존성 관리 파일 경로
> "Gradle"의 경우에는 build.gradle, "Maven"의 경우에는 pom.xml이 프로젝트 루트 디렉토리 최상단에 존재합니다. 이들을 수정하면 됩니다. 이 문서에서는 "Gradle"만 다룹니다.

[src/part2/ch08/app/src/build.gradle](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch08/app/src/build.gradle)
```gradle
// ...

dependencies {
    // ...
	implementation 'org.springframework.boot:spring-boot-starter-actuator'
	runtimeOnly 'io.micrometer:micrometer-registry-prometheus'
    // ...
}

// ...
```

그리고 `application.yml`에 다음을 적어주면 된다.

[src/part2/ch08/app/src/main/resources/application.yml](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch08/app/src/main/resources/application.yml)
```yml
spring:
  application:
    name: example
management:
  endpoints:
    web:
      exposure:
        include: "prometheus"
  metrics:
    tags:
      application: ${spring.application.name}
```

위 설정을 가지고 있을 때 스프링 부트 기반 WAS는 다음 엔드포인트를 제공한다. 

* http://localhost:8080/actuator/prometheus

위 엔드포인트를 들어가게 되면 다양한 애플리케이션에서 수집되는 메트릭들을 확인할 수 있다.

![02](./02.png)

수집되는 정보는 크게 다음과 같다.

* 로그 이벤트
* process cpu 정보
* hikaricp 풀 정보 (DB 연결 정보)
* http 서버 요청/응답 정보
* jvm 메모리 정보

아마 자바/스프링 부트 기반의 애플리케이션을 접하지 않은 사람이라면 위 정보가 익숙하지 않을 것이다. 그냥 자바 애플리케이션 모니터링 시 필요한 필수 정보라고 생각하자. 여기까지 하면 일단은 자바/스프링 부트 기반의 WAS에서 할 수 있는 설정은 완료하였다.

## Prometheus 설정

`Prometheus`는 역시 다음과 같이 설정 파일을 수정하면 된다. 서버 환경에서라면, `/etc/prometheus/prometheus.yml`을 수정하자.

[src/part2/ch08/prometheus/prometheus.yml](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch08/prometheus/prometheus.yml)
```yml
# my global config
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # By default, scrape targets every 15 seconds.

  external_labels:
    monitor: 'my-project'

rule_files:

scrape_configs:
  # ...

  - job_name: 'spring-boot-application'
    scrape_interval: 5s
    # 여기에서는 spring actuator가 활성화시키는 엔드포인트에서 데이터를 스크래핑한다.
    metrics_path: "/actuator/prometheus"

    static_configs:
      # 자바 애플리케이션이 실행되는 IP:PORT
      - targets: ["app:8080"]
```

잘 수집되는지 확인하려면 `Prometheus UI`에 다음을 쿼리해보자.

```
jvm_gc_max_data_size_bytes
```

그럼 다음과 같은 결과를 얻을 수 있다.

![03](./03.png)

## Spring Boot Application 서버 모니터링을 위한 Grafana 대시보드 구축

이제 대시보드를 구축한다. 다음 JSON 파일을 복사해서 대시보드를 임포트한다. (로컬 환경에는 이미 대시보드가 로드되어 있다.) 다음 링크로 가서 JSON 파일을 복사한다.

* [src/part2/ch08/grafana/dashboard.json](https://github.com/gurumee92/gurumee-book-prometheus/tree/master/src/part2/ch08/grafana/dashboard.json)

먼저 그라파나에 접속한다. 그 후 "+" 버튼을 누른다.

![04](./04.png)

그럼 아래와 같이 메뉴가 보이는데 "Import"를 누른다.

![05](./05.png)

그 후 위 링크에서 제공하고 있는 json 파일을 복사하여 붙여넣고 "Load"를 누른다.

![06](./06.png)

그럼 위와 같이 입력값들이 자동적으로 채워진다. "Import"를 누른다.

![07](./07.png)

그럼 다음 대시보드가 구축된다.

![08](./08.png)

### Basic Statistics

![09](./09.png)

다음 대시보드에서 확인할 수 있는 지표는 다음과 같다.

* uptime 
* start time
* heap 메모리 사용량
* non-heap 메모리 사용량
* process open file 지표
* process cpu 사용량
* system load 관련 지표

### JVM Statistics - Memory

![10](./10.png)

* JVM 로드된 클래스 개수, 로드되지 않은 클래스 개수
* JVM 버퍼 사용량 (direct, map)
* 쓰레드 관련 지표
* GC Memory 할당/Promote(GC 지역 사이를 돌아다니는)량

![11](./11.png)

* GC 힙 영역 공간 지표
* GC 논-힙 영역 공간 지표

### JVM Statistics - GC

![12](./12.png)

* GC 락 개수
* GC 락 시간 합계

### HikariCP Statistics

![13](./13.png)

히카리 풀은 자바의 DB 커넥션 풀 구현체 중 하나이다. 

* 히카리 풀 커넥션 개수
* 히카리 풀 커넥션 타임아웃 개수
* 평균 커넥션 생성 시간
* 평균 커넥션 사용 시간
* 평균 커넥션 획득 시간

### HTTP Statistics

![14](./14.png)

* HTTP 요청 개수
* HTTP 응답 시간

### Log Statistics

![15](./15.png)

* 로그 레벨 별 이벤트 발생 평균 개수
