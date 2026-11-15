# -*- coding: utf-8 -*-
# 逐行注释注入脚本：为业务代码每一行上方补齐注释（语法感知，不破坏代码）
# 跳过规则：空行、已有注释的行（上行为注释/本行以 // /* * <!-- -- 开头）、字符串/原始字面量内部、多行标签内部
import os, re, sys

ROOT = r"D:\Users\Lenovo\Desktop\shop\shop"

def trunc(s, n=110):
    s = s.strip()
    return s if len(s) <= n else s[:n] + "..."

def go_desc(s):
    # 按语句首词生成中文说明
    if s.startswith(('if ', 'if(')): return '// 条件判断：' + trunc(s)
    if s.startswith('} else') or s.startswith('else'): return '// 条件分支：' + trunc(s)
    if s.startswith('for '): return '// 循环处理：' + trunc(s)
    if s.startswith(('switch ', 'switch{', 'select')): return '// 多路分支：' + trunc(s)
    if s.startswith('case '): return '// 匹配分支：' + trunc(s)
    if s.startswith('default:'): return '// 默认分支'
    if s.startswith('return'): return '// 返回：' + trunc(s)
    if s.startswith('func '): return '// 函数定义：' + trunc(s)
    if s.startswith('type '): return '// 类型定义：' + trunc(s)
    if s.startswith(('var ', 'const ')): return '// 声明：' + trunc(s)
    if s.startswith('package '): return '// 包声明'
    if s.startswith('import') or s == 'import (': return '// 导入依赖包'
    if s.startswith(('go ', 'defer ')): return '// 异步/延迟执行：' + trunc(s)
    if s.startswith(('}', ')', ']')): return '// 结束当前代码块'
    if s.startswith(('break', 'continue', 'panic', 'fallthrough')): return '// 流程控制：' + trunc(s)
    if '=' in s.split('(')[0]: return '// 赋值/更新：' + trunc(s)
    return '// 执行：' + trunc(s)

def js_desc(s):
    if s.startswith(('if', 'if(')): return '// 条件判断：' + trunc(s)
    if s.startswith(('} else', 'else')): return '// 条件分支：' + trunc(s)
    if s.startswith('for'): return '// 循环处理：' + trunc(s)
    if s.startswith(('switch', 'case ')): return '// 分支处理：' + trunc(s)
    if s.startswith('return'): return '// 返回：' + trunc(s)
    if s.startswith(('function', 'const ', 'let ', 'var ')): return '// 声明：' + trunc(s)
    if s.startswith('}'): return '// 结束当前代码块'
    if re.match(r'^\w+\s*\(', s): return '// 调用：' + trunc(s)
    return '// 执行：' + trunc(s)

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

def is_comment_line(s, lang):
    if lang == 'go': return s.startswith('//') or s.startswith('/*') or s.startswith('*')
    if lang in ('js', 'vue_script'): return s.startswith('//') or s.startswith('/*') or s.startswith('*')
    if lang == 'html': return s.startswith('<!--')
    if lang == 'sql': return s.startswith('--')
    return False

def inject_go(path):
    lines = open(path, encoding='utf-8').read().splitlines()
    out, added = [], 0
    in_raw = False       # 反引号原始字符串
    in_block = False     # /* */ 块注释
    for ln in lines:
        st = ln.strip()
        if in_block:
            out.append(ln)
            if '*/' in ln: in_block = False
            continue
        if in_raw:
            out.append(ln)
            if '`' in ln: in_raw = False
            continue
        if st.startswith('/*') and '*/' not in st:
            in_block = True; out.append(ln); continue
        if '`' in ln:
            # 本行含反引号：按奇偶决定是否进入原始字符串
            if ln.count('`') % 2 == 1: in_raw = True
            # 反引号行本身可能也是语句，但如果上一行不是注释则补注释
            if st and not is_comment_line(st, 'go') and not (out and out[-1].strip().startswith('//')):
                out.append(go_desc(st)); added += 1
            out.append(ln); continue
        skip = (not st) or is_comment_line(st, 'go') or (out and out[-1].strip().startswith('//'))
        if not skip:
            out.append(go_desc(st)); added += 1
        out.append(ln)
    open(path, 'w', encoding='utf-8', newline='').write('\n'.join(out) + '\n')
    return added

def inject_vue_or_html(path):
    lines = open(path, encoding='utf-8').read().splitlines()
    out, added = [], 0
    mode = 'html'      # html | script | style
    in_tag = False     # 多行标签属性
    in_bt = False      # JS 模板字符串
    in_block = False
    for ln in lines:
        st = ln.strip()
        low = st.lower()
        if mode == 'html' and low.startswith('<script'):
            mode = 'script'
            out.append(ln); continue
        if mode == 'html' and low.startswith('<style'):
            mode = 'style'; out.append(ln); continue
        if mode == 'script' and low.startswith('</script>'):
            mode = 'html'; out.append(ln); continue
        if mode == 'style' and low.startswith('</style>'):
            mode = 'html'; out.append(ln); continue
        if mode == 'style':
            out.append(ln); continue
        if mode == 'html':
            if in_tag:
                out.append(ln)
                if '>' in ln: in_tag = False
                continue
            skip = (not st) or st.startswith('<!--') or (out and out[-1].strip().startswith('<!--'))
            if not skip and ('<' in st and '>' not in st):
                out.append(ln); in_tag = True; continue
            if not skip:
                out.append('<!-- ' + trunc(st) + ' -->'); added += 1
            out.append(ln); continue
        # script 段：JS 规则
        if in_block:
            out.append(ln)
            if '*/' in ln: in_block = False
            continue
        if in_bt:
            out.append(ln)
            if '`' in ln: in_bt = False
            continue
        if st.startswith('/*') and '*/' not in st:
            in_block = True; out.append(ln); continue
        if '`' in ln:
            if ln.count('`') % 2 == 1: in_bt = True
            skip = (not st) or is_comment_line(st, 'js') or (out and out[-1].strip().startswith('//'))
            if not skip:
                out.append(js_desc(st)); added += 1
            out.append(ln); continue
        skip = (not st) or is_comment_line(st, 'js') or (out and out[-1].strip().startswith('//'))
        if not skip:
            out.append(js_desc(st)); added += 1
        out.append(ln)
    open(path, 'w', encoding='utf-8', newline='').write('\n'.join(out) + '\n')
    return added

def inject_js(path):
    lines = open(path, encoding='utf-8').read().splitlines()
    out, added = [], 0
    in_bt = False
    in_block = False
    for ln in lines:
        st = ln.strip()
        if in_block:
            out.append(ln)
            if '*/' in ln: in_block = False
            continue
        if in_bt:
            out.append(ln)
            if '`' in ln: in_bt = False
            continue
        if st.startswith('/*') and '*/' not in st:
            in_block = True; out.append(ln); continue
        if '`' in ln:
            if ln.count('`') % 2 == 1: in_bt = True
            skip = (not st) or is_comment_line(st, 'js') or (out and out[-1].strip().startswith('//'))
            if not skip:
                out.append(js_desc(st)); added += 1
            out.append(ln); continue
        skip = (not st) or is_comment_line(st, 'js') or (out and out[-1].strip().startswith('//'))
        if not skip:
            out.append(js_desc(st)); added += 1
        out.append(ln)
    open(path, 'w', encoding='utf-8', newline='').write('\n'.join(out) + '\n')
    return added

def inject_sql(path):
    lines = open(path, encoding='utf-8').read().splitlines()
    out, added = [], 0
    for ln in lines:
        st = ln.strip()
        skip = (not st) or st.startswith('--') or st.startswith('/*') or (out and out[-1].strip().startswith('--'))
        if not skip:
            out.append(sql_desc(st)); added += 1
        out.append(ln)
    open(path, 'w', encoding='utf-8', newline='').write('\n'.join(out) + '\n')
    return added

total = 0
files_out = []
for base, dirs, files in os.walk(os.path.join(ROOT, 'server', 'shop')):
    dirs[:] = [d for d in dirs if d not in ('web', 'docs')]
    for f in files:
        p = os.path.join(base, f)
        if f.endswith('.go'):
            n = inject_go(p); total += n; files_out.append((p, n))
for f in os.listdir(os.path.join(ROOT, 'server', 'cmd', 'shop')):
    p = os.path.join(ROOT, 'server', 'cmd', 'shop', f)
    if f.endswith('.go'):
        n = inject_go(p); total += n; files_out.append((p, n))
for base, dirs, files in os.walk(os.path.join(ROOT, 'web', 'src', 'view', 'shop')):
    for f in files:
        p = os.path.join(base, f)
        if f.endswith('.vue'):
            n = inject_vue_or_html(p); total += n; files_out.append((p, n))
p = os.path.join(ROOT, 'web', 'src', 'api', 'shop.js')
if os.path.exists(p):
    n = inject_js(p); total += n; files_out.append((p, n))
p = os.path.join(ROOT, 'server', 'shop', 'web', 'web.go')
if os.path.exists(p):
    n = inject_go(p); total += n; files_out.append((p, n))
for f in ('buyer.html', 'payment.html'):
    p = os.path.join(ROOT, 'server', 'shop', 'web', f)
    if os.path.exists(p):
        n = inject_vue_or_html(p); total += n; files_out.append((p, n))
for base, dirs, files in os.walk(os.path.join(ROOT, 'deploy', 'shop')):
    for f in files:
        if f.endswith('.sql'):
            p = os.path.join(base, f)
            n = inject_sql(p); total += n; files_out.append((p, n))
for base, dirs, files in os.walk(os.path.join(ROOT, 'server', 'shop', 'migrations')):
    for f in files:
        if f.endswith('.sql'):
            p = os.path.join(base, f)
            n = inject_sql(p); total += n; files_out.append((p, n))

for p, n in files_out:
    print(('%6d  %s' % (n, os.path.relpath(p, ROOT))).encode('utf-8').decode('utf-8'))
print('TOTAL comments added:', total)
