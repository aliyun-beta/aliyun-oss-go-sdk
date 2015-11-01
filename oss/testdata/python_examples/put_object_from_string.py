from oss.oss_api import *
from oss import oss_xml_handler
endpoint="hostlocal:9999"
endpoint="oss-cn-hangzhou.aliyuncs.com:9999"
id, secret = "ayahghai0juiSie", "quitie*ph3Lah{F"
oss = OssAPI(endpoint, id, secret)

res=oss.put_object_from_string("bucket-name","object/name","wefpofjwefew")
