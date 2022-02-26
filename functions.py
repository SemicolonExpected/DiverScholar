# -*- coding: utf-8 -*-
"""
Created on Fri Feb 25 21:09:57 2022

@author: SemicolonExpected
"""
import math
import json
import requests

def genderize(name):
    try:
        #batch query
        for i in range(math.floor(len(name)/10)+1):
            if ((i+1)*10) < len(name): 
                queryString = '&name[]='.join(name[i*10:(i+1)*10])
            else:
                queryString = '&name[]='.join(name[i*10:])
            url = "https://api.genderize.io?name[]="+queryString
            try:
                response = requests.get(url, timeout=10)
            except:
                return(-1)
            if response.status_code == 200:
                result = json.loads(response.text)
                [i.pop('count') for i in result]
                return result
            else:
                return(-1)
    except:
        #single query
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

print(genderize(["peter","lois","stewie","peter","lois","stewie","peter","lois"])) 