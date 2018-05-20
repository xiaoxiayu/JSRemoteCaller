#coding=utf-8

import requests
import json
import base64
import time
import codecs
import getpass

__js_test_code = {}
def _parse_js_testcode(noCase=True):
    global __js_test_code
    fp = codecs.open('fxqa_api_testcode.js', 'r', 'utf8')
    str_lines = fp.readlines()
    fp.close()
    testcode_map = {}
    b_start = False
    testname = ''
    testcode = ''
    for line in str_lines:
        strip_line = line.strip().replace(' ', '')
        if strip_line.find('////TestFunction:') != -1:
            testname = strip_line[strip_line.find(':') + 1:]
            #print(testname)
            b_start = True
            continue
        elif strip_line.find('////TestEnd:') != -1:
            b_start = False
            continue
        #elif strip_line.find('//') != -1:
            #continue

        if b_start:
            testcode += line
        else:
            if noCase:
                testname = testname.lower()
            __js_test_code[testname] = testcode
            testcode = ''
            testname = ''
            


def fxqa_js_code_format(code_s):
    return base64.b64encode(code_s.encode('utf8')).decode('utf8')


def fxqa_get_teststatus(platform_name):
    if platform_name == 'windows':
        res_p = requests.get('http://127.0.0.1:9091/test-status')
    elif platform_name == 'ios':
        res_p = requests.get('http://127.0.0.1:9092/test-status')
    else:
        res_p = requests.get('http://127.0.0.1:9093/test-status')
    ret = json.loads(res_p.text)
    return ret['status']


def fxqa_get_log(platform_name, limit_t = 60):
    timeout = 0
    print("log")
    while timeout < limit_t:
        timeout += 1
        if platform_name == 'windows':
            res_p = requests.get('http://127.0.0.1:9091/log')
        elif platform_name == 'ios':
            res_p = requests.get('http://127.0.0.1:9092/log')
        else:
            res_p = requests.get('http://127.0.0.1:9093/log')
        
        print(res_p.text)
        if res_p.text == '':
            time.sleep(1)
            continue
        ret = json.loads(res_p.text)
        return (ret['Err'], ret['Info'])
    

def fxqa_run(platform_name, code_s):
    code_s = fxqa_js_code_format(code_s)
    data = {'code':code_s}
    if platform_name == 'windows':
        res_p = requests.post('http://127.0.0.1:9091/code', data=data)
        res_p = requests.get('http://127.0.0.1:9091/wait')
    elif platform_name == 'ios':
        res_p = requests.post('http://127.0.0.1:9092/code', data=data)
        res_p = requests.get('http://127.0.0.1:9092/wait')
    else:
        res_p = requests.post('http://127.0.0.1:9093/code', data=data)
        #print(res_p.text)
        res_p = requests.get('http://127.0.0.1:9093/wait')
    

def _get_testcode(js_func_name, js_func_param=None):
    global __js_test_code
    if js_func_param == None:
        call_func = ';%s();' % (js_func_name)
    else:
        call_func = ';%s(%s);' % (js_func_name, str(js_func_param))
    
    testcode = __js_test_code[js_func_name.lower()] + call_func
    #print(testcode)
    return testcode


def fxqa_test(platform_name, js_func_name, js_func_param=None):
    js = _get_testcode(js_func_name, js_func_param)
    fxqa_run(platform_name, js)
    return fxqa_get_log(platform_name)

    
def fxqa_test_all(platform_name):
    global __js_test_code
    for testname in __js_test_code.keys():
        #print(testname)
        fxqa_test(platform_name, testname)
        
_parse_js_testcode()

#fxqa_test_all()
##class User:
##    def __init__(self):
##        pass
##
##    def getUserToken(self):
##        js = fxqa_js_set('''FX.app.user.getUsearToken().then(function(result){fxqa_log(0, result);});''')
##        fxqa_run(js)
##        return fxqa_get_log()
##    
##	    
##    def getUserId(self):
##        js = fxqa_js_set('''FX.app.user.getUserId().then(function(result){fxqa_log(0, result);});''')
##        fxqa_run(js)
##        return fxqa_get_log()
##
##class App:
##    def __init__(self):
##        pass
##    
##    def getCurDoc(self):
##        js = _get_testcode('app_getCurDoc')
##        fxqa_run(js)
##        return fxqa_get_log()
##
##    def openDoc(self, docpath):
##        js = '''FX.app.openDoc('%s');fxqa_log(0, "openDoc runned");''' % docpath
##        fxqa_run(js)
##        return fxqa_get_log()
##
##    def closeDoc(self):
##        js = '''FX.app.getCurDoc().then(function(curDoc){FX.app.closeDoc(curDoc);fxqa_log(0, "closeDoc runned")});'''
##        fxqa_run(js)
##        return fxqa_get_log()
##
##    def signOut(self):
##        js = '''FX.app.signOut();fxqa_log(0, "signOut runned");'''
##        fxqa_run(js)
##        return fxqa_get_log()



