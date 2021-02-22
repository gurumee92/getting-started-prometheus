FROM adoptopenjdk/openjdk11

COPY ./build/libs/*.jar /application.jar

RUN mkdir /logs

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "/application.jar"]