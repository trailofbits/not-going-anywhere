This entire example is from a [go gotchas](https://yourbasic.org/golang/gotcha-append/) blog. The caveat at the end of their blog is so important, I listed it in the README too, just in case ;)

```
The scary case: It “worked” for me
In some Go implementations []byte("ba") only allocates two bytes, and then the code seems to work: the first string is "bad" and the second one "bag".

Unfortunately the code is still wrong, even though it seems to work. The program may behave differently when you run it in another environment.
```
