#!/usr/bin/env python3

import urllib.request
import os

try:
	os.mkdir('downloads')
except:
	pass

with open('sources.txt') as f:
	for line in f.readlines():
		url = line.strip()
		x = url.split('/')
		filename = f'{x[3]}_{x[4]}_{x[-1]}'
		print(f'downloading {filename}')

		response = urllib.request.urlopen(url)
		data = response.read()
		with open(f'downloads/{filename}', 'wb+') as src:
			src.write(data)
