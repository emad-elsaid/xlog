![](/docs/public/shortcode.png)
#extension

Xlog extension **shortcode** allow using blocks with custom language code that can render the content of the block with custom function

For example rendering an alert can use two different formats

# Short format

<pre>
/alert this is important
</pre>

/alert this is important

# Long format

<pre>
```alert
this is important
```
</pre>

```alert
this is important
```

# Default blocks

Shortcode extension includes couple default blocks:

## /alert

```alert
Computer science is the study of computation, automation, and information. Computer science spans theoretical disciplines (such as algorithms, theory of computation, information theory, and automation) to practical disciplines (including the design and implementation of hardware and software). Computer science is generally considered an area of academic research and distinct from computer programming.
```

## /info

```info
Computer science is the study of computation, automation, and information. Computer science spans theoretical disciplines (such as algorithms, theory of computation, information theory, and automation) to practical disciplines (including the design and implementation of hardware and software). Computer science is generally considered an area of academic research and distinct from computer programming.
```

## /success

```success
Computer science is the study of computation, automation, and information. Computer science spans theoretical disciplines (such as algorithms, theory of computation, information theory, and automation) to practical disciplines (including the design and implementation of hardware and software). Computer science is generally considered an area of academic research and distinct from computer programming.
```

## /warning

```warning
Computer science is the study of computation, automation, and information. Computer science spans theoretical disciplines (such as algorithms, theory of computation, information theory, and automation) to practical disciplines (including the design and implementation of hardware and software). Computer science is generally considered an area of academic research and distinct from computer programming.
```
