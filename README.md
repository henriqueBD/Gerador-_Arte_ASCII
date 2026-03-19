# Go ASCII Art Converter

A command-line utility written in **Go (Golang)** that transforms conventional PNG images into ASCII character-based compositions, generating a new image as a result.

This project was developed to deepen knowledge in digital image processing, file manipulation, and the efficient use of the language's standard library.

## Demonstration

| Before (Original Image) | After (ASCII Result) |
| :---: | :---: |
| <img src="images/pawn.png" width="300" alt="Original image before conversion"> | <img src="images/pawn_ASCII.png" width="300" alt="Image converted to ASCII"> |
| <img src="images/chess.png" width="300" alt="Original image before conversion"> | <img src="images/chess_ASCII.png" width="300" alt="Image converted to ASCII"> |

## Technical Skills

  * **File Manipulation and I/O:** Safe reading and writing of binary files (PNG images) using the `os` package.
  * **Native Image Processing:** Extensive use of Go's `image` library (without third-party external dependencies) for decoding, creating bounding boxes, and drawing (`image/draw`).
  * **Algorithms and Mathematical Logic:**
      * **Grayscale:** Conversion of colored pixels to luminosity (Grayscale).
      * **Downscaling (Nearest Neighbor):** Manual implementation of the nearest neighbor algorithm for image resizing.
      * **Color Quantization:** Mapping pixel intensities to dynamically match a limited character palette (ASCII blocks).
  * **Best Practices:** Code modularized into single-responsibility functions and error handling (safeguards for array boundaries).

## How to run

1.  Make sure you have [Go installed](https://go.dev/doc/install) on your machine.
2.  Clone this repository.
3.  Place the image you want to convert (in `.png` format) into the `images/` folder.
4.  Run the program:

<!-- end list -->

```bash
go run main.go
```

5.  Enter the name of the image (without the extension) when prompted by the terminal.
6.  The result will be saved in the `images/` folder with the `_ASCII` suffix.
