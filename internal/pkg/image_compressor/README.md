# Инструкция по настройке компрессора

### Обновление системы
```
sudo apt update
sudo apt upgrade
```

### Установка либ/драйверов для работы с изображениями
```
sudo apt-get install libjpeg-dev
sudo apt-get install libpng-dev
sudo apt-get install libwebp-dev
```

### Установка самой либы для ковертера
Качаем исходники отсюда [тык](https://github.com/libvips/libvips/releases). Затем команды
```
tar xf vips-8.12.0-rc1.tar.gz
cd vips-8.12.0/
./configure --with-png-includes --with-png-libraries --with-jpeg-includes --with-jpeg-libraries
```

Напртив JPEG, PNG и WEBP должны быть **yes**

```
make
sudo make install
sudo ldconfig
```

### Установка либы проекта

```
go get -u github.com/h2non/bimg
```