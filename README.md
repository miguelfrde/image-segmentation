image-segmentation
==================

[Try it online](http://image-segmentation.herokuapp.com/)

Image segmentation using Minimum Spanning Forests.

Implements:

- [Efficient Graph-Based Image Segmentation (GBS)](http://cs.brown.edu/~pff/papers/seg-ijcv.pdf)
- [An Efficient Parallel Algorithm for Graph-Based Image Segmentation (PHMSF)](http://algo2.iti.kit.edu/wassenberg/wassenberg09parallelSegmentation.pdf)

Also, as a helper for the second segmentation algorithm:

- [Block-based noise estimation using adaptive Gaussian filtering](http://ieeexplore.ieee.org/xpl/login.jsp?tp=&arnumber=1405723&url=http%3A%2F%2Fieeexplore.ieee.org%2Fxpls%2Fabs_all.jsp%3Farnumber%3D1405723)

Developed as a project for the course "Graph Theory, Networks and Applications" at [Mälardalen University](http://mdh.se/). The course's project report can be found [here](https://www.dropbox.com/s/gdtghavyyr1x7m1/report.pdf?dl=0).


## Setup

Prerequisites:

- [Go](https://golang.org/)
- [Node](http://nodejs.org/)
- [Bower](http://bower.io/)

```
$ go get github.com/miguelfrde/image-segmentation
$ cd $GODIR/src/github.com/miguelfrde/image-segmentation
$ bower install
$ mkdir tmp
```

The `tmp` directory needs to be created since result images will be served from it in Heroku.


## Run

```
$ PORT=8080 go run main.go
```

Go to `localhost:5000` and try some images!

## Test

```
$ go test ./...
```


## TODO

- Implement the parallel version of the PHMSF algorithm
- Implement more MST or non-graph based segmentation algorithms
