# Kafka System 

---

## Requisitos

* Docker
* Docker Compose
* go1.23.2 linux/amd64

---

## Passo a passo para rodar

### 1. Coloque a credencial da Cloud na raiz do projeto

Pegue o arquivo:

```
v1sensor-b8cd9b8b9de3.json
```

E coloque **na raiz do projeto** (mesmo nível do `docker-compose.yml`).

---

### 2. Exporte a credencial

No terminal, dentro da raiz do projeto:

```bash
export GOOGLE_APPLICATION_CREDENTIALS=$(pwd)/v1sensor-b8cd9b8b9de3.json
```

> Esse export é necessário para o Go conseguir acessar a Cloud.

---

### 3. Suba o Docker Compose

```bash
docker-compose up -d
```

Durante a subida, **verifique se aparece** algo como:

```
[+] Running 4/4
 ✔ Network kafka_system_default  Created                                                                                                                                                                              0.1s 
 ✔ Container zookeeper           Started                                                                                                                                                                              0.1s 
 ✔ Container kafka               Started                                                                                                                                                                              0.1s 
 ✔ Container kafka-init          Started  
```

#### ⚠️ Se NÃO aparecer essa linha

Execute:

```bash
docker-compose down
docker-compose up -d
```

---

### 4. Rode a aplicação Go

Com os containers já rodando:

```bash
go run main.go
```

---

## Observações

* O Kafka e o tópico são criados automaticamente pelo Docker
* Não é necessário rodar comandos `docker exec`

---

Pronto. O sistema já deve estar rodando.
