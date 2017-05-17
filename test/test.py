#!/usr/bin/env python
# coding=utf-8
import yaml
import json

print json.dumps(yaml.load(file("./test.yaml").read()), indent=4)
