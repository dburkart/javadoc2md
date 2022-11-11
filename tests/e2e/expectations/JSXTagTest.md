# JSXTagTest

## Definition

```java
public class JSXTagTest
```

## Overview

This test ensures that JSX tags that are not closed are closed after
transpilation. JSX does not have the same facilities for closing tags
as HTML does.

### `public void testUnclosedPTag()` {#testUnclosedPTag()}

This test ensures that we close a p tag.
<p/>
The above tag should be closed upon transpilation.

### `public void testTagsThatAreClosedOnSeparateLines()` {#testTagsThatAreClosedOnSeparateLines()}

<p>
This test ensures that we don't close tags which are actually
closed already.
</p>

### `public void testItalics()` {#testItalics()}

This <i>word</i> is italicized
<i>This</i> one is <i>too</i>

### `public void testATag()` {#testATag()}

[RFC-7231](https://tools.ietf.org/html/rfc7231#section-4.3)

### `public void testPreTag()` {#testPreTag()}

This tests that the pre tag gets converted properly

```java
int foo = bar;
foo * 2;
```

### `public void testCodeTag()` {#testCodeTag()}

This tests that the code tag gets converted properly

```java
int foo = bar;
foo * 2;
```

```java
int foo = bar;
```

### `public void testInlineCodeTag()` {#testInlineCodeTag()}

This tests an inline `code` tag.

