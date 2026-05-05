variable "GO_VERSION" {
  default = null
}

variable "DESTDIR" {
  default = "./bin"
}

target "_common" {
  args = {
    GO_VERSION = GO_VERSION
  }
}

group "default" {
  targets = ["test"]
}

target "test" {
  inherits = ["_common"]
  target = "test-coverage"
  output = ["${DESTDIR}/coverage"]
}

target "binary" {
  inherits = ["_common"]
  target = "binary"
  output = ["${DESTDIR}"]
}

target "vendor" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "update"
  output = ["."]
}

group "validate" {
  targets = ["lint", "vendor-validate"]
}

target "lint" {
  target = "lint"
  output = ["type=cacheonly"]
}

target "vendor-validate" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "validate"
  output = ["type=cacheonly"]
}
