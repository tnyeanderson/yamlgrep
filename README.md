# yamlgrep

```bash
go install github.com/tnyeanderson/yamlgrep@latest
```

Pipe a multidoc YAML to this command and give it grep-ish arguments, and only
the matching documents from the YAML will be printed. It's like `grep -C`, but
with the perfect amount of context.

Take the following `demo.yaml`:

```yaml
---
id: 1
first_name: Kin
last_name: Stockwell
email: kstockwell0@cocolog-nifty.com
---
id: 2
first_name: Serene
last_name: Edgeller
email: sedgeller1@google.com.br
---
id: 3
first_name: Hollis
last_name: Syrie
email: hsyrie2@feedburner.com
```

Let's find all the YAML documents that contain `ell`:

```
$ cat /tmp/demo.yaml | yamlgrep ell
---
id: 1
first_name: Kin
last_name: Stockwell
email: kstockwell0@cocolog-nifty.com
---
id: 2
first_name: Serene
last_name: Edgeller
email: sedgeller1@google.com.br
```

Under the hood, `yamlgrep` calls whatever grep is in the user's `$PATH`
(usually the system's). This means you can use all the same arguments you
normally would with grep. However, arguments like `-C` are ignored, because
this program only cares if the grep matched or not, and it ignores grep's
output (both to stdout and stderr).

## History and performance

This started as a simple bash script, which I thought was really elegant.
I wrote it quickly after being frustrated with YAML searching tools at work,
and it all worked fine with small YAML files at home, but once I tried using it
on a very large YAML response, I learned just how INCREDIBLY SLOW bash scripts
can be, even if they seem so nice. I was waiting *minutes* for my elegant
script to complete.

So I rewrote it in go and it became **1.5 billion times faster**. That is not
an exaggeration: the original script is retained in the `testdata` directory
and the benchmarks are available for you to run yourself:

```bash
$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/tnyeanderson/yamlgrep
cpu: AMD Ryzen 5 2600 Six-Core Processor
BenchmarkYAMLGrep-12         	1000000000	         0.7599 ns/op
BenchmarkYAMLGrep_Bash-12    	       1	1172747556 ns/op
PASS
ok  	github.com/tnyeanderson/yamlgrep	37.227s
```
