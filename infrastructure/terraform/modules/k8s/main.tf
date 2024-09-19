locals {
  folder_id   = "b1guh01d66fll6dkc716"
}

resource "yandex_kubernetes_cluster" "momo-store-k8s-cluster" {
 name = "momo-store-k8s-cluster"
 network_id = yandex_vpc_network.momo.id
 master {
    public_ip = true
   master_location {
     zone      = yandex_vpc_subnet.subnet_momo_1.zone
     subnet_id = yandex_vpc_subnet.subnet_momo_1.id
   }
 }
 service_account_id      = yandex_iam_service_account.k8s-account-claster.id
 node_service_account_id = yandex_iam_service_account.k8s-account-claster.id
   depends_on = [
     yandex_resourcemanager_folder_iam_member.k8s-clusters-agent,
     yandex_resourcemanager_folder_iam_member.vpc-public-admin,
     yandex_resourcemanager_folder_iam_member.images-puller,
     yandex_resourcemanager_folder_iam_member.encrypterDecrypter
   ]
}

resource "yandex_vpc_network" "momo" { 
  name = "momo" 
}

resource "yandex_vpc_subnet" "subnet_momo_1" {
 name           = "subnet_momo_1"
 v4_cidr_blocks = ["10.1.0.0/16"]
 zone           = "ru-central1-a"
 network_id     = yandex_vpc_network.momo.id
}

resource "yandex_iam_service_account" "k8s-account-claster" {
 name        = "k8s-account-claster"
 description = "K8S zonal service account"
}

resource "yandex_resourcemanager_folder_iam_member" "k8s-clusters-agent" {
  # Сервисному аккаунту назначается роль "k8s.clusters.agent".
  folder_id = local.folder_id
  role      = "k8s.clusters.agent"
  member    = "serviceAccount:${yandex_iam_service_account.k8s-account-claster.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "vpc-public-admin" {
  # Сервисному аккаунту назначается роль "vpc.publicAdmin".
  folder_id = local.folder_id
  role      = "vpc.publicAdmin"
  member    = "serviceAccount:${yandex_iam_service_account.k8s-account-claster.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "images-puller" {
  # Сервисному аккаунту назначается роль "container-registry.images.puller".
  folder_id = local.folder_id
  role      = "container-registry.images.puller"
  member    = "serviceAccount:${yandex_iam_service_account.k8s-account-claster.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "encrypterDecrypter" {
  # Сервисному аккаунту назначается роль "kms.keys.encrypterDecrypter".
  folder_id = local.folder_id
  role      = "kms.keys.encrypterDecrypter"
  member    = "serviceAccount:${yandex_iam_service_account.k8s-account-claster.id}"
}

resource "yandex_vpc_security_group" "k8s-public-services" {
  name        = "k8s-public-services"
  description = "Правила группы разрешают подключение к сервисам из интернета. Примените правила только для групп узлов."
  network_id  = yandex_vpc_network.momo.id
  ingress {
    protocol          = "TCP"
    description       = "Правило разрешает проверки доступности с диапазона адресов балансировщика нагрузки. Нужно для работы отказоустойчивого кластера Managed Service for Kubernetes и сервисов балансировщика."
    predefined_target = "loadbalancer_healthchecks"
    from_port         = 0
    to_port           = 65535
  }
  ingress {
    protocol          = "ANY"
    description       = "Правило разрешает взаимодействие мастер-узел и узел-узел внутри группы безопасности."
    predefined_target = "self_security_group"
    from_port         = 0
    to_port           = 65535
  }
  ingress {
    protocol          = "ANY"
    description       = "Правило разрешает взаимодействие под-под и сервис-сервис. Укажите подсети вашего кластера Managed Service for Kubernetes и сервисов."
    v4_cidr_blocks    = concat(yandex_vpc_subnet.subnet_momo_1.v4_cidr_blocks)
    from_port         = 0
    to_port           = 65535
  }
  ingress {
    protocol          = "ICMP"
    description       = "Правило разрешает отладочные ICMP-пакеты из внутренних подсетей."
    v4_cidr_blocks    = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
  }
  ingress {
    protocol          = "TCP"
    description       = "Правило разрешает входящий трафик из интернета на диапазон портов NodePort. Добавьте или измените порты на нужные вам."
    v4_cidr_blocks    = ["0.0.0.0/0"]
    from_port         = 30000
    to_port           = 32767
  }
  egress {
    protocol          = "ANY"
    description       = "Правило разрешает весь исходящий трафик. Узлы могут связаться с Yandex Container Registry, Yandex Object Storage, Docker Hub и т. д."
    v4_cidr_blocks    = ["0.0.0.0/0"]
    from_port         = 0
    to_port           = 65535
  }
}

# Создание группы узлов
resource "yandex_kubernetes_node_group" "k8s-momo-ng" {
  name        = "k8s-momo-ng"
  description = "Группа узлов магазина Момо"
  cluster_id  = yandex_kubernetes_cluster.momo-store-k8s-cluster.id
  version     = "1.27"
  instance_template {
    name = "momo-{instance.short_id}-{instance_group.id}"
    platform_id = "standard-v3"
    resources {
      cores         = 2
      core_fraction = 50
      memory        = 4
    }
    boot_disk {
      size = 64
      type = "network-ssd"
    }
    network_acceleration_type = "standard"
    network_interface {
      security_group_ids = [yandex_vpc_security_group.k8s-public-services.id]
      subnet_ids         = [yandex_vpc_subnet.subnet_momo_1.id]
      nat                = true
    }
    scheduling_policy {
      preemptible = true
    }
  }
  scale_policy {
    fixed_scale {
      size = 2
    }
  }
  deploy_policy {
    max_expansion   = 3
    max_unavailable = 1
  }
  maintenance_policy {
    auto_upgrade = true
    auto_repair  = true
    maintenance_window {
      start_time = "22:00"
      duration   = "10h"
    }
  }
  node_labels = {
    node-label1 = "node-value1"
  }
  labels = {
    "template-label1" = "template-value1"
  }
  allowed_unsafe_sysctls = ["kernel.msg*", "net.core.somaxconn"]
}

