# Build Evnironment using Windows mingw, vcpkg

export TRIPLET=x64-mingw-dynamic

export VCPKG_DEFAULT_TRIPLET=$TRIPLET
export VCPKG_DEFAULT_HOST_TRIPLET=$TRIPLET
export VCPKG_BASE=$HOME/Documents/Projekte/vcpkg

export PATH=$HOME/bin
export PATH=$PATH:/c/WINDOWS/system32
export PATH=$PATH:/c/WINDOWS/System32/Wbem
export PATH=$PATH:/c/WINDOWS/System32/WindowsPowerShell/v1.0
export PATH=$PATH:/C/msys64/mingw64/bin
export PATH=$PATH:/usr/local/bin:/usr/bin:/bin:/mingw64/bin

export PATH="$PATH:/c/Program Files/Go/bin"

export PATH=$PATH:$VCPKG_BASE
export PATH=$PATH:$VCPKG_BASE/installed/$TRIPLET/bin

#export CGO_CXXFLAGS="--std=c++11 --enable-threads=posix"
export CGO_CXXFLAGS="--std=c++11 -lpthread -pthread"
export CGO_CFLAGS='-m64 -O2 -pthread'
export CGO_ENABLED=1
export CGO_CPPFLAGS="-I$VCPKG_BASE/installed/$TRIPLET/include"
export CGO_LDFLAGS="-L$VCPKG_BASE/installed/$TRIPLET/lib"

export TESSDATA_PREFIX=$VCPKG_BASE/installed/x64-mingw-dynamic/share/tessdata

vcpkg install mchehab-zbar
vcpkg install tesseract
vcpkg install opencv
vcpkg install opencv[contrib]

# -DWITH_OBSENSOR=OFF

nano /C/Users/thomas.hoffmann/Documents/Projekte/vcpkg/scripts/buildsystems
# add `set(WITH_OBSENSOR OFF)`
# -buildmode=c-shared

go build -ldflags "-s -w --enable-threads=posix"

export cwd=`pwd`
cd $VCPKG_BASE/installed/$TRIPLET/lib
cp libopencv_core4.dll.a libopencv_core490.dll.a 
cp libopencv_calib3d4.dll.a libopencv_calib3d490.dll.a
cp libopencv_core490.dll.a libopencv_core490.dll.a
cp libopencv_dnn4.dll.a libopencv_dnn490.dll.a
cp libopencv_features2d4.dll.a libopencv_features2d490.dll.a
cp libopencv_flann4.dll.a libopencv_flann490.dll.a
cp libopencv_highgui4.dll.a libopencv_highgui490.dll.a
cp libopencv_imgcodecs4.dll.a libopencv_imgcodecs490.dll.a
cp libopencv_imgproc4.dll.a libopencv_imgproc490.dll.a
cp libopencv_ml4.dll.a libopencv_ml490.dll.a
cp libopencv_objdetect4.dll.a libopencv_objdetect490.dll.a
cp libopencv_photo4.dll.a libopencv_photo490.dll.a
cp libopencv_stitching4.dll.a libopencv_stitching490.dll.a
cp libopencv_video4.dll.a libopencv_video490.dll.a
cp libopencv_videoio4.dll.a libopencv_videoio490.dll.a
cp libleptonica-1.84.1.dll.a libleptonica.dll.a
cp libtesseract53.dll.a libtesseract.dll.a
cd $cwd


cd $VCPKG_BASE/installed/$TRIPLET/lib
cp libopencv_face4.dll.a libopencv_face490.dll.a
cp libopencv_xfeatures2d4.dll.a libopencv_xfeatures2d490.dll.a
cp libopencv_plot4.dll.a libopencv_plot490.dll.a
cp libopencv_tracking4.dll.a libopencv_tracking490.dll.a
cp libopencv_img_hash4.dll.a libopencv_img_hash490.dll.a
cp libopencv_bgsegm4.dll.a libopencv_bgsegm490.dll.a
cp libopencv_aruco4.dll.a libopencv_aruco490.dll.a
cp libopencv_wechat_qrcode4.dll.a libopencv_wechat_qrcode490.dll.a
cp libopencv_ximgproc4.dll.a libopencv_ximgproc490.dll.a
cd $cwd



