# Dolorosa - сервис онлайн контроля исходящих операций

- проводит проверки на операциях и принимает решение об отклонении
- собирает и отправляет лог о принятом решении в аналитический кликхаус через кафку
- взаимодействует с сервисом Nirvana (мастер сервис исключений)


## Test request 
```shell
  grpcurl -plaintext \                                                                                        
  -d '{
    "client_id": "test_client_id",
    "amount": 500000,
    "sender_phone": "+79271020000",
    "receiver_phone": "+79271020001",
    "receiver_bic": "044525225",
    "receiver_name": "Иванов Иван Иванович"
  }' \
  localhost:50051 \
  control.OnlineControl/CheckSBP
```

