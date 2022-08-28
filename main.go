/*
 * Copyright (c) 2022, Dana Burkart <dana.burkart@gmail.com>
 *
 * SPDX-License-Identifier: BSD-2-Clause
 */

package main

import "fmt"

func main() {
	s := BeginScanning("Foo", `  /**
	   * Creates a <code>ReadAheadInputStream</code> with the specified buffer size and read-ahead
	   * threshold
	   *
	   * @param inputStream The underlying input stream.
	   * @param bufferSizeInBytes The buffer size.
	   */
	  public ReadAheadInputStream(
		  InputStream inputStream, int bufferSizeInBytes) {
		Preconditions.checkArgument(bufferSizeInBytes > 0,
			"bufferSizeInBytes should be greater than 0, but the value is " + bufferSizeInBytes);
		activeBuffer = ByteBuffer.allocate(bufferSizeInBytes);
		readAheadBuffer = ByteBuffer.allocate(bufferSizeInBytes);
		this.underlyingInputStream = inputStream;
		activeBuffer.flip();
		readAheadBuffer.flip();
	  }`)

	d := ParseDocument(s, "Foo/Bar/Baz")
	fmt.Println(d)
}
