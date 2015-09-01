#!/bin/sh

sudo docker run -d --name sitehub -v /home/user/devel/sitehub:/var/sitehub -p 8080:8080 prinsmike/sitehub:0.1
