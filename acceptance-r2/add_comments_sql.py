# -*- coding: utf-8 -*-
# SQL 注入（容错版）：跳过非 UTF-8（UTF-16 备份转储）文件
import os

ROOT = r"D:\Users\Lenovo\Desktop\shop\shop"

def trunc(s, n=110):
    s = s.strip()
    return s if len(s) <= n else s[:n] + "..."

def sql_desc(s):
    up = s.upper()
    if up.startswith('CREATE '): return '-- 建表/建对象：' + trunc(s)
    if up.startswith(('INSERT', 'REPLACE')): return '-- 插入数据：' + trunc(s)
    if up.startswith('UPDATE '): return '-- 更新数据：' + trunc(s)
    if up.startswith('SELECT'): return '-- 查询：' + trunc(s)
    if up.startswith('DELETE'): return '-- 删除数据：' + trunc(s)
    if up.startswith(('SET ', 'USE ')): return '-- 环境设置：' + trunc(s)
    if s.startswith((')', ';')): return '-- 语句结束'
    return '-- SQL：' + trunc(s)

total, skipped = 0, []
targets = []
for sub in ('deploy', 'shop', 'shop'):
    pass
for root, dirs, files in os.walk(os.path.join(ROOT, 'deploy', 'shop')):
    for f in files:
        if f.endswith('.sql'): targets.append(os.path.join(root, f))
for root, dirs, files in os.walk(os.path.join(ROOT, 'server', 'shop', 'migrations')):
    for f in files:
        if f.endswith('.sql'): targets.append(os.path.join(root, f))

for p in targets:
    try:
        lines = open(p, encoding='utf-8').read().splitlines()
    except UnicodeDecodeError:
        skipped.append(p); continue
    out, added = [], 0
    for ln in lines:
        st = ln.strip()
        skip = (not st) or st.startswith('--') or st.startswith('/*') or (out and out[-1].strip().startswith('--'))
        if not skip:
            out.append(sql_desc(st)); added += 1
        out.append(ln)
    open(p, 'w', encoding='utf-8', newline='').write('\n'.join(out) + '\n')
    total += added
    print(('%6d  %s' % (added, os.path.relpath(p, ROOT))))
print('TOTAL sql comments added:', total)
for s in skipped:
    print('SKIPPED non-utf8:', os.path.relpath(s, ROOT))
