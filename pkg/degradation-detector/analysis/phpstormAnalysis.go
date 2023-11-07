package analysis

import detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"

func GeneratePhpStormSettings() []detector.Settings {
  tests := []string{"drupal8-master-with-plugin/inspection",
    "shopware/inspection",
    "b2c-demo-shop/inspection",
    "magento/inspection",
    "wordpress/inspection",
    "laravel-io/inspection",
    "mediawiki/inspection", "php-cs-fixer/inspection", "proxyManager/inspection",
    "akaunting/inspection", "aggregateStitcher/inspection", "prestaShop/inspection", "kunstmaanBundlesCMS/inspection",
    "mpdf/localInspection", "WI_65655/localInspection", "WI_59961/localInspection", "bitrix/localInspection", "WI_65893/localInspection",
    "b2c-demo-shop/indexing",
    "bitrix/indexing",
    "oro/indexing",
    "ilias/indexing",
    "magento2/indexing",
    "drupal8-master-with-plugin/indexing",
    "laravel-io/indexing",
    "wordpress/indexing",
    "mediawiki/indexing",
    "WI_66681/indexing",
    "akaunting/indexing", "aggregateStitcher/indexing", "prestaShop/indexing", "kunstmaanBundlesCMS/indexing", "shopware/indexing",
    "WI_39333-5x/indexing",
    "php-cs-fixer/indexing",
    "many_classes/indexing",
    "magento/indexing",
    "proxyManager/indexing",
    "dql/indexing",
    "tcpdf/indexing",
    "WI_51645/indexing",
    "empty_project/indexing", "complex_meta/indexing", "WI_53502-10x/indexing", "many_array_access/indexing-10x", "WI_66279-10x/indexing",
    "many_classes/completion/classes",
    "magento2/completion/function_var",
    "magento2/completion/function_stlr",
    "magento2/completion/classes",
    "dql/completion",
    "WI_64694/completion",
    "WI_58919/completion",
    "WI_58807/completion",
    "WI_58306/completion",
    "mpdf/inlineRename",
  }
  settings := make([]detector.Settings, 0, 100)
  for _, test := range tests {
    metrics := getMetricFromTestName(test)
    for _, metric := range metrics {
      settings = append(settings, detector.Settings{
        Db:          "perfint",
        Table:       "phpstorm",
        Channel:     "phpstorm-performance-degradations",
        Branch:      "master",
        Machine:     "intellij-linux-hw-hetzner%",
        Test:        test,
        Metric:      metric,
        ProductLink: "phpstorm",
      })
    }

  }
  return settings
}
