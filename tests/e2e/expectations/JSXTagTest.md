# JSXTagTest

## Definition

```java
public class JSXTagTest
```

## Overview

This test ensures that JSX tags that are not closed are closed after
transpilation. JSX does not have the same facilities for closing tags
as HTML does.

### `public void testUnclosedPTag()` {#testUnclosedPTag}

This test ensures that we close a p tag.
<p/>
The above tag should be closed upon transpilation.

### `public void testTagsThatAreClosedOnSeparateLines()` {#testTagsThatAreClosedOnSeparateLines}

<p>
This test ensures that we don't close tags which are actually
closed already.
</p>

