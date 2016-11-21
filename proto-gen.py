#!/usr/bin/env python3


import os
import subprocess

for root, dirs, files in os.walk('.'):
    if ".git" in root or ".idea" in root:
        continue

    if True in map(lambda f:f.endswith(".proto"), files):
        subprocess.run("protoc {DIR}/*.proto --go_out={DIR}".format(DIR=root), shell=True)
