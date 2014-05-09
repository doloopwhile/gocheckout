go-checkout
===========

Checkout specific revision of go package in GOPATH.
Supports git, hg and bzr.

Installation
------------
    go get github.com/doloopwhile/gocheckout

Usage
-----
    gocheckout [options] <package_name> <revision>

For example, if your product require version to build version 0.1 of [martini](https://github.com/go-martini/martini)

    go get -d github.com/go-martini/martini

    gocheckout github.com/go-martini/martini v0.1

Author
------
OMOTO Kenji doloopwhile@gmail.com

Some codes in gocheckout is picked from [gom](https://github.com/mattn/gom)
by Yasuhiro Matsumoto mattn.jp@gmail.com.
