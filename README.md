# Fossil SCM delta compression algorithm
======================================

Format:
http://www.fossil-scm.org/index.html/doc/tip/www/delta_format.wiki

Algorithm:
http://www.fossil-scm.org/index.html/doc/tip/www/delta_encoder_algorithm.wiki

Original implementation:
http://www.fossil-scm.org/index.html/artifact/d1b0598adcd650b3551f63b17dfc864e73775c3d

> This is a port from the original C and Javascript implementation. See references below.

Other implementations:

- [Haxe](https://github.com/endel/fossil-delta-hx)
- [Python](https://github.com/ggicci/python-fossil-delta)
- [JavaScript](https://github.com/dchest/fossil-delta-js) ([Online demo](https://dchest.github.io/fossil-delta-js/))

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

