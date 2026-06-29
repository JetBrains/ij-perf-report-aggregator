export const phpStormUnitTestGroups: { label: string; projects: string[] }[] = [
  {
    label: "Control Flow",
    projects: [
      "com.jetbrains.php.PhpControlFlowBuilderTest.testLoadPerformance$1 - LoadPerformance$1",
      "com.jetbrains.php.PhpControlFlowBuilderTest.testLoadPerformance$2 - LoadPerformance$2",
    ],
  },
  {
    label: "Formatter",
    projects: [
      "com.jetbrains.php.PhpFormatterTest.testPerformance - testPerformance",
      "com.jetbrains.php.PhpFormatterTest.testPerformance1 - testPerformance1",
      "com.jetbrains.php.PhpFormatterTest.testPerformance2 - testPerformance2",
      "com.jetbrains.php.PhpFormatterTest.testPerformance3 - testPerformance3",
      "com.jetbrains.php.PhpFormatterTest.testPerformance4 - testPerformance4",
      "com.jetbrains.php.PhpFormatterTest.testPerformanceWi30101 - testPerformanceWi30101",
    ],
  },
  {
    label: "Typing",
    projects: [
      "com.jetbrains.php.PhpTypingPerformanceTest.testTypingInHeredocLiteralWithVariables - PHP typing in heredoc literal",
      "com.jetbrains.php.PhpTypingPerformanceTest.testTypingInHeredocLiteral - PHP typing in heredoc literal",
      "com.jetbrains.php.PhpTypingPerformanceTest.testTypingInStringLiteral - PHP typing in heredoc literal",
      "com.jetbrains.php.PhpTypingPerformanceTest.testTypingInLargeMixedFile - PHP typing in large mixed file",
    ],
  },
  {
    label: "Completion",
    projects: ["com.jetbrains.php.completion.PhpAccessorsCompletionTest.testPerformanceOnLotsOfFields - Php getters/setters completion"],
  },
  {
    label: "Lang Util",
    projects: ["com.jetbrains.php.lang.PhpLangUtilTest.testIsPhpIdentifierPerformance - IsPhpIdentifierPerformance"],
  },
  {
    label: "Phar File System",
    projects: [
      "com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePhar - indexing _phar file",
      "com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharGz - indexing _phar_gz file",
      "com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharBz2 - indexing _phar_bz2 file",
      "com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharTar - indexing _phar_tar file",
      "com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharTarGz - indexing _phar_tar_gz file",
      "com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharTarBz2 - indexing _phar_tar_bz2 file",
    ],
  },
]
