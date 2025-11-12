pkgname=pydev
pkgver=0.0.1
pkgrel=2
pkgdesc="pydev"
arch=('x86_64')
url="https://github.com/nnovak-pylon/pydev"
makedepends=('git')
license=('none')
source=("git+ssh://git@github.com/nnovak-pylon/${pkgname}.git")
sha256sums=('SKIP')

build() {
  cd "$srcdir/$pkgname"
  go build .
}

package() {
  mkdir -pm755 "$pkgdir/usr/bin"
  install -Dm755 "$srcdir/$pkgname/$pkgname" "$pkgdir/usr/local/bin/$pkgname"
}
