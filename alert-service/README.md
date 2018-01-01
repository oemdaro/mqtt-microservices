# alert-service

The Alert Service written in Java. It will send SMS alert to user using Twilio SMS service as long as Temperature or Humidity go high.

## Run Application

Replace parameters in `resources/app.properties` to match your Twilio account.

```properties
#
app.worker.pool                 =   100
app.threshold.temperature       =   30.0
app.threshold.humidity          =   60.0
#
twilio.account.sid              =   YOUR_TWILIO_ACCOUNT_SID
twilio.auth.token               =   YOUR_TWILIO_TOKEN
twilio.from                     =   YOUR_TWILIO_PHONE_NUMBER
twilio.to                       =   YOUR_PHONE_NUMBER
#
kafka.peers                     =   localhost:9092
kafka.topics                    =   mqtt.data
kafka.group.id                  =   alert-consumer
kafka.enable.auto.commit        =   false
kafka.auto.commit.interval.ms   =   1000
kafka.session.timeout.ms        =   30000
kafka.key.deserializer          =   org.apache.kafka.common.serialization.StringDeserializer
kafka.value.deserializer        =   org.apache.kafka.common.serialization.StringDeserializer
kafka.max.pool.records          =   80
#
```

And then run the following command

```bash
$ ./gradlew clean installDist
$ ./build/install/alert-service/bin/alert-service
```
