# -*- coding: utf-8 -*-
"""
Created on Fri Feb 25 21:09:57 2022

@author: SemicolonExpected
"""

import json
import requests

def genderize(name):
    url = "https://api.genderize.io?name="+name
    try:
        response = requests.get(url, timeout=10)
    except:
        return(-1)
    if response.status_code == 200:
        result = json.loads(response.text)
        result.pop('count')
        return result
    else:
        return(-1)
