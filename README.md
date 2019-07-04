# Fossil SCM delta compression algorithm
======================================

Format:
http://www.fossil-scm.org/index.html/doc/tip/www/delta_format.wiki

Algorithm:
http://www.fossil-scm.org/index.html/doc/tip/www/delta_encoder_algorithm.wiki

Original implementation:
http://www.fossil-scm.org/index.html/artifact/d1b0598adcd650b3551f63b17dfc864e73775c3d


Usage
-----

### Create(origin string, target string) []rune

Returns a delta (as `Array` of runes) from origin to target.

### Apply(origin string, delta []rune, verifyChecksum bool) ([]rune, error)

Returns target (as `Array` of runes) by applying delta to origin.

Return error if it fails to apply the delta
(e.g. if it was corrupted).

Argument `verifyChecksum` is use disable checksum verification.

### OutputSize(delta []rune) (int, error)

Returns a size of target for this delta.

Return error if it can't read the size from delta.


#LICENSE
-------

Copyright 2019 Phat Tran (Golang port)
Copyright 2014 Dmitry Chestnykh (JavaScript port)
Copyright 2007 D. Richard Hipp  (original C version)
All rights reserved.

Redistribution and use in source and binary forms, with or
without modification, are permitted provided that the
following conditions are met:

  1. Redistributions of source code must retain the above
     copyright notice, this list of conditions and the
     following disclaimer.

  2. Redistributions in binary form must reproduce the above
     copyright notice, this list of conditions and the
     following disclaimer in the documentation and/or other
     materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE AUTHORS ``AS IS'' AND ANY EXPRESS
OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHORS OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation
are those of the authors and contributors and should not be interpreted
as representing official policies, either expressed or implied, of anybody
else.
