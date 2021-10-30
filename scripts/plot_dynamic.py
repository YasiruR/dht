#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Mon Oct 25 13:50:08 2021

@author: yasi
"""

from matplotlib import pyplot as plt
import numpy as np
import math

keys_join=['join-10', 'join-20', 'join-30', 'join-40', 'join-50']
keys_leave=['leave-5', 'leave-10', 'join-15', 'leave-20', 'leave-25']
keys_crash=['1', '2', '3', '5', '10']

join_10=[0.236247, 0.237322, 0.22953]
join_20=[0.579841, 0.584309, 0.567491]
join_30=[0.959497, 1.098406, 0.983137]
join_40=[1.415455, 1.447214, 1.431279]
join_50=[1.869647, 1.877155, 1.85583]

leave_5=[0.070959, 0.072373, 0.073]
leave_10=[0.149947, 0.160596, 0.159415]
leave_15=[0.228879, 0.255471, 0.273161]
leave_20=[0.323369, 0.304179, 0.324312]
leave_25=[0.405809, 0.405604, 0.413812]

crash_1=[5.049562, 5.047081, 5.04918]
crash_2=[5.611188, 5.527882, 5.440613]
crash_3=[5.537222, 5.464855, 5.466797]
crash_5=[5.57084, 5.927389, 5.915173]
crash_10=[5.772832, 5.75469, 5.812839]

join=[np.mean(join_10), np.mean(join_20), np.mean(join_30), np.mean(join_40), np.mean(join_50)]
join_err=[np.std(join_10) / math.sqrt(3), np.std(join_20) / math.sqrt(3), np.std(join_30) / math.sqrt(3), np.std(join_40) / math.sqrt(3), np.std(join_50) / math.sqrt(3), ]

leave=[np.mean(leave_5), np.mean(leave_10), np.mean(leave_15), np.mean(leave_20), np.mean(leave_25)]
leave_err=[np.std(leave_5) / math.sqrt(3), np.std(leave_10) / math.sqrt(3), np.std(leave_15) / math.sqrt(3), np.std(leave_20) / math.sqrt(3), np.std(leave_25) / math.sqrt(3), ]

crash=[np.mean(crash_1), np.mean(crash_2), np.mean(crash_3), np.mean(crash_5), np.mean(crash_10)]
crash_err=[np.std(crash_1) / math.sqrt(3), np.std(crash_2) / math.sqrt(3), np.std(crash_3) / math.sqrt(3), np.std(crash_5) / math.sqrt(3), np.std(crash_10) / math.sqrt(3), ]

#define chart 
fig_join_leave, ax = plt.subplots()

#create chart
ax.bar(x=keys_join, #x-coordinates of bars
       height=join, #height of bars
       yerr=join_err, #error bar width
       capsize=25, #length of error bar caps
       ecolor='red',
       color='seagreen',
       label='join') 

ax.bar(x=keys_leave, #x-coordinates of bars
       height=leave, #height of bars
       yerr=leave_err, #error bar width
       capsize=25, #length of error bar caps
       ecolor='red',
       color='royalblue',
       label='leave') 

ax.set_ylabel('average time (s)')
ax.set_xlabel('experiment (type and number of nodes)')
plt.legend()
fig_join_leave.set_size_inches(12, 7)
plt.show()
fig_join_leave.savefig('join_and_leave_bar.pdf', bbox_inches="tight")

#define chart 
fig_crash, ax = plt.subplots()

ax.bar(x=keys_crash, #x-coordinates of bars
       height=crash, #height of bars
       yerr=crash_err, #error bar width
       capsize=25, #length of error bar caps
       ecolor='red',
       color='grey',
       width=0.6)

ax.set_ylabel('corrected average time (s)')
ax.set_xlabel('experiment (number of nodes crashed)')
#fig_crash.set_size_inches(6, 4)
plt.show()
fig_crash.savefig('crash_bar.pdf', bbox_inches="tight")


