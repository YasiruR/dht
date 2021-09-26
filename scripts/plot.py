#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Sun Sep 26 09:16:38 2021

@author: yasi
"""

from matplotlib import pyplot as plt
import numpy as np
import math

keys_1 = ["(1,16)", "(1,100)", "(1,1000)"]
keys_2 = ["(2,16)", "(2,100)", "(2,1000)"]
keys_4 = ["(4,16)", "(4,100)", "(4,500)"]
keys_8 = ["(8,16)", "(8, 100)", "(8,500)"]
keys_16 = ["(16,16)", "(16,100)", "(16,500)"]


n_1_16_put = [0.4375, 0.3125, 0.3125]
n_1_100_put = [0.09, 0.12, 0.13]
n_1_1000_put = [0.049, 0.05, 0.067]

n_2_16_put = [0.375, 0.375, 0.4375]
n_2_100_put = [0.12, 0.12, 0.13]
n_2_1000_put = [0.055, 0.06, 0.067]

n_4_16_put = [0.625, 0.625, 0.5625]
n_4_100_put = [0.19, 0.2, 0.22]
n_4_500_put = [0.11, 0.118, 0.144]

n_8_16_put = [0.9375, 1, 0.875]
n_8_100_put = [0.28, 0.26, 0.29]
n_8_500_put = [0.122, 0.18, 0.12]

n_16_16_put = [1.4375, 1.4375, 1.5]
n_16_100_put = [0.38, 0.34, 0.4]
n_16_500_put = [0.134, 0.158, 0.146]

n_1_put = [np.mean(n_1_16_put), np.mean(n_1_100_put), np.mean(n_1_1000_put)]
n_1_put_err = [np.std(n_1_16_put) / math.sqrt(3), np.std(n_1_100_put) / math.sqrt(3), np.std(n_1_1000_put) / math.sqrt(3)]

n_2_put = [np.mean(n_2_16_put), np.mean(n_2_100_put), np.mean(n_2_1000_put)]
n_2_put_err = [np.std(n_2_16_put) / math.sqrt(3), np.std(n_2_100_put) / math.sqrt(3), np.std(n_2_1000_put) / math.sqrt(3)]

n_4_put = [np.mean(n_4_16_put), np.mean(n_4_100_put), np.mean(n_4_500_put)]
n_4_put_err = [np.std(n_4_16_put) / math.sqrt(3), np.std(n_4_100_put) / math.sqrt(3), np.std(n_4_500_put) / math.sqrt(3)]

n_8_put = [np.mean(n_8_16_put), np.mean(n_8_100_put), np.mean(n_8_500_put)]
n_8_put_err = [np.std(n_8_16_put) / math.sqrt(3), np.std(n_8_100_put) / math.sqrt(3), np.std(n_8_500_put) / math.sqrt(3)]

n_16_put = [np.mean(n_16_16_put), np.mean(n_16_100_put), np.mean(n_16_500_put)]
n_16_put_err = [np.std(n_16_16_put) / math.sqrt(3), np.std(n_16_100_put) / math.sqrt(3), np.std(n_16_500_put) / math.sqrt(3)]

#define chart 
fig_put, ax = plt.subplots()

#create chart
ax.bar(x=keys_1, #x-coordinates of bars
       height=n_1_put, #height of bars
       yerr=n_1_put_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='seagreen',
       label='1-node') 

ax.bar(x=keys_2, #x-coordinates of bars
       height=n_2_put, #height of bars
       yerr=n_2_put_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='royalblue',
       label='2-nodes') 

ax.bar(x=keys_4, #x-coordinates of bars
       height=n_4_put, #height of bars
       yerr=n_4_put_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='grey',
       label='4-nodes') 

ax.bar(x=keys_8, #x-coordinates of bars
       height=n_8_put, #height of bars
       yerr=n_8_put_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='sienna',
       label='8-nodes')

ax.bar(x=keys_16, #x-coordinates of bars
       height=n_16_put, #height of bars
       yerr=n_16_put_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='mediumpurple',
       label='16-nodes') 

ax.set_ylabel('average time (s)')
ax.set_xlabel('experiment (number of nodes, number of requests)')
plt.legend()
fig_put.set_size_inches(12, 12)
plt.show()
fig_put.savefig('put.pdf', bbox_inches="tight")


n_1_16_get = [0.25, 0.25, 0.3125]
n_1_100_get = [0.12, 0.12, 0.08]
n_1_1000_get = [0.053, 0.057, 0.064]

n_2_16_get = [0.375, 0.375, 0.375]
n_2_100_get = [0.11, 0.13, 0.14]
n_2_1000_get = [0.057, 0.064, 0.068]

n_4_16_get = [0.5625, 0.5625, 0.625]
n_4_100_get = [0.19, 0.18, 0.2]
n_4_500_get = [0.084, 0.12, 0.098]

n_8_16_get = [0.9375, 0.9375, 0.875]
n_8_100_get = [0.23, 0.25, 0.26]
n_8_500_get = [0.148, 0.158, 0.19]

n_16_16_get = [1.4375, 1.4375, 1.75]
n_16_100_get = [0.39, 0.35, 0.41]
n_16_500_get = [0.156, 0.158, 0.194]

n_1_get = [np.mean(n_1_16_get), np.mean(n_1_100_get), np.mean(n_1_1000_get)]
n_1_get_err = [np.std(n_1_16_get) / math.sqrt(3), np.std(n_1_100_get) / math.sqrt(3), np.std(n_1_1000_get) / math.sqrt(3)]

n_2_get = [np.mean(n_2_16_get), np.mean(n_2_100_get), np.mean(n_2_1000_get)]
n_2_get_err = [np.std(n_2_16_get) / math.sqrt(3), np.std(n_2_100_get) / math.sqrt(3), np.std(n_2_1000_get) / math.sqrt(3)]

n_4_get = [np.mean(n_4_16_get), np.mean(n_4_100_get), np.mean(n_4_500_get)]
n_4_get_err = [np.std(n_4_16_get) / math.sqrt(3), np.std(n_4_100_get) / math.sqrt(3), np.std(n_4_500_get) / math.sqrt(3)]

n_8_get = [np.mean(n_8_16_get), np.mean(n_8_100_get), np.mean(n_8_500_get)]
n_8_get_err = [np.std(n_8_16_get) / math.sqrt(3), np.std(n_8_100_get) / math.sqrt(3), np.std(n_8_500_get) / math.sqrt(3)]

n_16_get = [np.mean(n_16_16_get), np.mean(n_16_100_get), np.mean(n_16_500_get)]
n_16_get_err = [np.std(n_16_16_get) / math.sqrt(3), np.std(n_16_100_get) / math.sqrt(3), np.std(n_16_500_get) / math.sqrt(3)]


#define chart 
fig_get, ax = plt.subplots()

#create chart
ax.bar(x=keys_1, #x-coordinates of bars
       height=n_1_get, #height of bars
       yerr=n_1_get_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='seagreen',
       label='1-node') 

ax.bar(x=keys_2, #x-coordinates of bars
       height=n_2_get, #height of bars
       yerr=n_2_get_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='royalblue',
       label='2-nodes') 

ax.bar(x=keys_4, #x-coordinates of bars
       height=n_4_get, #height of bars
       yerr=n_4_get_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='grey',
       label='4-nodes') 

ax.bar(x=keys_8, #x-coordinates of bars
       height=n_8_get, #height of bars
       yerr=n_8_get_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='sienna',
       label='8-nodes')

ax.bar(x=keys_16, #x-coordinates of bars
       height=n_16_get, #height of bars
       yerr=n_16_get_err, #error bar width
       capsize=8, #length of error bar caps
       ecolor='red',
       color='mediumpurple',
       label='16-nodes') 

ax.set_ylabel('average time (s)')
ax.set_xlabel('experiment (number of nodes, number of requests)')
plt.legend()
fig_get.set_size_inches(12, 12)
plt.show()
fig_get.savefig('get.pdf', bbox_inches="tight")


keys_a = [16, 100, 1000]
keys_b = [16, 100, 500]

plt.plot(keys_a, n_1_get, label="1-node", color="seagreen")
plt.plot(keys_a, n_2_get, label="2-nodes", color="royalblue")
plt.plot(keys_b, n_4_get, label="4-nodes", color="grey")
plt.plot(keys_b, n_8_get, label="8-nodes", color="sienna")
plt.plot(keys_b, n_16_get, label="16-nodes", color="mediumpurple")
plt.legend()
plt.xlabel('number of requests')
plt.ylabel('average time (s)')
plt.savefig('get_line.pdf', bbox_inches="tight")
plt.show()


plt.plot(keys_a, n_1_put, label="1-node", color="seagreen")
plt.plot(keys_a, n_2_put, label="2-nodes", color="royalblue")
plt.plot(keys_b, n_4_put, label="4-nodes", color="grey")
plt.plot(keys_b, n_8_put, label="8-nodes", color="sienna")
plt.plot(keys_b, n_16_put, label="16-nodes", color="mediumpurple")
plt.legend()
plt.xlabel('number of requests')
plt.ylabel('average time (s)')
plt.savefig('put_line.pdf', bbox_inches="tight")
plt.show()



























