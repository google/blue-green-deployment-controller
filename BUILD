load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "helpers.go",
        "main.go",
    ],
    importpath = "k8s.io/bgd-controller",
    visibility = ["//visibility:private"],
    deps = [
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/api/extensions/v1beta1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/bgd-controller/pkg/apis/demo/v1:go_default_library",
        "//vendor/k8s.io/bgd-controller/pkg/client/clientset/versioned:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/typed/extensions/v1beta1:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/k8s.io/client-go/util/retry:go_default_library",
    ],
)

go_binary(
    name = "sample-controller",
    embed = [":go_default_library"],
    importpath = "k8s.io/sample-controller",
    visibility = ["//visibility:public"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//staging/src/k8s.io/bgd-controller/pkg/apis/demo:all-srcs",
        "//staging/src/k8s.io/bgd-controller/pkg/client/clientset/versioned:all-srcs",
        "//staging/src/k8s.io/bgd-controller/pkg/client/informers/externalversions:all-srcs",
        "//staging/src/k8s.io/bgd-controller/pkg/client/listers/demo/v1:all-srcs",
        "//staging/src/k8s.io/bgd-controller/pkg/signals:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "bgd-controller",
    embed = [":go_default_library"],
    importpath = "k8s.io/bgd-controller",
    visibility = ["//visibility:public"],
)
