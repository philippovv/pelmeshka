# Momo Store aka Пельменная №2

<img width="900" alt="image" src="https://storage.yandexcloud.net/picture-for-site/front1.jpg">


### https://best-store-dumpling.ru/

## Frontend

```bash
npm install
NODE_ENV=production VUE_APP_API_URL=http://localhost:8081 npm run serve
```

## Backend

```bash
go run ./cmd/api
go test -v ./... 
```

## [Momo-store]
https://best-store-dumpling.ru -- по этой ссылке можно перейти на сайт

## CI/CD

- ссылка на репозиторий https://gitlab.praktikum-services.ru/std-028-22/dumplings.git
- Развертывание приложения осуществляется через [Downstream pipeline](https://docs.gitlab.com/ee/ci/pipelines/downstream_pipelines.html#parent-child-pipelines), что повышает эффективность и упрощает процесс CI/CD.
- При внесении изменений в соответствующие директории автоматически инициируются пайплайны для компонентов `backend`, `frontend` и `инфраструктурных компонентов`, что обеспечивает быструю реакцию на изменения кода.
- В текущем релизе процесс развертывания `backend` и `frontend` проходят через последовательные этапы `сборки`, `тестирования`, `релиза` и `деплоя` в ходе которых на выходе сохраняется для каждого образ контейнера.
- Модуль momo-helm-chart для k8s проходит стадии релиза (в хранилище nexus) и деплоя в production environment (Kubernetes).
- Вся интеграция осуществляется в соответствии с принципами trunk-based development, что способствует повышению качества кода и ускорению процессов.

## Versioning

- Мы следуем стандартам [SemVer 2.0.0](https://semver.org/lang/ru/) для управления версиями нашего программного обеспечения.
- Мажорные и минорные версии приложения обновляются вручную в файлах `.gitlab-ci.yaml...`, в соответствующей переменной `VERSION`.
- Автоматическое изменение патч-версий осуществляется на основе переменной `CI_PIPELINE_ID`.
- Для инфраструктуры версия приложения фиксируется вручную в чарте `infrastructure/helm-chart/momo-store-chart/Chart.yaml`.

## Infrastructure

- Кодовая база хранится и управляется в `GitLab` ---> [Gitlab](https://gitlab.praktikum-services.ru/)
- Все чарты для Kubernetes хранятся в `Nexus` ---> [Nexus](https://nexus.praktikum-services.ru/)
- Для обеспечения качества кода и устранения проблем на ранних этапах производятся UNIT тесты и тесты на`SonarQube` ---> [SonarQube](https://sonarqube.praktikum-services.ru/)
- Docker-образы хранятся в `GitLab Container Registry` ---> [Container Registry](https://gitlab.praktikum-services.ru/std-028-22/dumplings/container_registry)
- Для хранения состояния Terraform и статических файлов используется `Yandex Object Storage` ---> [Yandex Object Storage](https://cloud.yandex.ru/services/storage)
- Развертывание и управление приложениями в productions осуществляется на `Yandex Managed Service for Kubernetes` ---> [Yandex Managed Service for Kubernetes](https://cloud.yandex.ru/services/managed-kubernetes)

## Init kubernetes

### Скрипт terraform выполняет следующие задачи:
- Создает K8S-кластер
- Создает сеть и подсети
- Настраивает сетевые правила
- Cоздаёт сервисных пользователей и настраивает их права
- Cоздаёт бакет для загрузки terrform состояний в S3 Object Storage
- Поднимает группу нод

### Для использования скрипта необходимо:

- Клонировать репозиторий на машину с установленным [`terraform`](https://developer.hashicorp.com/terraform/downloads)
- Установить [`консоль Yandex`](https://cloud.yandex.ru/docs/cli/operations/install-cli), настроить свой [`профиль в Yandex Cloud`](https://yandex.cloud/ru/docs/cli/operations/profile/profile-create)
- Установить [`Kubectl`](https://kubernetes.io/ru/docs/tasks/tools/install-kubectl/) для работы с K8S
- Установить [`Helm CLI`](https://helm.sh/docs/intro/install/)
- Создать сервисный аккаунт с ролью `editor` через консоль Yandex Cloud, создать статический ключ доступа, сохранить секретные ключи в переменные окружения 

```
export AWS_ACCESS_KEY_ID = Идентификатор ключа
export AWS_SECRET_ACCESS_KEY = Секретный ключ
```

- Получить [iam-token](https://cloud.yandex.ru/docs/iam/operations/iam-token/create), сохранить в переменной окружения 
```
export YC_TOKEN = Полученный ключ
```
- Выполнить следующие комманды:

```
cd infrastructure\terraform
terraform init
terraform apply
```

### Подключение к К8S-кластеру:
```
yc managed-kubernetes cluster get-credentials <имя_или_идентификатор_кластера> --external
kubectl cluster-info
```

## Init production. Работа с kubectl и helm charts

- Создать базовый namespace
```
kubectl create namespace momo-store
```
- Установить [`Vertical Pod Autoscaler`] 
    cd /tmp && \
    git clone https://github.com/kubernetes/autoscaler.git && \
    cd autoscaler/vertical-pod-autoscaler/hack && \
    ./vpa-up.sh

- Проверка, что vpa установлен успешно:
    kubectl describe vpa

- В интерфейсе Yandex Cloud ---> Cloud DNS необходимо создать зону по имени хоста.
    Там будут добавляться записи для корректного хостинга трафика.

- Установить `Ingress-контроллер NGINX`. 
    Установка:  helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx && \
                helm repo update && \
                helm install ingress-nginx ingress-nginx/ingress-nginx

- Установить `ExternalDNS c плагином для Yandex Cloud DNS`.
    ExternalDNS позволяет автоматически создавать DNS-записи в Yandex Cloud DNS. 
    ExternalDNS делает ресурсы Kubernetes доступными для обнаружения через общедоступные DNS-серверы.
        Установка:  export HELM_EXPERIMENTAL_OCI=1 && \
                    helm pull oci://cr.yandex/yc-marketplace/yandex-cloud/externaldns/chart/externaldns \
                    --version 0.5.1 \
                    --untar && \
                    helm install \
                    --namespace <пространство_имен> \
                    --create-namespace \
                    --set config.folder_id=<идентификатор_каталога_с_DNS-зоной> \
                    --set-file config.auth.json=<путь_к_файлу_с_авторизованным_ключом_сервисного_аккаунта> \
                    externaldns ./externaldns/
    Более подробно ознакомиться с работой `ExternalDNS c плагином для Yandex Cloud DNS` можно перейджя по ссылке (https://yandex.cloud/ru/docs/managed-kubernetes/operations/applications/externaldns)

- Установить `cert-manager`. На текущий момент выполняется в ручную, следующими командами:
    Установка: kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.1/cert-manager.yaml
    Првоерка, что после установки, запущено 3 пода: kubectl get pods -n cert-manager --watch

- Создать `ClusterIssuer`.
     С помощью ClusterIssuer можно выпускать сертификаты Let's Encrypt®. Сертификаты будут выпускаться после прохождения проверки HTTP-01 с помощью установленного ранее Ingress-контроллера.
        Установка:  cd infrastructure\helm_charts
                    kubectl apply -f http-clusterissuer.yaml

```
## [Monitoring]
- Написаны отдельные helm чарты для `grafana`, `prometheus`.
- Данные для входа по умолчанию: `admin:admin` 

- Grafana Dashboard
    Для получение домена иказанного в Ingress, требуется ввести команду: kubectl get ingress grafana -n <namespace>
    В текущей реализации на дашборде прикреплены две панели графиков, с отслеживаением бизнес метрик.
- Prometheus 
    Для получение домена иказанного в Ingress, требуется ввести команду: kubectl get ingress prometheus -n <namespace>

```
- Для ускореняи работы backend, картинки хранятся в S3 bucket. 

## Backlog
- Поднять и настроить `ArgoCD` для синхронизации состояний приложений
- Настроить сбор метрик кластера k8s , доработать текущий бизнес dashboard
- Настроить `AlertManager` и Prometheus rules для отправки уведомлений
- Поднять `Vault` для хранения секретов
- Параметризировать Helm чарты приложений
- Параметризировать Terraform манифесты
- Написать pipline для автоматизации развертывания инфраструктуры через манифесты Terraform
- Публикация образов frontend и backend в репозиторий Nexus
- Развернуть и настроить сбор логов с помощью Loki.
- Настроить динамическое версионирование хелм чартов, которые отправляются в Nexus
- Сформирвоать helm chart для установки: `Ingress-контроллер NGINX`, `ExternalDNS c плагином для Yandex Cloud DNS`, `cert-manager`, `ClusterIssuer`. Добавить в pipline.