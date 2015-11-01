from oss.oss_api import *
from oss import oss_xml_handler
endpoint="127.0.0.1:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)
oss.debug=True

res=oss.copy_object("source-bucket","source-object","target-bucket","target-object")
print "%s\n%s" % (res.status,res.getheaders())
