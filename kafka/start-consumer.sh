#!/bin/bash

cd ./consumer
cd ./starter
./starter --t_endpoint=192.168.49.2:30880

cd ../worker
./worker --t_endpoint=192.168.49.2:30880 --k_endpoint=192.168.49.2:32059
