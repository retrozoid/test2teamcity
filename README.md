```
go install github.com/retrozoid/test2teamcity/cmd@latest
```

```
% ~/go/bin/cmd go test -v
##teamcity[testStarted name='TestXxx' flowId='1' captureStandardOutput='false']
##teamcity[testStdOut name='TestXxx' flowId='1' out='test_test.go:10: hello']
##teamcity[testFinished name='TestXxx' flowId='1' duration='500']
##teamcity[testStdOut name='TestXxx' flowId='1' out='PASS']
##teamcity[testStdOut name='TestXxx' flowId='1' out='ok         github.com/retrozoid/test2teamcity      1.485s']
```