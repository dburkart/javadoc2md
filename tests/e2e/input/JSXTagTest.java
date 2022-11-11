/**
 * This test ensures that JSX tags that are not closed are closed after
 * transpilation. JSX does not have the same facilities for closing tags
 * as HTML does.
 *
 */
public class JSXTagTest {

    /**
     * This test ensures that we close a p tag.
     * <p>
     * The above tag should be closed upon transpilation.
     */
    public void testUnclosedPTag() {

    }

    /**
     * <p>
     * This test ensures that we don't close tags which are actually
     * closed already.
     * </p>
     */
    public void testTagsThatAreClosedOnSeparateLines() {

    }

    /**
     * This <i>word</i> is italicized
     * <i>This</i> one is <i>too</i>
     */
    public void testItalics() {

    }

    /**
     * <a href="https://tools.ietf.org/html/rfc7231#section-4.3">RFC-7231</a>
     */
    public void testATag() {

    }

    /**
     * This tests that the pre tag gets converted properly
     *
     * <pre>
     * int foo = bar;
     * foo * 2;
     * </pre>
     */
    public void testPreTag() {

    }

    /**
     * This tests that the code tag gets converted properly
     *
     * <code>
     * int foo = bar;
     * foo * 2;
     * </code>
     *
     * <code>
     * int foo = bar;
     * </code>
     */
    public void testCodeTag() {

    }

    /**
     * This tests an inline <code>code</code> tag.
     */
    public void testInlineCodeTag() {

    }
}