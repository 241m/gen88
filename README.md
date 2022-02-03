# gen88

Tool to generate digital art from SHA256 hashes of public domain works.

1. Run `./bin/gutenberg.sh NAME URL` to download a digital work from Project
   Gutenberg. The URL must be a link to the work's plain text UTF-8 format.

2. Run `./bin/generate.sh` to generate art in SVG format.

The art is generated using symbols from SVG files in the `symbols` directory.
