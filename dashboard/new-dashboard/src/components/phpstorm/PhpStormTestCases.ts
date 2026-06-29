export const phpStormUnitTestGroups: { label: string; projects: string[] }[] = [
  {
    label: "Injection",
    projects: ["com.jetbrains.php.slowTests.PhpInjectionSlowTest.testDeepConcatenation"],
  },
  {
    label: "Resolve",
    projects: [
      "com.jetbrains.php.slowTests.PhpResolveSlowTest.testGotoClassCanWorkInDumbMode",
      "com.jetbrains.php.slowTests.PhpResolveSlowTest.testNoSOEOnIndexingDeepMemberAccess",
      "com.jetbrains.php.slowTests.PhpResolveSlowTest.testTooBigControlFlow",
      "com.jetbrains.php.slowTests.PhpResolveSlowTest.testBigControlFlow",
    ],
  },
  {
    label: "Control Flow",
    projects: [
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpControlFlowBuilderTest.testLoadPerformance$1",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpControlFlowBuilderTest.testLoadPerformance$2",
    ],
  },
  {
    label: "Editor",
    projects: [
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpEditorTest.testPerformance",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpEditorTest.testPerformance1",
    ],
  },
  {
    label: "Formatter",
    projects: [
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testPerformance",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testPerformance1",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testPerformance2",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testPerformance3",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testPerformance4",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testForceIfBracesPerformance",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpFormatterTest.testPerformanceWi30101",
    ],
  },
  {
    label: "Typing",
    projects: [
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpTypingPerformanceTest.testTypingInHeredocLiteralWithVariables",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpTypingPerformanceTest.testTypingInHeredocLiteral",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpTypingPerformanceTest.testTypingInStringLiteral",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.PhpTypingPerformanceTest.testTypingInLargeMixedFile",
    ],
  },
  {
    label: "Completion",
    projects: [
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.completion.PhpAccessorsCompletionTest.testPerformanceOnLotsOfFields",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.completion.PhpTokenBasedCompletionTest.testPhpTokenBasedKeywordsCompletionPerformance",
    ],
  },
  {
    label: "Lang Util",
    projects: ["com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.lang.PhpLangUtilTest.testIsPhpIdentifierPerformance"],
  },
  {
    label: "Lexer",
    projects: ["com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.lexer.PhpLexerRegressionTest.testBackslashHeavyDoubleQuotedStringLexingPerformance"],
  },
  {
    label: "Phar File System",
    projects: [
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharTarGz",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharGz",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePhar",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharBz2",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharTar",
      "com.jetbrains.php.slowTests.PhpPerformanceTestSuite: com.jetbrains.php.phar.PharFileSystemPerformanceTest.testPerformancePharTarBz2",
    ],
  },
]
