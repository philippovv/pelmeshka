module "s3bucket" {
  source    = "./modules/s3bucket"
}

module "k8s" {
  source    = "./modules/k8s"
}